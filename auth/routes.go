package auth

import "gitlab.com/kskitek/arecar/user-service/http_boundary"

func (u *handler) Routes() []*http_boundary.Route {
	return []*http_boundary.Route{
		{
			Methods: []string{"POST"},
			Path:    "/login",
			Handler: u.handleLogin,
		},
	}
}
