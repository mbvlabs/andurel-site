package controllers

import (
	"errors"
	"net/http"

	documentation "andurel-site/docs"
	"andurel-site/internal/hypermedia"
	"andurel-site/router"
	"andurel-site/router/routes"
	"andurel-site/views"

	"github.com/labstack/echo/v5"
)

type Docs struct {
	site *documentation.Site
}

func NewDocs(site *documentation.Site) Docs {
	return Docs{site: site}
}

func (d Docs) RegisterRoutes(r *router.Router) error {
	definitions := []echo.Route{
		{
			Method:  http.MethodGet,
			Path:    routes.DocsSearch.Path(),
			Name:    routes.DocsSearch.Name(),
			Handler: d.Search,
		},
		{
			Method:  http.MethodGet,
			Path:    routes.DocsIndex.Path(),
			Name:    routes.DocsIndex.Name(),
			Handler: d.Index,
		},
		{
			Method:  http.MethodHead,
			Path:    routes.DocsIndex.Path(),
			Name:    routes.DocsIndex.Name() + ".head",
			Handler: d.Index,
		},
		{
			Method:  http.MethodGet,
			Path:    routes.DocsVersion.Path(),
			Name:    routes.DocsVersion.Name(),
			Handler: d.Version,
		},
		{
			Method:  http.MethodHead,
			Path:    routes.DocsVersion.Path(),
			Name:    routes.DocsVersion.Name() + ".head",
			Handler: d.Version,
		},
		{
			Method:  http.MethodGet,
			Path:    routes.DocsShow.Path(),
			Name:    routes.DocsShow.Name(),
			Handler: d.Show,
		},
		{
			Method:  http.MethodHead,
			Path:    routes.DocsShow.Path(),
			Name:    routes.DocsShow.Name() + ".head",
			Handler: d.Show,
		},
	}

	var errs []error
	for _, definition := range definitions {
		if _, err := r.AddRoute(definition); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (d Docs) Index(etx *echo.Context) error {
	return etx.Redirect(http.StatusPermanentRedirect, d.site.LatestURL())
}

func (d Docs) Version(etx *echo.Context) error {
	url, ok := d.site.VersionURL(etx.Param("version"))
	if !ok {
		return renderNotFound(etx)
	}
	return etx.Redirect(http.StatusPermanentRedirect, url)
}

func (d Docs) Show(etx *echo.Context) error {
	document, ok := d.site.Find(etx.Param("version"), etx.Param("slug"))
	if !ok {
		return renderNotFound(etx)
	}
	return hypermedia.RenderPage(etx, views.Documentation(d.site, document))
}

func (d Docs) Search(etx *echo.Context) error {
	return etx.JSON(http.StatusOK, d.site.SearchIndex())
}

func renderNotFound(etx *echo.Context) error {
	return hypermedia.RenderPage(etx, views.NotFound(), hypermedia.WithStatus(http.StatusNotFound))
}
