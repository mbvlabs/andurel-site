// Package controllers provides HTTP handlers for the web application.
package controllers

import (
	"andurel-site/controllers/api"
	documentation "andurel-site/docs"
	"andurel-site/router"

	"go.uber.org/fx"
)

var otherCache = NewCacheBuilder[string]().WithSize(2).Build

var constructors = fx.Provide(
	otherCache,
	documentation.New,
	NewPages,
	NewDocs,
	NewAssets,
	api.NewAPI,
	NewSessions,
	NewRegistrations,
	NewConfirmations,
	NewResetPasswords,
)

var Module = fx.Module(
	"controllers",
	constructors,
	fx.Invoke(func(r *router.Router, c Pages) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c Docs) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c Assets) error {
		return c.RegisterRoutes(r)
	}),
)
