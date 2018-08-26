package auth

import (
	"net/http"

	"github.com/kskitek/user-service/http_boundary"
	"github.com/kskitek/user-service/user"
)

type Service interface {
	Login(string, string) (string, *http_boundary.ApiError)
}

type service struct {
	userDao       user.Dao
	authenticator Authenticator
}

func (a *service) Login(name string, password string) (string, *http_boundary.ApiError) {
	matching, err := a.userDao.MatchPassword(name, password)
	if err != nil {
		return "", &http_boundary.ApiError{Message: "Error when checking password: " + err.Error(), StatusCode: http.StatusInternalServerError}
	}
	if !matching {
		return "", &http_boundary.ApiError{Message: "Invalid username or password", StatusCode: http.StatusNotFound}
	}

	token, err := a.authenticator.GetToken(name, nil)
	if err != nil {
		return "", &http_boundary.ApiError{Message: err.Error(), StatusCode: http.StatusInternalServerError}
	}
	return token, nil
}
