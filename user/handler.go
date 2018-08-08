package user

import (
	"gitlab.com/kskitek/arecar/user-service/http_boundary"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"strconv"
)

type handler struct {
	s Service
}

type UserResponse struct {
	*User
	*http_boundary.Response
}

func (u *handler) handleUserGet(w http.ResponseWriter, r *http.Request) {
	var err *http_boundary.ApiError
	id := mux.Vars(r)["id"]
	selfHref := r.URL.Path

	intId, parseErr := strconv.ParseInt(id, 10, 64)
	if parseErr != nil {
		httpErr := &http_boundary.HttpError{Href: &http_boundary.Link{Href: selfHref}, ApiError: err}
		http_boundary.RespondWithError(httpErr, w)
	}
	user, err := u.s.GetUser(intId)

	if err != nil {
		httpErr := &http_boundary.HttpError{Href: &http_boundary.Link{Href: selfHref}, ApiError: err}
		http_boundary.RespondWithError(httpErr, w)
	} else {
		response := userLink(selfHref)
		responsePayload := &UserResponse{User: user, Response: response}
		http_boundary.Respond(responsePayload, r.URL.Path, http.StatusOK, w)
	}
}

func (u *handler) handleUserAdd(w http.ResponseWriter, r *http.Request) {
	var err *http_boundary.ApiError
	user := &User{}
	selfHref := r.URL.Path

	decodeErr := json.NewDecoder(r.Body).Decode(user)
	if decodeErr != nil {
		err = &http_boundary.ApiError{Message: decodeErr.Error(), StatusCode: http.StatusUnprocessableEntity}
	} else {
		user, err = u.s.AddUser(user)
	}

	if err != nil {
		httpErr := &http_boundary.HttpError{Href: &http_boundary.Link{Href: selfHref}, ApiError: err}
		http_boundary.RespondWithError(httpErr, w)
	} else {
		selfHref += "/" + user.Id
		response := userLink(selfHref)
		responsePayload := &UserResponse{User: user, Response: response}
		http_boundary.Respond(responsePayload, r.URL.Path, http.StatusCreated, w)
	}
}

func (u *handler) handleUserDelete(w http.ResponseWriter, r *http.Request) {
	var err *http_boundary.ApiError
	id := mux.Vars(r)["id"]
	selfHref := r.URL.Path

	intId, parseErr := strconv.ParseInt(id, 10, 64)
	if parseErr != nil {
		httpErr := &http_boundary.HttpError{Href: &http_boundary.Link{Href: selfHref}, ApiError: err}
		http_boundary.RespondWithError(httpErr, w)
	}
	err = u.s.DeleteUser(intId)

	if err != nil {
		httpErr := &http_boundary.HttpError{Href: &http_boundary.Link{Href: selfHref}, ApiError: err}
		http_boundary.RespondWithError(httpErr, w)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}

}

func userLink(selfHref string) *http_boundary.Response {
	response := &http_boundary.Response{
		Href: &http_boundary.Link{Href: selfHref},
		Links: []*http_boundary.Link{
			{Name: "delete", Href: selfHref, Method: "DELETE"},
		},
	}
	return response
}
