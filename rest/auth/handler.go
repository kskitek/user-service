package auth

import (
	"encoding/json"
	"net/http"

	"github.com/kskitek/user-service/auth"
	"github.com/kskitek/user-service/server"
	"github.com/kskitek/user-service/user"
)

func NewHandler(service auth.Service) server.Handler {
	return &handler{
		s: service,
	}
}

type handler struct {
	s auth.Service
}
type LoginResponse struct {
	Token string `json:"token"`
	*server.Response
}

func (a *handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var err *server.ApiError
	u := &user.User{}
	selfHref := r.URL.Path
	var token string

	decodeErr := json.NewDecoder(r.Body).Decode(u)
	if decodeErr != nil {
		err = &server.ApiError{Message: decodeErr.Error(), StatusCode: http.StatusUnprocessableEntity}
	} else {
		token, err = a.s.Login(u.Name, u.Password)
	}

	if err != nil {
		httpErr := &server.HttpError{Href: &server.Link{Href: selfHref}, ApiError: err}
		server.RespondWithError(httpErr, w)
	} else {
		selfHref += "/" + u.Id
		response := authLink(selfHref)
		responsePayload := &LoginResponse{Token: token, Response: response}
		server.Respond(responsePayload, r.URL.Path, http.StatusCreated, w)
	}
}

func authLink(selfHref string) *server.Response {
	response := &server.Response{
		Href: &server.Link{Href: selfHref},
		Links: []*server.Link{
			{Name: "logout", Href: selfHref, Method: "POST"},
		},
	}
	return response
}
