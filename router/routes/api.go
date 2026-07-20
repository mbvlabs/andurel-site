package routes

import (
	"andurel-site/internal/routing"
)

const APIPrefix = "/api"

var Health = routing.NewSimpleRoute(
	"/health",
	"api.health",
	APIPrefix,
)
