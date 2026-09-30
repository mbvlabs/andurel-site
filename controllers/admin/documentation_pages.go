package admin

import (
	"errors"
	"net/http"
	"strconv"

	"andurel-site/models"
	"andurel-site/router"
	"andurel-site/router/middleware"
	"andurel-site/router/routes"
	"andurel-site/services"

	"github.com/labstack/echo/v5"
	"github.com/mbvlabs/andurel/pkg/inertia"
	"github.com/mbvlabs/andurel/pkg/kiks"
	"github.com/mbvlabs/andurel/pkg/telemetry"
)

type DocumentationPages struct {
	documentation services.Documentation
	renderer      *inertia.Renderer
}

func NewDocumentationPages(documentation services.Documentation, renderer *inertia.Renderer) DocumentationPages {
	return DocumentationPages{documentation: documentation, renderer: renderer}
}

func (dp DocumentationPages) RegisterRoutes(r *router.Router) error {
	var errs []error

	adminOnly := []echo.MiddlewareFunc{middleware.AdminOnly()}

	_, err := r.AddRoute(routes.AdminDocumentationPageShow, echo.Route{
		Method:      http.MethodGet,
		Handler:     dp.Show,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(routes.AdminDocumentationPageUpdate, echo.Route{
		Method:      http.MethodPut,
		Handler:     dp.Update,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(routes.AdminDocumentationPagePublish, echo.Route{
		Method:      http.MethodPost,
		Handler:     dp.Publish,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(routes.AdminDocumentationPageRestore, echo.Route{
		Method:      http.MethodPost,
		Handler:     dp.Restore,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(routes.AdminDocumentationPagePreview, echo.Route{
		Method:      http.MethodPost,
		Handler:     dp.Preview,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func (dp DocumentationPages) Show(etx *echo.Context) error {
	pageID, err := strconv.ParseInt(etx.Param("id"), 10, 64)
	if err != nil {
		return dp.renderer.Page(etx, "Errors/BadRequest", inertia.Props{}).Render()
	}

	userID, err := currentAdminUserID(etx)
	if err != nil {
		return dp.renderer.Page(etx, "Errors/InternalError", inertia.Props{}).Render()
	}

	editor, err := dp.documentation.PageEditor(etx.Request().Context(), pageID, userID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return dp.renderer.Page(etx, "Errors/NotFound", inertia.Props{}).Render()
		}
		return dp.renderer.Page(etx, "Errors/InternalError", inertia.Props{}).Render()
	}

	revisions := make([]DocumentationRevisionData, 0, len(editor.Revisions))
	for _, revision := range editor.Revisions {
		revisions = append(revisions, newDocumentationRevisionData(revision))
	}

	var published *DocumentationRevisionData
	if editor.Published != nil {
		data := newDocumentationRevisionData(*editor.Published)
		published = &data
	}

	return dp.renderer.Page(
		etx,
		"Admin/DocumentationPage/Show",
		inertia.FromStruct(DocumentationPageItemProps{
			Version:   newDocumentationVersionData(editor.Version),
			Page:      newDocumentationPageData(editor.Page),
			Draft:     newDocumentationRevisionData(editor.Draft),
			Published: published,
			Revisions: revisions,
		}),
	).Render()
}

func (dp DocumentationPages) Update(etx *echo.Context) error {
	ctx, span := telemetry.From(etx, "admin.documentation_pages.update")
	defer span.End()

	pageID, err := strconv.ParseInt(etx.Param("id"), 10, 64)
	if err != nil {
		return dp.renderer.Page(etx, "Errors/BadRequest", inertia.Props{}).Render()
	}

	userID, err := currentAdminUserID(etx)
	if err != nil {
		return dp.renderer.Page(etx, "Errors/InternalError", inertia.Props{}).Render()
	}

	var payload UpdateDocumentationPageFormPayload
	if err := etx.Bind(&payload); err != nil {
		telemetry.Error(ctx, "could not parse UpdateDocumentationPageFormPayload", "error", err)
		return dp.renderer.Page(etx, "Errors/BadRequest", inertia.Props{}).Render()
	}

	if _, err := dp.documentation.SavePageDraft(ctx, services.SavePageDraftInput{
		PageID:       pageID,
		Slug:         payload.Slug,
		Title:        payload.Title,
		MetaTitle:    payload.MetaTitle,
		Description:  payload.Description,
		BodyMarkdown: payload.BodyMarkdown,
		CreatedBy:    userID,
	}); err != nil {
		flashAdminError(etx, err)
		return dp.renderer.Redirect(etx, routes.AdminDocumentationPageShow.URL(pageID), http.StatusSeeOther)
	}

	kiks.AddFlash(etx.Request().Context(), kiks.FlashSuccess, "Draft saved")
	return dp.renderer.Redirect(etx, routes.AdminDocumentationPageShow.URL(pageID), http.StatusSeeOther)
}

func (dp DocumentationPages) Publish(etx *echo.Context) error {
	ctx, span := telemetry.From(etx, "admin.documentation_pages.publish")
	defer span.End()

	pageID, err := strconv.ParseInt(etx.Param("id"), 10, 64)
	if err != nil {
		return dp.renderer.Page(etx, "Errors/BadRequest", inertia.Props{}).Render()
	}

	userID, err := currentAdminUserID(etx)
	if err != nil {
		return dp.renderer.Page(etx, "Errors/InternalError", inertia.Props{}).Render()
	}

	var payload UpdateDocumentationPageFormPayload
	if err := etx.Bind(&payload); err != nil {
		telemetry.Error(ctx, "could not parse UpdateDocumentationPageFormPayload", "error", err)
		return dp.renderer.Page(etx, "Errors/BadRequest", inertia.Props{}).Render()
	}

	if payload.Title != "" || payload.BodyMarkdown != "" || payload.Slug != "" {
		if _, err := dp.documentation.SavePageDraft(ctx, services.SavePageDraftInput{
			PageID:       pageID,
			Slug:         payload.Slug,
			Title:        payload.Title,
			MetaTitle:    payload.MetaTitle,
			Description:  payload.Description,
			BodyMarkdown: payload.BodyMarkdown,
			CreatedBy:    userID,
		}); err != nil {
			flashAdminError(etx, err)
			return dp.renderer.Redirect(etx, routes.AdminDocumentationPageShow.URL(pageID), http.StatusSeeOther)
		}
	}

	if _, err := dp.documentation.PublishPage(ctx, pageID, userID); err != nil {
		flashAdminError(etx, err)
		return dp.renderer.Redirect(etx, routes.AdminDocumentationPageShow.URL(pageID), http.StatusSeeOther)
	}

	kiks.AddFlash(etx.Request().Context(), kiks.FlashSuccess, "Page released")
	return dp.renderer.Redirect(etx, routes.AdminDocumentationPageShow.URL(pageID), http.StatusSeeOther)
}

func (dp DocumentationPages) Restore(etx *echo.Context) error {
	ctx, span := telemetry.From(etx, "admin.documentation_pages.restore")
	defer span.End()

	pageID, err := strconv.ParseInt(etx.Param("id"), 10, 64)
	if err != nil {
		return dp.renderer.Page(etx, "Errors/BadRequest", inertia.Props{}).Render()
	}

	revisionID, err := strconv.ParseInt(etx.Param("revisionId"), 10, 64)
	if err != nil {
		return dp.renderer.Page(etx, "Errors/BadRequest", inertia.Props{}).Render()
	}

	userID, err := currentAdminUserID(etx)
	if err != nil {
		return dp.renderer.Page(etx, "Errors/InternalError", inertia.Props{}).Render()
	}

	if _, err := dp.documentation.RestoreRevision(ctx, pageID, revisionID, userID); err != nil {
		flashAdminError(etx, err)
		return dp.renderer.Redirect(etx, routes.AdminDocumentationPageShow.URL(pageID), http.StatusSeeOther)
	}

	kiks.AddFlash(etx.Request().Context(), kiks.FlashSuccess, "Opened as a new draft")
	return dp.renderer.Redirect(etx, routes.AdminDocumentationPageShow.URL(pageID), http.StatusSeeOther)
}

func (dp DocumentationPages) Preview(etx *echo.Context) error {
	ctx, span := telemetry.From(etx, "admin.documentation_pages.preview")
	defer span.End()

	pageID, err := strconv.ParseInt(etx.Param("id"), 10, 64)
	if err != nil {
		return etx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid page"})
	}

	if _, err := currentAdminUserID(etx); err != nil {
		return etx.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthenticated"})
	}

	var payload PreviewDocumentationPageFormPayload
	if err := etx.Bind(&payload); err != nil {
		telemetry.Error(ctx, "could not parse PreviewDocumentationPageFormPayload", "error", err)
		return etx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}

	html, headings, err := dp.documentation.RenderPreview(ctx, pageID, payload.BodyMarkdown)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return etx.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
		}
		telemetry.Error(ctx, "could not render documentation preview", "error", err)
		return etx.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "could not render markdown"})
	}

	return etx.JSON(http.StatusOK, map[string]any{
		"html":     html,
		"headings": headings,
	})
}
