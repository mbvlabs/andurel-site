package routes

import (
	"github.com/mbvlabs/andurel/pkg/routing"
)

type APIDocumentationParams struct {
	Slug string `param:"slug"`
}

var APIDocumentationIndex = routing.NewSimpleRoute(
	"/documentations",
	"api.documentations.index",
	APIPrefix,
)

var APIDocumentationCreate = routing.NewSimpleRoute(
	"/documentations",
	"api.documentations.create",
	APIPrefix,
)

var APIDocumentationShow = routing.NewRouteWithParams[APIDocumentationParams](
	"/documentations/:slug",
	"api.documentations.show",
	APIPrefix,
)

var APIDocumentationUpdate = routing.NewRouteWithParams[APIDocumentationParams](
	"/documentations/:slug",
	"api.documentations.update",
	APIPrefix,
)
