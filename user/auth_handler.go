package user

import (
	"gitlab.com/kskitek/arecar/user-service/http_boundary"
	"net/http"
	"encoding/json"
)

type LoginResponse struct {
	Token string `json:"token"`
	*http_boundary.Response
}

// TODO create some kind of factory to provide daos
func (a *UserHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var err *http_boundary.ApiError
	user := &User{}
	selfHref := r.URL.Path
	var token string

	decodeErr := json.NewDecoder(r.Body).Decode(user)
	if decodeErr != nil {
		err = &http_boundary.ApiError{Message: decodeErr.Error(), StatusCode: http.StatusUnprocessableEntity}
	} else {
		token, err = a.auth.Login(user.Name, user.Password)
	}

	if err != nil {
		httpErr := &http_boundary.HttpError{Href: &http_boundary.Link{Href: selfHref}, ApiError: err}
		http_boundary.RespondWithError(httpErr, w)
	} else {
		selfHref += "/" + user.Id
		response := authLink(selfHref)
		responsePayload := &LoginResponse{Token: token, Response: response}
		http_boundary.Respond(responsePayload, r.URL.Path, http.StatusCreated, w)
	}
}

func authLink(selfHref string) *http_boundary.Response {
	response := &http_boundary.Response{
		Href: &http_boundary.Link{Href: selfHref},
		Links: []*http_boundary.Link{
			{Name: "logout", Href: selfHref, Method: "POST"},
		},
	}
	return response
}
