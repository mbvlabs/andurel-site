// Package controllers provides HTTP handlers for the web application.
package controllers

import (
	"andurel-site/config"
	"andurel-site/controllers/admin"
	"andurel-site/controllers/api"
	"andurel-site/docs"
	"andurel-site/router"
	"andurel-site/views"

	"go.uber.org/fx"
)

func configureViews(cfg config.App) {
	views.ConfigureHead(cfg.ProjectName, cfg.BaseURL())
}

var otherCache = NewCacheBuilder[string]().WithSize(2).Build

var constructors = fx.Provide(
	otherCache,
	NewPages,
	NewAssets,
	api.NewAPI,
	NewSessions,
	NewRegistrations,
	NewConfirmations,
	NewResetPasswords,
	docs.New,
	NewDocumentations,
	admin.NewDashboards,
	admin.NewDocumentationVersions,
	admin.NewDocumentationPages,
	admin.NewApiTokens,
	api.NewDocumentations,
)

var Module = fx.Module(
	"controllers",
	constructors,
	fx.Invoke(configureViews),
	fx.Invoke(func(r *router.Router, c Pages) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c Assets) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c api.API) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c Sessions) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c Registrations) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c Confirmations) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c ResetPasswords) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c Documentations) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c admin.Dashboards) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c admin.DocumentationVersions) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c admin.DocumentationPages) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c admin.ApiTokens) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c api.Documentations) error {
		return c.RegisterRoutes(r)
	}),
)
