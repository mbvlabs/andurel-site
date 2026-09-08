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
}

func NewDocumentations(renderer *inertia.Renderer) Documentations {
	return Documentations{renderer: renderer}
}

func (d Documentations) RegisterRoutes(r *router.Router) error {
	var errs []error
	var err error

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

func (d Documentations) Show(etx *echo.Context) error {
	versionName := etx.Param("version")
	slug := etx.Param("slug")
	_, section, page, ok := docs.Find(versionName, slug)
	if ok {
		return d.renderer.Page(etx, "Documentation/Show", inertia.Props{
			"versions":       docs.Navigation(),
			"currentVersion": versionName,
			"currentSlug":    page.Slug,
			"currentSection": section.Title,
			"title":          page.Title,
			"description":    page.Description,
		}).SSR().Render()
	}

	_, section, page, _ = docs.Find("latest", "latest")
		return d.renderer.Page(etx, "Documentation/Show", inertia.Props{
			"versions":       docs.Navigation(),
			"currentVersion": versionName,
			"currentSlug":    page.Slug,
			"currentSection": section.Title,
			"title":          page.Title,
			"description":    page.Description,
		}).SSR().Render()
}
