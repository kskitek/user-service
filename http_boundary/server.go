package http_boundary

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	Router *mux.Router
}

func NewServer(routes []*Route) *Server {
	s := &Server{Router: mux.NewRouter()}
	s.addRoutes(routes)
	return s
}

func (s *Server) addRoutes(routes []*Route) {
	for _, r := range routes {
		httpHandler := r.Handler(s)
		s.Router.HandleFunc(r.Path, httpHandler).Methods(r.Methods...)
	}
}

type Route struct {
	Path    string
	Handler ServerHandlerFunc
	Methods []string
}

type ServerHandlerFunc = func(*Server) http.HandlerFunc
