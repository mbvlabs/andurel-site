package routes

import (
	"andurel-site/config"

	"github.com/mbvlabs/andurel/pkg/routing"
)

const AdminDocumentationVersionPrefix = "/docs"

type AdminDocumentationRevisionParams struct {
	ID         string `param:"id"`
	RevisionID string `param:"revisionId"`
}

var AdminDocumentationVersionIndex = routing.NewSimpleRoute(
	"",
	"admin.documentation_versions.index",
	AdminDocumentationVersionPrefix,
	routing.InertiaRoute(),
	routing.Host(config.HostAdmin),
)
var AdminDocumentationVersionNew = routing.NewSimpleRoute(
	"/new",
	"admin.documentation_versions.new",
	AdminDocumentationVersionPrefix,
	routing.InertiaRoute(),
	routing.Host(config.HostAdmin),
)
var AdminDocumentationVersionCreate = routing.NewSimpleRoute(
	"",
	"admin.documentation_versions.create",
	AdminDocumentationVersionPrefix,
	routing.InertiaRoute(),
	routing.Host(config.HostAdmin),
)
var AdminDocumentationVersionShow = routing.NewRouteWithBigSerialID(
	"/:id",
	"admin.documentation_versions.show",
	AdminDocumentationVersionPrefix,
	routing.InertiaRoute(),
	routing.Host(config.HostAdmin),
)
var AdminDocumentationVersionUpdate = routing.NewRouteWithBigSerialID(
	"/:id",
	"admin.documentation_versions.update",
	AdminDocumentationVersionPrefix,
	routing.InertiaRoute(),
	routing.Host(config.HostAdmin),
)
var AdminDocumentationVersionCreatePage = routing.NewRouteWithBigSerialID(
	"/:id/pages",
	"admin.documentation_versions.pages.create",
	AdminDocumentationVersionPrefix,
	routing.InertiaRoute(),
	routing.Host(config.HostAdmin),
)
var AdminDocumentationVersionUpdateNav = routing.NewRouteWithBigSerialID(
	"/:id/nav",
	"admin.documentation_versions.nav.update",
	AdminDocumentationVersionPrefix,
	routing.InertiaRoute(),
	routing.Host(config.HostAdmin),
)
var AdminDocumentationVersionPublishNav = routing.NewRouteWithBigSerialID(
	"/:id/nav/publish",
	"admin.documentation_versions.nav.publish",
	AdminDocumentationVersionPrefix,
	routing.InertiaRoute(),
	routing.Host(config.HostAdmin),
)
