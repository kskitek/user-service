package user

import (
	"gitlab.com/kskitek/arecar/user-service/http_boundary"
	"net/http"
	"encoding/json"
)

type UserHandler struct {
	uc UseCases
}

type UserResponse struct {
	*User
	*http_boundary.Response
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		uc: &ucs{dao: &MongoDao{}},
	}
}

func (u *UserHandler) Routes() []*http_boundary.Route {
	return []*http_boundary.Route{
		{
			Methods: []string{"POST"},
			Path:    "/user",
			Handler: u.handleUserAdd,
		},
	}
}

func (u *UserHandler) handleUserAdd(w http.ResponseWriter, r *http.Request) {
	var err *http_boundary.ApiError
	user := &User{}

	decodeErr := json.NewDecoder(r.Body).Decode(user)
	if decodeErr != nil {
		err = &http_boundary.ApiError{Message: decodeErr.Error(), StatusCode: http.StatusUnprocessableEntity}
	} else {
		user, err = u.uc.AddUser(user)
	}

	selfHref := r.URL.Path
	if err != nil {
		httpErr := &http_boundary.HttpError{Href: &http_boundary.Link{Href: selfHref}, ApiError: err}
		http_boundary.RespondWithError(httpErr, w)
	} else {
		selfHref += "/" + user.Id
		response := &http_boundary.Response{
			Href: &http_boundary.Link{Href: selfHref},
			Links: []*http_boundary.Link{
				{Name: "delete", Href: selfHref, Method: "DELETE"},
			},
		}
		responsePayload := &UserResponse{User: user, Response: response}
		http_boundary.Respond(responsePayload, r.URL.Path, http.StatusCreated, w)
	}
}
