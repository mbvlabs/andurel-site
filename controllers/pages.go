package controllers

import (
	"errors"
	"net/http"

	"andurel-site/internal/hypermedia"
	"andurel-site/router"
	"andurel-site/router/routes"
	"andurel-site/views"

	"github.com/labstack/echo/v5"
)

type Pages struct{}

func NewPages() Pages {
	return Pages{}
}

func (p Pages) RegisterRoutes(r *router.Router) error {
	errs := []error{}

	_, err := r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.HomePage.Path(),
		Name:    routes.HomePage.Name(),
		Handler: p.Home,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodHead,
		Path:    routes.HomePage.Path(),
		Name:    routes.HomePage.Name() + ".head",
		Handler: p.Home,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_ = r.AddRouteNotFound(p.NotFound)

	return errors.Join(errs...)
}

func (p Pages) Home(etx *echo.Context) error {
	return hypermedia.RenderPage(etx, views.Welcome{}.Page())
}

func (p Pages) NotFound(etx *echo.Context) error {
	return renderNotFound(etx)
}
