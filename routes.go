package main

import (
	"gitlab.com/kskitek/arecar/user-service/http_boundary"
	"gitlab.com/kskitek/arecar/user-service/user"
	"gitlab.com/kskitek/arecar/user-service/auth"
)

var handlers = []http_boundary.Handler{
	user.NewHandler(),
	auth.NewHandler(),
}

func getRoutes() []*http_boundary.Route {
	var routes []*http_boundary.Route

	for _, h := range handlers {
		routes = append(routes, h.Routes()...)
	}

	return routes
}
