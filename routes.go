package main

import (
	"github.com/kskitek/user-service/auth"
	"github.com/kskitek/user-service/http_boundary"
	"github.com/kskitek/user-service/user"
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
