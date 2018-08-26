package auth

import "github.com/kskitek/user-service/http_boundary"

func (u *handler) Routes() []*http_boundary.Route {
	return []*http_boundary.Route{
		{
			Methods: []string{"POST"},
			Path:    "/login",
			Handler: u.handleLogin,
		},
	}
}
