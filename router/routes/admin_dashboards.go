package routes

import (
	"andurel-site/config"

	"github.com/mbvlabs/andurel/pkg/routing"
)

var AdminDashboardHome = routing.NewSimpleRoute(
	"/",
	"admin.dashboards.home",
	"",
	routing.InertiaRoute(),
	routing.Host(config.HostAdmin),
)
