package auth

import (
	"net/http"
	"time"

	"github.com/kskitek/user-service/event"
	"github.com/kskitek/user-service/server"
	"github.com/kskitek/user-service/user"
	"github.com/sirupsen/logrus"
)

const (
	AuthTopic = "user-service.v1.auth"
)

type Service interface {
	Login(string, string) (string, *server.ApiError)
}

type service struct {
	userDao       user.Dao
	authenticator Authenticator
	notifier      event.Notifier
}

func (a *service) Login(name string, password string) (string, *server.ApiError) {
	matching, err := a.userDao.MatchPassword(name, password)
	if err != nil {
		return "", &server.ApiError{Message: "Error when checking password: " + err.Error(), StatusCode: http.StatusInternalServerError}
	}
	if !matching {
		return "", &server.ApiError{Message: "Invalid username or password", StatusCode: http.StatusNotFound}
	}

	token, err := a.authenticator.GetToken(name, nil)
	if err != nil {
		return "", &server.ApiError{Message: err.Error(), StatusCode: http.StatusInternalServerError}
	}
	n := event.Notification{When: time.Now(), Token: token, Payload: name}
	err = a.notifier.Notify(AuthTopic+".login", n)
	if err != nil {
		logrus.WithError(err).Error("unable to notify about login")
	}
	return token, nil
}
