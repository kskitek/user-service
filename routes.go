package main

import (
	"gitlab.com/kskitek/arecar/user-service/http_boundary"
	"gitlab.com/kskitek/arecar/user-service/user"
)

func getRoutes() []*http_boundary.Route {
	var routes []*http_boundary.Route

	routes = append(routes, user.Routes()...)

	return routes
}
