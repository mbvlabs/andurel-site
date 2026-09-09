package routes

import (
	"github.com/mbvlabs/andurel/pkg/routing"
)

const DocumentationPrefix = "/docs"

type DocumentationVersionParams struct {
	Version string `param:"version"`
}

type DocumentationShowParams struct {
	Version string `param:"version"`
	Slug    string `param:"slug"`
}

var DocumentationIndex = routing.NewSimpleRoute(
	"",
	"documentations.index",
	DocumentationPrefix,
)

var DocumentationVersion = routing.NewRouteWithParams[DocumentationVersionParams](
	"/:version",
	"documentations.version",
	DocumentationPrefix,
)

var DocumentationShow = routing.NewRouteWithParams[DocumentationShowParams](
	"/:version/:slug",
	"documentations.show",
	DocumentationPrefix,
	routing.InertiaRoute(),
)
