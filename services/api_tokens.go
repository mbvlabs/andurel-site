package services

import (
	"context"
	"errors"
	"strings"
	"time"
	"uuid"

	"andurel-site/config"
	"andurel-site/models"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mbvlabs/andurel/pkg/validation"
)

const minAPITokenLength = 16

var (
	ErrAPITokenInvalid = errors.New("api token invalid")
	ErrAPITokenExpired = errors.New("api token expired")
	ErrAPITokenTaken   = errors.New("api token already in use")
)

type APITokens struct {
	tokens models.Tokens
	secret string
}

func NewAPITokens(tokens models.Tokens, authCfg config.Auth) APITokens {
	return APITokens{
		tokens: tokens,
		secret: authCfg.TokenSigningKey,
	}
}

type CreateAPITokenInput struct {
	Name      string
	Token     string
	ExpiresAt pgtype.Timestamptz
	CreatedBy uuid.UUID
}

type UpdateAPITokenInput struct {
	ID        uuid.UUID
	Name      string
	Token     string
	ExpiresAt pgtype.Timestamptz
}

type IssuedAPIToken struct {
	Token models.Token
	Meta  models.APITokenMeta
	Plain string
}

func (a APITokens) List(ctx context.Context) ([]models.Token, error) {
	return a.tokens.ListByScope(ctx, models.ScopeAPI)
}

func (a APITokens) Find(ctx context.Context, id uuid.UUID) (models.Token, error) {
	token, err := a.tokens.Find(ctx, id)
	if err != nil {
		return models.Token{}, err
	}
	if token.Scope != models.ScopeAPI {
		return models.Token{}, models.ErrNotFound
	}

	return token, nil
}

func (a APITokens) Create(ctx context.Context, input CreateAPITokenInput) (IssuedAPIToken, error) {
	name := strings.TrimSpace(input.Name)
	b := validation.NewBuilder()
	b.Required("name", name)
	b.MaxLen("name", name, 255)
	if !input.ExpiresAt.Valid {
		b.Required("expiresAt", "")
	}
	if !b.Errors().Empty() {
		return IssuedAPIToken{}, b.Errors()
	}

	plain, err := normalizeAPITokenValue(input.Token)
	if err != nil {
		return IssuedAPIToken{}, err
	}

	hash := models.HashForStorage(plain, a.secret)
	if err := a.ensureHashAvailable(ctx, hash, uuid.UUID{}); err != nil {
		return IssuedAPIToken{}, err
	}

	meta := models.APITokenMeta{
		Name:      name,
		CreatedBy: input.CreatedBy,
		Prefix:    models.TokenPrefix(plain),
	}
	raw, err := meta.Bytes()
	if err != nil {
		return IssuedAPIToken{}, err
	}

	created, issued, err := a.tokens.Issue(ctx, a.secret, models.ScopeAPI, input.ExpiresAt, raw, plain)
	if err != nil {
		return IssuedAPIToken{}, err
	}

	return IssuedAPIToken{Token: created, Meta: meta, Plain: issued}, nil
}

func (a APITokens) Update(ctx context.Context, input UpdateAPITokenInput) (IssuedAPIToken, error) {
	name := strings.TrimSpace(input.Name)
	b := validation.NewBuilder()
	b.Required("name", name)
	b.MaxLen("name", name, 255)
	if !input.ExpiresAt.Valid {
		b.Required("expiresAt", "")
	}
	if !b.Errors().Empty() {
		return IssuedAPIToken{}, b.Errors()
	}

	existing, err := a.Find(ctx, input.ID)
	if err != nil {
		return IssuedAPIToken{}, err
	}

	meta, err := existing.APIMeta()
	if err != nil {
		return IssuedAPIToken{}, err
	}

	hash := existing.Hash
	plain := ""
	if strings.TrimSpace(input.Token) != "" {
		plain, err = normalizeAPITokenValue(input.Token)
		if err != nil {
			return IssuedAPIToken{}, err
		}
		hash = models.HashForStorage(plain, a.secret)
		if err := a.ensureHashAvailable(ctx, hash, existing.ID); err != nil {
			return IssuedAPIToken{}, err
		}
		meta.Prefix = models.TokenPrefix(plain)
	}
	meta.Name = name

	raw, err := meta.Bytes()
	if err != nil {
		return IssuedAPIToken{}, err
	}

	updated, err := a.tokens.Update(ctx, models.UpdateTokenData{
		ID:        existing.ID,
		Scope:     models.ScopeAPI,
		ExpiresAt: input.ExpiresAt,
		Hash:      hash,
		MetaData:  raw,
	})
	if err != nil {
		return IssuedAPIToken{}, err
	}

	return IssuedAPIToken{Token: updated, Meta: meta, Plain: plain}, nil
}

func (a APITokens) Destroy(ctx context.Context, id uuid.UUID) error {
	if _, err := a.Find(ctx, id); err != nil {
		return err
	}

	return a.tokens.Destroy(ctx, id)
}

func (a APITokens) Authenticate(ctx context.Context, plain string) (models.Token, error) {
	plain = strings.TrimSpace(plain)
	if plain == "" {
		return models.Token{}, ErrAPITokenInvalid
	}

	token, err := a.tokens.FindByScopeAndHash(ctx, a.secret, models.ScopeAPI, plain)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return models.Token{}, ErrAPITokenInvalid
		}
		return models.Token{}, err
	}
	if token.IsExpired(time.Now()) {
		return models.Token{}, ErrAPITokenExpired
	}

	if err := a.touchLastUsed(ctx, token); err != nil {
		return models.Token{}, err
	}

	return token, nil
}

func (a APITokens) touchLastUsed(ctx context.Context, token models.Token) error {
	meta, err := token.APIMeta()
	if err != nil {
		return err
	}
	now := time.Now()
	meta.LastUsedAt = &now
	raw, err := meta.Bytes()
	if err != nil {
		return err
	}

	_, err = a.tokens.Update(ctx, models.UpdateTokenData{
		ID:        token.ID,
		Scope:     token.Scope,
		ExpiresAt: token.ExpiresAt,
		Hash:      token.Hash,
		MetaData:  raw,
	})
	return err
}

func (a APITokens) ensureHashAvailable(ctx context.Context, hash string, ignoreID uuid.UUID) error {
	existing, err := a.tokens.FindByScopeHash(ctx, models.ScopeAPI, hash)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil
		}
		return err
	}
	if existing.ID == ignoreID {
		return nil
	}

	return ErrAPITokenTaken
}

func normalizeAPITokenValue(value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		generated, err := models.GenerateSecureToken()
		if err != nil {
			return "", err
		}
		return generated, nil
	}

	b := validation.NewBuilder()
	if len(trimmed) < minAPITokenLength {
		b.Add("token", "min", "must be at least 16 characters")
	}
	if !b.Errors().Empty() {
		return "", b.Errors()
	}

	return trimmed, nil
}
