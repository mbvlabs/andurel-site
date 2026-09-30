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

type DocumentationVersions struct {
	documentation services.Documentation
	renderer      *inertia.Renderer
}

func NewDocumentationVersions(documentation services.Documentation, renderer *inertia.Renderer) DocumentationVersions {
	return DocumentationVersions{documentation: documentation, renderer: renderer}
}

func (dv DocumentationVersions) RegisterRoutes(r *router.Router) error {
	var errs []error

	adminOnly := []echo.MiddlewareFunc{middleware.AdminOnly()}

	_, err := r.AddRoute(routes.AdminDocumentationVersionIndex, echo.Route{
		Method:      http.MethodGet,
		Handler:     dv.Index,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(routes.AdminDocumentationVersionNew, echo.Route{
		Method:      http.MethodGet,
		Handler:     dv.New,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(routes.AdminDocumentationVersionCreate, echo.Route{
		Method:      http.MethodPost,
		Handler:     dv.Create,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(routes.AdminDocumentationVersionShow, echo.Route{
		Method:      http.MethodGet,
		Handler:     dv.Show,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(routes.AdminDocumentationVersionUpdate, echo.Route{
		Method:      http.MethodPut,
		Handler:     dv.Update,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(routes.AdminDocumentationVersionCreatePage, echo.Route{
		Method:      http.MethodPost,
		Handler:     dv.CreatePage,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(routes.AdminDocumentationVersionUpdateNav, echo.Route{
		Method:      http.MethodPut,
		Handler:     dv.UpdateNav,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(routes.AdminDocumentationVersionPublishNav, echo.Route{
		Method:      http.MethodPost,
		Handler:     dv.PublishNav,
		Middlewares: adminOnly,
	})
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func (dv DocumentationVersions) Index(etx *echo.Context) error {
	versions, err := dv.documentation.ListVersions(etx.Request().Context())
	if err != nil {
		return dv.renderer.Page(etx, "Errors/InternalError", inertia.Props{}).Render()
	}

	items := make([]DocumentationVersionData, 0, len(versions))
	for _, version := range versions {
		items = append(items, newDocumentationVersionData(version))
	}

	return dv.renderer.Page(
		etx,
		"Admin/DocumentationVersion/Index",
		inertia.FromStruct(DocumentationVersionIndexProps{Versions: items}),
	).Render()
}

func (dv DocumentationVersions) New(etx *echo.Context) error {
	return dv.renderer.Page(etx, "Admin/DocumentationVersion/Create", inertia.Props{}).Render()
}

func (dv DocumentationVersions) Create(etx *echo.Context) error {
	ctx, span := telemetry.From(etx, "admin.documentation_versions.create")
	defer span.End()

	userID, err := currentAdminUserID(etx)
	if err != nil {
		return dv.renderer.Page(etx, "Errors/InternalError", inertia.Props{}).Render()
	}

	var payload CreateDocumentationVersionFormPayload
	if err := etx.Bind(&payload); err != nil {
		telemetry.Error(ctx, "could not parse CreateDocumentationVersionFormPayload", "error", err)
		return dv.renderer.Page(etx, "Errors/BadRequest", inertia.Props{}).Render()
	}

	version, err := dv.documentation.CreateVersion(ctx, userID, services.CreateDocumentationVersionInput{
		Slug:     payload.Slug,
		Label:    payload.Label,
		IsLatest: payload.IsLatest,
		Position: payload.Position,
	})
	if err != nil {
		flashAdminError(etx, err)
		return dv.renderer.Redirect(etx, routes.AdminDocumentationVersionNew.URL(), http.StatusSeeOther)
	}

	kiks.AddFlash(etx.Request().Context(), kiks.FlashSuccess, "Version created")
	return dv.renderer.Redirect(etx, routes.AdminDocumentationVersionShow.URL(version.ID), http.StatusSeeOther)
}

func (dv DocumentationVersions) Show(etx *echo.Context) error {
	versionID, err := strconv.ParseInt(etx.Param("id"), 10, 64)
	if err != nil {
		return dv.renderer.Page(etx, "Errors/BadRequest", inertia.Props{}).Render()
	}

	userID, err := currentAdminUserID(etx)
	if err != nil {
		return dv.renderer.Page(etx, "Errors/InternalError", inertia.Props{}).Render()
	}

	workspace, err := dv.documentation.VersionWorkspace(etx.Request().Context(), versionID, userID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return dv.renderer.Page(etx, "Errors/NotFound", inertia.Props{}).Render()
		}
		return dv.renderer.Page(etx, "Errors/InternalError", inertia.Props{}).Render()
	}

	pages := make([]DocumentationPageSummaryData, 0, len(workspace.Pages))
	for _, page := range workspace.Pages {
		pages = append(pages, newDocumentationPageSummaryData(page))
	}

	return dv.renderer.Page(
		etx,
		"Admin/DocumentationVersion/Show",
		inertia.FromStruct(DocumentationVersionItemProps{
			Version:  newDocumentationVersionData(workspace.Version),
			Pages:    pages,
			NavDraft: newDocumentationNavDraftData(workspace.NavDraft),
		}),
	).Render()
}

func (dv DocumentationVersions) Update(etx *echo.Context) error {
	ctx, span := telemetry.From(etx, "admin.documentation_versions.update")
	defer span.End()

	versionID, err := strconv.ParseInt(etx.Param("id"), 10, 64)
	if err != nil {
		return dv.renderer.Page(etx, "Errors/BadRequest", inertia.Props{}).Render()
	}

	var payload UpdateDocumentationVersionFormPayload
	if err := etx.Bind(&payload); err != nil {
		telemetry.Error(ctx, "could not parse UpdateDocumentationVersionFormPayload", "error", err)
		return dv.renderer.Page(etx, "Errors/BadRequest", inertia.Props{}).Render()
	}

	version, err := dv.documentation.UpdateVersion(ctx, services.UpdateDocumentationVersionInput{
		ID:       versionID,
		Slug:     payload.Slug,
		Label:    payload.Label,
		IsLatest: payload.IsLatest,
		Position: payload.Position,
	})
	if err != nil {
		flashAdminError(etx, err)
		return dv.renderer.Redirect(etx, routes.AdminDocumentationVersionShow.URL(versionID), http.StatusSeeOther)
	}

	kiks.AddFlash(etx.Request().Context(), kiks.FlashSuccess, "Version saved")
	return dv.renderer.Redirect(etx, routes.AdminDocumentationVersionShow.URL(version.ID), http.StatusSeeOther)
}

func (dv DocumentationVersions) CreatePage(etx *echo.Context) error {
	ctx, span := telemetry.From(etx, "admin.documentation_versions.create_page")
	defer span.End()

	versionID, err := strconv.ParseInt(etx.Param("id"), 10, 64)
	if err != nil {
		return dv.renderer.Page(etx, "Errors/BadRequest", inertia.Props{}).Render()
	}

	userID, err := currentAdminUserID(etx)
	if err != nil {
		return dv.renderer.Page(etx, "Errors/InternalError", inertia.Props{}).Render()
	}

	var payload CreateDocumentationPageFormPayload
	if err := etx.Bind(&payload); err != nil {
		telemetry.Error(ctx, "could not parse CreateDocumentationPageFormPayload", "error", err)
		return dv.renderer.Page(etx, "Errors/BadRequest", inertia.Props{}).Render()
	}

	page, err := dv.documentation.CreatePage(ctx, services.CreateDocumentationPageInput{
		VersionID: versionID,
		Slug:      payload.Slug,
		Title:     payload.Title,
		CreatedBy: userID,
	})
	if err != nil {
		flashAdminError(etx, err)
		return dv.renderer.Redirect(etx, routes.AdminDocumentationVersionShow.URL(versionID), http.StatusSeeOther)
	}

	kiks.AddFlash(etx.Request().Context(), kiks.FlashSuccess, "Page created")
	return dv.renderer.Redirect(etx, routes.AdminDocumentationPageShow.URL(page.ID), http.StatusSeeOther)
}

func (dv DocumentationVersions) UpdateNav(etx *echo.Context) error {
	ctx, span := telemetry.From(etx, "admin.documentation_versions.update_nav")
	defer span.End()

	versionID, err := strconv.ParseInt(etx.Param("id"), 10, 64)
	if err != nil {
		return dv.renderer.Page(etx, "Errors/BadRequest", inertia.Props{}).Render()
	}

	userID, err := currentAdminUserID(etx)
	if err != nil {
		return dv.renderer.Page(etx, "Errors/InternalError", inertia.Props{}).Render()
	}

	var payload UpdateDocumentationNavFormPayload
	if err := etx.Bind(&payload); err != nil {
		telemetry.Error(ctx, "could not parse UpdateDocumentationNavFormPayload", "error", err)
		return dv.renderer.Page(etx, "Errors/BadRequest", inertia.Props{}).Render()
	}

	if _, err := dv.documentation.SaveNavDraft(ctx, services.SaveNavDraftInput{
		VersionID: versionID,
		Tree:      mapNavNodesToModels(payload.Tree),
		CreatedBy: userID,
	}); err != nil {
		flashAdminError(etx, err)
		return dv.renderer.Redirect(etx, routes.AdminDocumentationVersionShow.URL(versionID), http.StatusSeeOther)
	}

	kiks.AddFlash(etx.Request().Context(), kiks.FlashSuccess, "Nav draft saved")
	return dv.renderer.Redirect(etx, routes.AdminDocumentationVersionShow.URL(versionID), http.StatusSeeOther)
}

func (dv DocumentationVersions) PublishNav(etx *echo.Context) error {
	ctx, span := telemetry.From(etx, "admin.documentation_versions.publish_nav")
	defer span.End()

	versionID, err := strconv.ParseInt(etx.Param("id"), 10, 64)
	if err != nil {
		return dv.renderer.Page(etx, "Errors/BadRequest", inertia.Props{}).Render()
	}

	userID, err := currentAdminUserID(etx)
	if err != nil {
		return dv.renderer.Page(etx, "Errors/InternalError", inertia.Props{}).Render()
	}

	var payload UpdateDocumentationNavFormPayload
	if err := etx.Bind(&payload); err != nil {
		telemetry.Error(ctx, "could not parse UpdateDocumentationNavFormPayload", "error", err)
		return dv.renderer.Page(etx, "Errors/BadRequest", inertia.Props{}).Render()
	}

	if _, err := dv.documentation.SaveNavDraft(ctx, services.SaveNavDraftInput{
		VersionID: versionID,
		Tree:      mapNavNodesToModels(payload.Tree),
		CreatedBy: userID,
	}); err != nil {
		flashAdminError(etx, err)
		return dv.renderer.Redirect(etx, routes.AdminDocumentationVersionShow.URL(versionID), http.StatusSeeOther)
	}

	if _, err := dv.documentation.PublishNav(ctx, versionID, userID); err != nil {
		flashAdminError(etx, err)
		return dv.renderer.Redirect(etx, routes.AdminDocumentationVersionShow.URL(versionID), http.StatusSeeOther)
	}

	kiks.AddFlash(etx.Request().Context(), kiks.FlashSuccess, "Nav released")
	return dv.renderer.Redirect(etx, routes.AdminDocumentationVersionShow.URL(versionID), http.StatusSeeOther)
}
