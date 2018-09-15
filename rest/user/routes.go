package user

import "github.com/kskitek/user-service/server"

func (h *handler) Routes() []*server.Route {
	return []*server.Route{
		{
			Methods: []string{"GET"},
			Path:    "/user/{id}",
			Handler: h.handleUserGet,
		},
		{
			Methods: []string{"POST"},
			Path:    "/user",
			Handler: h.handleUserAdd,
		},
		{
			Methods: []string{"DELETE"},
			Path:    "/user/{id}",
			Handler: h.handleUserDelete,
		},
	}
}
