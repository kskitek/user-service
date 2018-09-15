package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kskitek/user-service/server"
	"github.com/kskitek/user-service/user"
)

func NewHandler(service user.Service) server.Handler {
	return &handler{
		s: service,
	}
}

type handler struct {
	s user.Service
}

type Response struct {
	*user.User
	*server.Response
}

func (h *handler) handleUserGet(w http.ResponseWriter, r *http.Request) {
	var err *server.ApiError
	id := mux.Vars(r)["id"]
	selfHref := r.URL.Path

	intId, parseErr := stringToId(id)
	if parseErr != nil {
		httpErr := &server.HttpError{Href: &server.Link{Href: selfHref}, ApiError: parseErr}
		server.RespondWithError(httpErr, w)
	}
	u, err := h.s.Get(intId)

	if err != nil {
		httpErr := &server.HttpError{Href: &server.Link{Href: selfHref}, ApiError: err}
		server.RespondWithError(httpErr, w)
	} else {
		response := userLink(selfHref)
		responsePayload := &Response{User: u, Response: response}
		server.Respond(responsePayload, r.URL.Path, http.StatusOK, w)
	}
}

func stringToId(id string) (int64, *server.ApiError) {
	intId, parseErr := strconv.ParseInt(id, 10, 64)
	if parseErr != nil {
		return 0, &server.ApiError{Message: fmt.Sprintf("cannot parse '%s' as user id", id), StatusCode: http.StatusBadRequest}
	}
	return intId, nil
}

func (h *handler) handleUserAdd(w http.ResponseWriter, r *http.Request) {
	var err *server.ApiError
	u := &user.User{}
	selfHref := r.URL.Path

	decodeErr := json.NewDecoder(r.Body).Decode(u)
	if decodeErr != nil {
		err = &server.ApiError{Message: decodeErr.Error(), StatusCode: http.StatusUnprocessableEntity}
	} else {
		u, err = h.s.Add(u)
	}

	if err != nil {
		httpErr := &server.HttpError{Href: &server.Link{Href: selfHref}, ApiError: err}
		server.RespondWithError(httpErr, w)
	} else {
		selfHref += "/" + u.Id
		response := userLink(selfHref)
		responsePayload := &Response{User: u, Response: response}
		server.Respond(responsePayload, r.URL.Path, http.StatusCreated, w)
	}
}

func (h *handler) handleUserDelete(w http.ResponseWriter, r *http.Request) {
	var err *server.ApiError
	id := mux.Vars(r)["id"]
	selfHref := r.URL.Path

	intId, parseErr := stringToId(id)
	if parseErr != nil {
		httpErr := &server.HttpError{Href: &server.Link{Href: selfHref}, ApiError: parseErr}
		server.RespondWithError(httpErr, w)
	}
	err = h.s.Delete(intId)

	if err != nil {
		httpErr := &server.HttpError{Href: &server.Link{Href: selfHref}, ApiError: err}
		server.RespondWithError(httpErr, w)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}

}

func userLink(selfHref string) *server.Response {
	response := &server.Response{
		Href: &server.Link{Href: selfHref},
		Links: []*server.Link{
			{Name: "delete", Href: selfHref, Method: "DELETE"},
		},
	}
	return response
}
