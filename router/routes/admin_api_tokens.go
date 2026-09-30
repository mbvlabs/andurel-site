package routes

import (
	"andurel-site/config"

	"github.com/mbvlabs/andurel/pkg/routing"
)

const AdminApiTokenPrefix = "/tokens"

var AdminApiTokenIndex = routing.NewSimpleRoute(
	"",
	"admin.api_tokens.index",
	AdminApiTokenPrefix,
	routing.InertiaRoute(),
	routing.Host(config.HostAdmin),
)
var AdminApiTokenShow = routing.NewRouteWithUUIDID(
	"/:id",
	"admin.api_tokens.show",
	AdminApiTokenPrefix,
	routing.InertiaRoute(),
	routing.Host(config.HostAdmin),
)
var AdminApiTokenNew = routing.NewSimpleRoute(
	"/new",
	"admin.api_tokens.new",
	AdminApiTokenPrefix,
	routing.InertiaRoute(),
	routing.Host(config.HostAdmin),
)
var AdminApiTokenCreate = routing.NewSimpleRoute(
	"",
	"admin.api_tokens.create",
	AdminApiTokenPrefix,
	routing.InertiaRoute(),
	routing.Host(config.HostAdmin),
)
var AdminApiTokenUpdate = routing.NewRouteWithUUIDID(
	"/:id",
	"admin.api_tokens.update",
	AdminApiTokenPrefix,
	routing.InertiaRoute(),
	routing.Host(config.HostAdmin),
)
var AdminApiTokenDestroy = routing.NewRouteWithUUIDID(
	"/:id",
	"admin.api_tokens.destroy",
	AdminApiTokenPrefix,
	routing.InertiaRoute(),
	routing.Host(config.HostAdmin),
)
