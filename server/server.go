package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

var Server = &server{router: mux.NewRouter()}

type server struct {
	router *mux.Router
}

type Route struct {
	Path    string
	Handler http.HandlerFunc
	Methods []string
}

func (s *server) Handler() http.Handler {
	return nethttp.Middleware(opentracing.GlobalTracer(), s.router)
}

func (s *server) AddRoutes(routes []*Route) {
	for _, r := range routes {
		s.router.HandleFunc(r.Path, handleAround(r.Handler)).Methods(r.Methods...)
	}
}

func handleAround(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()
		defer func() {
			logrus.WithFields(logrus.Fields{"m": r.Method, "p": r.URL.Path, "t": time.Since(begin)}).Info("Request")
		}()

		h.ServeHTTP(w, r)
	}
}
