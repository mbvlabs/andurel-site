package api

import (
	"errors"
	"net/http"
	"uuid"

	"andurel-site/models"
	"andurel-site/router"
	"andurel-site/router/middleware"
	"andurel-site/router/routes"
	"andurel-site/services"

	"github.com/labstack/echo/v5"
	"github.com/mbvlabs/andurel/pkg/validation"
)

type Documentations struct {
	documentation services.Documentation
	tokens        services.APITokens
}

func NewDocumentations(documentation services.Documentation, tokens services.APITokens) Documentations {
	return Documentations{
		documentation: documentation,
		tokens:        tokens,
	}
}

func (d Documentations) RegisterRoutes(r *router.Router) error {
	var errs []error
	auth := []echo.MiddlewareFunc{middleware.APITokenAuth(d.tokens)}

	_, err := r.AddRoute(routes.APIDocumentationIndex, echo.Route{
		Method:      http.MethodGet,
		Handler:     d.Index,
		Middlewares: auth,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(routes.APIDocumentationShow, echo.Route{
		Method:      http.MethodGet,
		Handler:     d.Show,
		Middlewares: auth,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(routes.APIDocumentationCreate, echo.Route{
		Method:      http.MethodPost,
		Handler:     d.Create,
		Middlewares: auth,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(routes.APIDocumentationUpdate, echo.Route{
		Method:      http.MethodPut,
		Handler:     d.Update,
		Middlewares: auth,
	})
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

type documentationPagePayload struct {
	Slug         string `json:"slug"`
	Title        string `json:"title"`
	MetaTitle    string `json:"metaTitle"`
	Description  string `json:"description"`
	BodyMarkdown string `json:"bodyMarkdown"`
}

type documentationNavPayload struct {
	Title    string                    `json:"title"`
	Slug     string                    `json:"slug"`
	Children []documentationNavPayload `json:"children"`
}

type documentationPayload struct {
	Slug     string                     `json:"slug"`
	Label    string                     `json:"label"`
	IsLatest bool                       `json:"isLatest"`
	Position int32                      `json:"position"`
	Publish  *bool                      `json:"publish,omitempty"`
	Pages    []documentationPagePayload `json:"pages"`
	Nav      []documentationNavPayload  `json:"nav"`
}

type documentationListItemPayload struct {
	Slug      string `json:"slug"`
	Label     string `json:"label"`
	IsLatest  bool   `json:"isLatest"`
	Position  int32  `json:"position"`
	PageCount int    `json:"pageCount"`
}

func (d Documentations) Index(etx *echo.Context) error {
	items, err := d.documentation.ListDocumentations(etx.Request().Context())
	if err != nil {
		return apiError(etx, http.StatusInternalServerError, "could not list documentations")
	}

	payload := make([]documentationListItemPayload, 0, len(items))
	for _, item := range items {
		payload = append(payload, documentationListItemPayload{
			Slug:      item.Slug,
			Label:     item.Label,
			IsLatest:  item.IsLatest,
			Position:  item.Position,
			PageCount: item.PageCount,
		})
	}

	return etx.JSON(http.StatusOK, map[string]any{"documentations": payload})
}

func (d Documentations) Show(etx *echo.Context) error {
	snapshot, err := d.documentation.GetDocumentation(etx.Request().Context(), etx.Param("slug"))
	if err != nil {
		return mapDocumentationAPIError(etx, err)
	}

	return etx.JSON(http.StatusOK, map[string]any{"documentation": snapshotPayload(snapshot)})
}

func (d Documentations) Create(etx *echo.Context) error {
	actorID, err := apiActorID(etx)
	if err != nil {
		return apiError(etx, http.StatusUnauthorized, "invalid token")
	}

	var payload documentationPayload
	if err := etx.Bind(&payload); err != nil {
		return apiError(etx, http.StatusBadRequest, "invalid json")
	}

	snapshot, err := d.documentation.CreateDocumentation(etx.Request().Context(), actorID, replaceInput(payload))
	if err != nil {
		return mapDocumentationAPIError(etx, err)
	}

	return etx.JSON(http.StatusCreated, map[string]any{"documentation": snapshotPayload(snapshot)})
}

func (d Documentations) Update(etx *echo.Context) error {
	actorID, err := apiActorID(etx)
	if err != nil {
		return apiError(etx, http.StatusUnauthorized, "invalid token")
	}

	var payload documentationPayload
	if err := etx.Bind(&payload); err != nil {
		return apiError(etx, http.StatusBadRequest, "invalid json")
	}

	snapshot, err := d.documentation.UpdateDocumentation(
		etx.Request().Context(),
		actorID,
		etx.Param("slug"),
		replaceInput(payload),
	)
	if err != nil {
		return mapDocumentationAPIError(etx, err)
	}

	return etx.JSON(http.StatusOK, map[string]any{"documentation": snapshotPayload(snapshot)})
}

func apiActorID(etx *echo.Context) (uuid.UUID, error) {
	token, ok := middleware.APITokenFromContext(etx.Request().Context())
	if !ok {
		return uuid.UUID{}, errors.New("unauthenticated")
	}

	meta, err := token.APIMeta()
	if err != nil {
		return uuid.UUID{}, err
	}

	return meta.CreatedBy, nil
}

func replaceInput(payload documentationPayload) services.ReplaceDocumentationInput {
	publish := true
	if payload.Publish != nil {
		publish = *payload.Publish
	}

	pages := make([]services.DocumentationPageSnapshot, 0, len(payload.Pages))
	for _, page := range payload.Pages {
		pages = append(pages, services.DocumentationPageSnapshot{
			Slug:         page.Slug,
			Title:        page.Title,
			MetaTitle:    page.MetaTitle,
			Description:  page.Description,
			BodyMarkdown: page.BodyMarkdown,
		})
	}

	return services.ReplaceDocumentationInput{
		Slug:     payload.Slug,
		Label:    payload.Label,
		IsLatest: payload.IsLatest,
		Position: payload.Position,
		Publish:  publish,
		Pages:    pages,
		Nav:      navInput(payload.Nav),
	}
}

func navInput(nodes []documentationNavPayload) []services.DocumentationNavSnapshotNode {
	out := make([]services.DocumentationNavSnapshotNode, 0, len(nodes))
	for _, node := range nodes {
		out = append(out, services.DocumentationNavSnapshotNode{
			Title:    node.Title,
			Slug:     node.Slug,
			Children: navInput(node.Children),
		})
	}

	return out
}

func snapshotPayload(snapshot services.DocumentationSnapshot) documentationPayload {
	pages := make([]documentationPagePayload, 0, len(snapshot.Pages))
	for _, page := range snapshot.Pages {
		pages = append(pages, documentationPagePayload{
			Slug:         page.Slug,
			Title:        page.Title,
			MetaTitle:    page.MetaTitle,
			Description:  page.Description,
			BodyMarkdown: page.BodyMarkdown,
		})
	}

	return documentationPayload{
		Slug:     snapshot.Slug,
		Label:    snapshot.Label,
		IsLatest: snapshot.IsLatest,
		Position: snapshot.Position,
		Pages:    pages,
		Nav:      navPayloadFromSnapshot(snapshot.Nav),
	}
}

func navPayloadFromSnapshot(nodes []services.DocumentationNavSnapshotNode) []documentationNavPayload {
	out := make([]documentationNavPayload, 0, len(nodes))
	for _, node := range nodes {
		out = append(out, documentationNavPayload{
			Title:    node.Title,
			Slug:     node.Slug,
			Children: navPayloadFromSnapshot(node.Children),
		})
	}

	return out
}

func mapDocumentationAPIError(etx *echo.Context, err error) error {
	if ve, ok := validation.As(err); ok && !ve.Empty() {
		return etx.JSON(http.StatusUnprocessableEntity, map[string]any{
			"error":  "validation failed",
			"fields": ve,
		})
	}
	switch {
	case errors.Is(err, models.ErrNotFound):
		return apiError(etx, http.StatusNotFound, "documentation not found")
	case errors.Is(err, services.ErrSlugTaken):
		return apiError(etx, http.StatusConflict, "slug already in use")
	case errors.Is(err, services.ErrNavPageMissing):
		return apiError(etx, http.StatusUnprocessableEntity, "nav references an unknown page slug")
	default:
		return apiError(etx, http.StatusInternalServerError, "could not process documentation")
	}
}

func apiError(etx *echo.Context, status int, message string) error {
	return etx.JSON(status, map[string]string{"error": message})
}
