package auth

import (
	"gitlab.com/kskitek/arecar/user-service/http_boundary"
	"net/http"
	"gitlab.com/kskitek/arecar/user-service/user"
)

type Service interface {
	Login(string, string) (string, *http_boundary.ApiError)
}

type service struct {
	userDao       user.Dao
	authenticator Authenticator
}

func (a *service) Login(name string, password string) (string, *http_boundary.ApiError) {
	u, err := a.userDao.GetByName(name)
	if err != nil {
		return "", &http_boundary.ApiError{Message: "Invalid username or password", StatusCode: http.StatusNotFound}
	}
	if u == nil || u.Password != password {
		return "", &http_boundary.ApiError{Message: "Invalid username or password", StatusCode: http.StatusNotFound}
	}

	token, err := a.authenticator.GetToken(u.Name, nil)
	if err != nil {
		return "", &http_boundary.ApiError{Message: err.Error(), StatusCode: http.StatusInternalServerError}
	}
	return token, nil
}
