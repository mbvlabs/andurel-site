package models

// andurel:table documentation_versions

import (
	"andurel-site/models/internal/queries"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mbvlabs/andurel/pkg/storage"
	"github.com/mbvlabs/andurel/pkg/validation"
)

type DocumentationVersions struct {
	queries *queries.Queries
}

func NewDocumentationVersions(db storage.Connection) DocumentationVersions {
	return DocumentationVersions{queries: queries.New(db)}
}

// WithTx returns a copy that runs queries inside tx.
func (dv DocumentationVersions) WithTx(tx storage.Transaction) DocumentationVersions {
	return DocumentationVersions{queries: queries.New(tx)}
}

type DocumentationVersion struct {
	ID             int64              `andurel:"id"`
	CreatedAt      pgtype.Timestamptz `andurel:"created_at"`
	UpdatedAt      pgtype.Timestamptz `andurel:"updated_at"`
	Slug           string             `andurel:"slug"`
	Label          string             `andurel:"label"`
	IsLatest       bool               `andurel:"is_latest"`
	Position       int32              `andurel:"position"`
	PublishedNavId pgtype.Int8        `andurel:"published_nav_id"`
}

func (dv DocumentationVersions) Find(ctx context.Context, id int64) (DocumentationVersion, error) {
	entity, err := dv.queries.GetDocumentationVersion[DocumentationVersion](ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DocumentationVersion{}, ErrNotFound
		}
		return DocumentationVersion{}, err
	}

	return entity, nil
}

type CreateDocumentationVersionData struct {
	Slug           string
	Label          string
	IsLatest       bool
	Position       int32
	PublishedNavId pgtype.Int8
}

func (dv DocumentationVersions) Create(ctx context.Context, data CreateDocumentationVersionData) (DocumentationVersion, error) {
	entity := DocumentationVersion{
		CreatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
		Slug:           data.Slug,
		Label:          data.Label,
		IsLatest:       data.IsLatest,
		Position:       data.Position,
		PublishedNavId: data.PublishedNavId,
	}

	if err := validation.Validate(&entity); err != nil {
		return DocumentationVersion{}, errors.Join(ErrDomainValidation, err)
	}

	return dv.queries.CreateDocumentationVersion[DocumentationVersion](ctx, queries.CreateDocumentationVersionParams{
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
		Slug:           entity.Slug,
		Label:          entity.Label,
		IsLatest:       entity.IsLatest,
		Position:       entity.Position,
		PublishedNavID: entity.PublishedNavId,
	})
}

type UpdateDocumentationVersionData struct {
	ID             int64
	Slug           string
	Label          string
	IsLatest       bool
	Position       int32
	PublishedNavId pgtype.Int8
}

func (dv DocumentationVersions) Update(ctx context.Context, data UpdateDocumentationVersionData) (DocumentationVersion, error) {
	entity := DocumentationVersion{
		ID:             data.ID,
		UpdatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
		Slug:           data.Slug,
		Label:          data.Label,
		IsLatest:       data.IsLatest,
		Position:       data.Position,
		PublishedNavId: data.PublishedNavId,
	}

	if err := validation.Validate(&entity); err != nil {
		return DocumentationVersion{}, errors.Join(ErrDomainValidation, err)
	}

	row, err := dv.queries.UpdateDocumentationVersion[DocumentationVersion](ctx, queries.UpdateDocumentationVersionParams{
		ID:             entity.ID,
		UpdatedAt:      entity.UpdatedAt,
		Slug:           entity.Slug,
		Label:          entity.Label,
		IsLatest:       entity.IsLatest,
		Position:       entity.Position,
		PublishedNavID: entity.PublishedNavId,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DocumentationVersion{}, ErrNotFound
		}
		return DocumentationVersion{}, err
	}

	return row, nil
}

func (dv DocumentationVersions) Destroy(ctx context.Context, id int64) error {
	if err := dv.queries.DeleteDocumentationVersion(ctx, id); err != nil {
		return err
	}

	return nil
}

func (dv DocumentationVersions) All(ctx context.Context) ([]DocumentationVersion, error) {
	return dv.queries.ListDocumentationVersions[DocumentationVersion](ctx).All()
}

type PaginatedDocumentationVersions struct {
	DocumentationVersions []DocumentationVersion
	TotalCount            int64
	Page                  int64
	PageSize              int64
	TotalPages            int64
}

func (dv DocumentationVersions) Paginate(ctx context.Context, page, pageSize int64) (PaginatedDocumentationVersions, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	totalCount, err := dv.queries.CountDocumentationVersions(ctx)
	if err != nil {
		return PaginatedDocumentationVersions{}, err
	}

	entities, err := dv.queries.ListDocumentationVersions[DocumentationVersion](ctx).
		Limit(int32(pageSize)).
		Offset(int32(offset)).
		All()
	if err != nil {
		return PaginatedDocumentationVersions{}, err
	}

	totalPages := (int64(totalCount) + pageSize - 1) / pageSize

	return PaginatedDocumentationVersions{
		DocumentationVersions: entities,
		TotalCount:            int64(totalCount),
		Page:                  page,
		PageSize:              pageSize,
		TotalPages:            totalPages,
	}, nil
}

func (dv DocumentationVersions) Upsert(ctx context.Context, id int64, data CreateDocumentationVersionData) (DocumentationVersion, error) {
	entity := DocumentationVersion{
		ID:             id,
		CreatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
		Slug:           data.Slug,
		Label:          data.Label,
		IsLatest:       data.IsLatest,
		Position:       data.Position,
		PublishedNavId: data.PublishedNavId,
	}

	if err := validation.Validate(&entity); err != nil {
		return DocumentationVersion{}, errors.Join(ErrDomainValidation, err)
	}

	return dv.queries.UpsertDocumentationVersion[DocumentationVersion](ctx, queries.UpsertDocumentationVersionParams{
		ID:             entity.ID,
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
		Slug:           entity.Slug,
		Label:          entity.Label,
		IsLatest:       entity.IsLatest,
		Position:       entity.Position,
		PublishedNavID: entity.PublishedNavId,
	})
}

func (dv DocumentationVersions) FindBySlug(ctx context.Context, slug string) (DocumentationVersion, error) {
	entity, err := dv.queries.GetDocumentationVersionBySlug[DocumentationVersion](ctx, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DocumentationVersion{}, ErrNotFound
		}
		return DocumentationVersion{}, err
	}

	return entity, nil
}

func (dv DocumentationVersions) AllByPosition(ctx context.Context) ([]DocumentationVersion, error) {
	return dv.queries.ListDocumentationVersionsByPosition[DocumentationVersion](ctx)
}

func (dv DocumentationVersions) ClearLatest(ctx context.Context) error {
	return dv.queries.ClearDocumentationVersionLatest(ctx, pgtype.Timestamptz{Time: time.Now(), Valid: true})
}
