package admin

import (
	"errors"
	"net/http"
	"strings"
	"time"
	"uuid"

	"andurel-site/models"
	"andurel-site/router"
	"andurel-site/router/middleware"
	"andurel-site/router/routes"
	"andurel-site/services"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v5"
	"github.com/mbvlabs/andurel/pkg/inertia"
	"github.com/mbvlabs/andurel/pkg/kiks"
	"github.com/mbvlabs/andurel/pkg/telemetry"
)

type ApiTokens struct {
	tokens   services.APITokens
	renderer *inertia.Renderer
}

func NewApiTokens(tokens services.APITokens, renderer *inertia.Renderer) ApiTokens {
	return ApiTokens{tokens: tokens, renderer: renderer}
}

func (at ApiTokens) RegisterRoutes(r *router.Router) error {
	var errs []error
	adminOnly := []echo.MiddlewareFunc{middleware.AdminOnly()}

	_, err := r.AddRoute(routes.AdminApiTokenIndex, echo.Route{
		Method:      http.MethodGet,
		Handler:     at.Index,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(routes.AdminApiTokenShow, echo.Route{
		Method:      http.MethodGet,
		Handler:     at.Show,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(routes.AdminApiTokenNew, echo.Route{
		Method:      http.MethodGet,
		Handler:     at.New,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(routes.AdminApiTokenCreate, echo.Route{
		Method:      http.MethodPost,
		Handler:     at.Create,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(routes.AdminApiTokenUpdate, echo.Route{
		Method:      http.MethodPut,
		Handler:     at.Update,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(routes.AdminApiTokenDestroy, echo.Route{
		Method:      http.MethodDelete,
		Handler:     at.Destroy,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

type ApiTokenData struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	TokenPrefix string `json:"tokenPrefix"`
	ExpiresAt   string `json:"expiresAt"`
	LastUsedAt  string `json:"lastUsedAt"`
	CreatedAt   string `json:"createdAt"`
}

type ApiTokenIndexProps struct {
	Tokens []ApiTokenData `json:"tokens"`
}

type ApiTokenItemProps struct {
	Token         ApiTokenData `json:"token"`
	RevealedToken string       `json:"revealedToken"`
}

type CreateApiTokenFormPayload struct {
	Name      string `json:"name"`
	Token     string `json:"token"`
	ExpiresAt string `json:"expiresAt"`
}

type UpdateApiTokenFormPayload struct {
	Name      string `json:"name"`
	Token     string `json:"token"`
	ExpiresAt string `json:"expiresAt"`
}

func newApiTokenData(entity models.Token) ApiTokenData {
	meta, err := entity.APIMeta()
	if err != nil {
		meta = models.APITokenMeta{}
	}
	lastUsed := ""
	if meta.LastUsedAt != nil && !meta.LastUsedAt.IsZero() {
		lastUsed = meta.LastUsedAt.UTC().Format(time.RFC3339)
	}

	return ApiTokenData{
		ID:          entity.ID.String(),
		Name:        meta.Name,
		TokenPrefix: meta.Prefix,
		ExpiresAt:   timestamptzString(entity.ExpiresAt),
		LastUsedAt:  lastUsed,
		CreatedAt:   timestamptzString(entity.CreatedAt),
	}
}

func parseTokenID(etx *echo.Context) (uuid.UUID, error) {
	return uuid.Parse(etx.Param("id"))
}

func (at ApiTokens) Index(etx *echo.Context) error {
	tokens, err := at.tokens.List(etx.Request().Context())
	if err != nil {
		return at.renderer.Page(etx, "Errors/InternalError", inertia.Props{}).Render()
	}

	items := make([]ApiTokenData, 0, len(tokens))
	for _, token := range tokens {
		items = append(items, newApiTokenData(token))
	}

	return at.renderer.Page(
		etx,
		"Admin/ApiToken/Index",
		inertia.FromStruct(ApiTokenIndexProps{Tokens: items}),
	).Render()
}

func (at ApiTokens) Show(etx *echo.Context) error {
	tokenID, err := parseTokenID(etx)
	if err != nil {
		return at.renderer.Page(etx, "Errors/BadRequest", inertia.Props{}).Render()
	}

	token, err := at.tokens.Find(etx.Request().Context(), tokenID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return at.renderer.Page(etx, "Errors/NotFound", inertia.Props{}).Render()
		}
		return at.renderer.Page(etx, "Errors/InternalError", inertia.Props{}).Render()
	}

	return at.renderer.Page(
		etx,
		"Admin/ApiToken/Show",
		inertia.FromStruct(ApiTokenItemProps{Token: newApiTokenData(token)}),
	).Render()
}

func (at ApiTokens) New(etx *echo.Context) error {
	return at.renderer.Page(etx, "Admin/ApiToken/Create", inertia.Props{}).Render()
}

func (at ApiTokens) Create(etx *echo.Context) error {
	ctx, span := telemetry.From(etx, "admin.api_tokens.create")
	defer span.End()

	userID, err := currentAdminUserID(etx)
	if err != nil {
		return at.renderer.Page(etx, "Errors/InternalError", inertia.Props{}).Render()
	}

	var payload CreateApiTokenFormPayload
	if err := etx.Bind(&payload); err != nil {
		telemetry.Error(ctx, "could not parse CreateApiTokenFormPayload", "error", err)
		return at.renderer.Page(etx, "Errors/BadRequest", inertia.Props{}).Render()
	}

	issued, err := at.tokens.Create(ctx, services.CreateAPITokenInput{
		Name:      payload.Name,
		Token:     payload.Token,
		ExpiresAt: parseOptionalTimestamptz(payload.ExpiresAt),
		CreatedBy: userID,
	})
	if err != nil {
		flashAdminError(etx, err)
		return at.renderer.Redirect(etx, routes.AdminApiTokenNew.URL(), http.StatusSeeOther)
	}

	kiks.AddFlash(etx.Request().Context(), kiks.FlashSuccess, "Token created. Copy it now — it will not be shown again.")
	return at.renderer.Page(
		etx,
		"Admin/ApiToken/Show",
		inertia.FromStruct(ApiTokenItemProps{
			Token:         newApiTokenData(issued.Token),
			RevealedToken: issued.Plain,
		}),
	).Render()
}

func (at ApiTokens) Update(etx *echo.Context) error {
	ctx, span := telemetry.From(etx, "admin.api_tokens.update")
	defer span.End()

	tokenID, err := parseTokenID(etx)
	if err != nil {
		return at.renderer.Page(etx, "Errors/BadRequest", inertia.Props{}).Render()
	}

	var payload UpdateApiTokenFormPayload
	if err := etx.Bind(&payload); err != nil {
		telemetry.Error(ctx, "could not parse UpdateApiTokenFormPayload", "error", err)
		return at.renderer.Page(etx, "Errors/BadRequest", inertia.Props{}).Render()
	}

	issued, err := at.tokens.Update(ctx, services.UpdateAPITokenInput{
		ID:        tokenID,
		Name:      payload.Name,
		Token:     payload.Token,
		ExpiresAt: parseOptionalTimestamptz(payload.ExpiresAt),
	})
	if err != nil {
		flashAdminError(etx, err)
		return at.renderer.Redirect(etx, routes.AdminApiTokenShow.URL(tokenID), http.StatusSeeOther)
	}

	if issued.Plain != "" {
		kiks.AddFlash(etx.Request().Context(), kiks.FlashSuccess, "Token updated. Copy the new value now — it will not be shown again.")
		return at.renderer.Page(
			etx,
			"Admin/ApiToken/Show",
			inertia.FromStruct(ApiTokenItemProps{
				Token:         newApiTokenData(issued.Token),
				RevealedToken: issued.Plain,
			}),
		).Render()
	}

	kiks.AddFlash(etx.Request().Context(), kiks.FlashSuccess, "Token saved")
	return at.renderer.Redirect(etx, routes.AdminApiTokenShow.URL(issued.Token.ID), http.StatusSeeOther)
}

func (at ApiTokens) Destroy(etx *echo.Context) error {
	tokenID, err := parseTokenID(etx)
	if err != nil {
		return at.renderer.Page(etx, "Errors/BadRequest", inertia.Props{}).Render()
	}

	if err := at.tokens.Destroy(etx.Request().Context(), tokenID); err != nil {
		flashAdminError(etx, err)
		return at.renderer.Redirect(etx, routes.AdminApiTokenIndex.URL(), http.StatusSeeOther)
	}

	kiks.AddFlash(etx.Request().Context(), kiks.FlashSuccess, "Token deleted")
	return at.renderer.Redirect(etx, routes.AdminApiTokenIndex.URL(), http.StatusSeeOther)
}

func parseOptionalTimestamptz(value string) pgtype.Timestamptz {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return pgtype.Timestamptz{}
	}

	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02T15:04",
		"2006-01-02 15:04",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if parsed, err := time.ParseInLocation(layout, trimmed, time.Local); err == nil {
			return pgtype.Timestamptz{Time: parsed, Valid: true}
		}
	}

	return pgtype.Timestamptz{}
}
