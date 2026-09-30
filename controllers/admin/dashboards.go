package admin

import (
	"errors"
	"net/http"

	"andurel-site/router"
	"andurel-site/router/middleware"
	"andurel-site/router/routes"

	"github.com/labstack/echo/v5"
	"github.com/mbvlabs/andurel/pkg/inertia"
)

type Dashboards struct {
	renderer *inertia.Renderer
}

func NewDashboards(renderer *inertia.Renderer) Dashboards {
	return Dashboards{renderer: renderer}
}

func (d Dashboards) RegisterRoutes(r *router.Router) error {
	var errs []error

	_, err := r.AddRoute(routes.AdminDashboardHome, echo.Route{
		Method:  http.MethodGet,
		Handler: d.Home,
		Middlewares: []echo.MiddlewareFunc{
			middleware.AdminOnly(),
		},
	})
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func (d Dashboards) Home(etx *echo.Context) error {
	return d.renderer.Page(etx, "Admin/Dashboard/Home", inertia.Props{}).Render()
}
