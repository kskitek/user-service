package auth

import "github.com/kskitek/user-service/server"

func (u *handler) Routes() []*server.Route {
	return []*server.Route{
		{
			Methods: []string{"POST"},
			Path:    "/login",
			Handler: u.handleLogin,
		},
	}
}
