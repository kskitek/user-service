package user

import "gitlab.com/kskitek/arecar/user-service/http_boundary"

func (u *UserHandler) Routes() []*http_boundary.Route {
	return []*http_boundary.Route{
		{
			Methods: []string{"GET"},
			Path:    "/user/{id}",
			Handler: u.handleUserGet,
		},
		{
			Methods: []string{"POST"},
			Path:    "/user",
			Handler: u.handleUserAdd,
		},
	}
}
