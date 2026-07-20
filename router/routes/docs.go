package routes

import "andurel-site/internal/routing"

type DocsVersionParams struct {
	Version string `param:"version"`
}

type DocsShowParams struct {
	Version string `param:"version"`
	Slug    string `param:"slug"`
}

var DocsIndex = routing.NewSimpleRoute(
	"/docs",
	"docs.index",
	"",
)

var DocsSearch = routing.NewSimpleRoute(
	"/docs/search.json",
	"docs.search",
	"",
)

var DocsVersion = routing.NewRouteWithParams[DocsVersionParams](
	"/docs/:version",
	"docs.version",
	"",
)

var DocsShow = routing.NewRouteWithParams[DocsShowParams](
	"/docs/:version/:slug",
	"docs.show",
	"",
)
