package routes

import (
	"andurel-site/config"

	"github.com/mbvlabs/andurel/pkg/routing"
)

const AdminDocumentationPagePrefix = "/docs/pages"

var AdminDocumentationPageShow = routing.NewRouteWithBigSerialID(
	"/:id",
	"admin.documentation_pages.show",
	AdminDocumentationPagePrefix,
	routing.InertiaRoute(),
	routing.Host(config.HostAdmin),
)
var AdminDocumentationPageUpdate = routing.NewRouteWithBigSerialID(
	"/:id",
	"admin.documentation_pages.update",
	AdminDocumentationPagePrefix,
	routing.InertiaRoute(),
	routing.Host(config.HostAdmin),
)
var AdminDocumentationPagePublish = routing.NewRouteWithBigSerialID(
	"/:id/publish",
	"admin.documentation_pages.publish",
	AdminDocumentationPagePrefix,
	routing.InertiaRoute(),
	routing.Host(config.HostAdmin),
)
var AdminDocumentationPageRestore = routing.NewRouteWithParams[AdminDocumentationRevisionParams](
	"/:id/revisions/:revisionId/restore",
	"admin.documentation_pages.restore",
	AdminDocumentationPagePrefix,
	routing.InertiaRoute(),
	routing.Host(config.HostAdmin),
)
var AdminDocumentationPagePreview = routing.NewRouteWithBigSerialID(
	"/:id/preview",
	"admin.documentation_pages.preview",
	AdminDocumentationPagePrefix,
	routing.InertiaRoute(),
	routing.Host(config.HostAdmin),
)
