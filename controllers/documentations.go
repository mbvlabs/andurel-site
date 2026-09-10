package controllers

import (
	"andurel-site/docs"
	"andurel-site/router"
	"andurel-site/router/routes"
	"errors"
	"net/http"

	"github.com/mbvlabs/andurel/pkg/inertia"

	"github.com/labstack/echo/v5"
)

type Documentations struct {
	renderer *inertia.Renderer
	site     *docs.Site
}

func NewDocumentations(renderer *inertia.Renderer, site *docs.Site) Documentations {
	return Documentations{renderer: renderer, site: site}
}

func (d Documentations) RegisterRoutes(r *router.Router) error {
	var errs []error
	var err error

	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.DocumentationIndex.Path(),
		Name:    routes.DocumentationIndex.Name(),
		Handler: d.Index,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.DocumentationVersion.Path(),
		Name:    routes.DocumentationVersion.Name(),
		Handler: d.Version,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.DocumentationShow.Path(),
		Name:    routes.DocumentationShow.Name(),
		Handler: d.Show,
	})
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func (d Documentations) Index(etx *echo.Context) error {
	return etx.Redirect(http.StatusPermanentRedirect, d.site.LatestURL())
}

func (d Documentations) Version(etx *echo.Context) error {
	url, ok := d.site.VersionURL(etx.Param("version"))
	if !ok {
		return d.notFound(etx)
	}
	return etx.Redirect(http.StatusPermanentRedirect, url)
}

func (d Documentations) Show(etx *echo.Context) error {
	document, ok := d.site.Find(etx.Param("version"), etx.Param("slug"))
	if !ok {
		return d.notFound(etx)
	}

	return d.renderer.Page(etx, "Documentation/Show", inertia.Props{
		"versions":       docs.Navigation(),
		"currentVersion": document.Version,
		"currentSlug":    document.Slug,
		"currentSection": document.Section,
		"title":          document.Title,
		"description":    document.Description,
		"html":           document.HTML,
		"headings":       document.Headings,
		"parent":         document.Parent,
		"previous":       document.Previous,
		"next":           document.Next,
	}).SSR().Render()
}

func (d Documentations) notFound(etx *echo.Context) error {
	return d.renderer.Page(etx, "Errors/NotFound", inertia.Props{}).
		Status(http.StatusNotFound).
		Render()
}
