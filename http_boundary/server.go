package http_boundary

import (
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
	"fmt"
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
		s.Router.HandleFunc(r.Path, r.Handler).Methods(r.Methods...)
	}
}

type Route struct {
	Path    string
	Handler http.HandlerFunc
	Methods []string
}

type Response struct {
	Href  *Link   `json:"self"`
	Links []*Link `json:"_links"`
}

type Link struct {
	Name   string `json:"name,omitempty"`
	Href   string `json:"href"`
	Method string `json:"method,omitempty"`
	Rel    string `json:"rel,omitempty"`
}

func Respond(responsePayload interface{}, selfHref string, okStatusCode int, w http.ResponseWriter) {
	bytes, marshalErr := json.Marshal(responsePayload)
	if marshalErr != nil {
		httpErr := &HttpError{Href: &Link{Href: selfHref}, ApiError: &ApiError{marshalErr.Error(), http.StatusInternalServerError}}
		RespondWithError(httpErr, w)
	} else {
		w.WriteHeader(okStatusCode)
		w.Write(bytes)
	}
}

func RespondWithError(err *HttpError, w http.ResponseWriter) {
	// TODO log errors
	if err != nil {
		bytes, jsonErr := json.Marshal(err)
		if jsonErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(err.StatusCode)
			w.Write(bytes)
			return
		}
	}
}

type HttpError struct {
	Href *Link `json:"self"`
	*ApiError
}

type ApiError struct {
	Message    string `json:"message,omitempty"`
	StatusCode int    `json:"statusCode,omitempty"`
}

func (a ApiError) Error() string {
	return fmt.Sprintf("HTTPErr: %d=%s", a.StatusCode, a.Message)
}
