package http_boundary

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Router *mux.Router
}

func NewServer(routes []*Route) *Server {
	s := &Server{Router: mux.NewRouter()}
	s.addRoutes(routes)
	return s
}

type Route struct {
	Path    string
	Handler http.HandlerFunc
	Methods []string
}

func (s *Server) addRoutes(routes []*Route) {
	for _, r := range routes {
		s.Router.HandleFunc(r.Path, handleAround(r.Handler)).Methods(r.Methods...)
	}
}

func handleAround(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()
		h.ServeHTTP(w, r)
		logrus.WithField("m", r.Method).WithField("p", r.URL.Path).WithField("t", time.Since(begin)).Info("Request")
	}
}
