package user

import (
	auth_token "gitlab.com/kskitek/arecar/user-service/auth"
	"gitlab.com/kskitek/arecar/user-service/http_boundary"
	"net/http"
)

type Auth interface {
	Login(string, string) (string, *http_boundary.ApiError)
}

type auth struct {
	userDao       Dao
	authenticator auth_token.Authenticator
}

func (a *auth) Login(name string, password string) (string, *http_boundary.ApiError) {
	user, err := a.userDao.GetByName(name)
	if err != nil {
		return "", &http_boundary.ApiError{Message: "Invalid username or password", StatusCode: http.StatusNotFound}
	}
	if user == nil || user.Password != password {
		return "", &http_boundary.ApiError{Message: "Invalid username or password", StatusCode: http.StatusNotFound}
	}

	token, err := a.authenticator.GetToken(user.Name, nil)
	if err != nil {
		return "", &http_boundary.ApiError{Message: err.Error(), StatusCode: http.StatusInternalServerError}
	}
	return token, nil
}
