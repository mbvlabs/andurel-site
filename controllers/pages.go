package controllers

import (
	"errors"
	"net/http"

	"andurel-site/router"
	"andurel-site/router/routes"

	"github.com/mbvlabs/andurel/pkg/inertia"

	"github.com/labstack/echo/v5"
)

type Pages struct {
	renderer *inertia.Renderer
}

func NewPages(renderer *inertia.Renderer) Pages {
	return Pages{renderer: renderer}
}

func (p Pages) RegisterRoutes(r *router.Router) error {
	errs := []error{}

	_, err := r.AddRoute(routes.HomePage, echo.Route{
		Method:  http.MethodGet,
		Handler: p.Home,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.AddRoute(routes.HomePage, echo.Route{
		Method:  http.MethodHead,
		Name:    routes.HomePage.Name() + ".head",
		Handler: p.Home,
	})
	if err != nil {
		errs = append(errs, err)
	}

	if err = r.AddRouteNotFoundEachHost(p.NotFound); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func (p Pages) Home(etx *echo.Context) error {
	return p.renderer.Page(etx, "Home", inertia.Props{}).SSR().Render()
}

func (p Pages) NotFound(etx *echo.Context) error {
	return p.renderer.Page(etx, "Errors/NotFound", inertia.Props{}).Render()
}
