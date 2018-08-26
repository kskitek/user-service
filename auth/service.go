package auth

import (
	"net/http"
	"time"

	"github.com/kskitek/user-service/event"
	"github.com/kskitek/user-service/http_boundary"
	"github.com/kskitek/user-service/user"
)

const (
	AuthTopic = "user-service.v1.auth"
)

type Service interface {
	Login(string, string) (string, *http_boundary.ApiError)
}

type service struct {
	userDao       user.Dao
	authenticator Authenticator
	notifier      event.Notifier
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
	n := event.Notification{When: time.Now(), Token: token, Payload: name}
	a.notifier.Notify(AuthTopic+".login", n)
	return token, nil
}
