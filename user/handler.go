package user

import (
	"gitlab.com/kskitek/arecar/user-service/http_boundary"
	"net/http"
)

func Routes() []*http_boundary.Route {
	return []*http_boundary.Route{
		{
			Methods: []string{"GET"},
			Path:    "/user/{id}",
			Handler: handleUserAdd,
		},
	}
}

func handleUserAdd(_ *http_boundary.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}
