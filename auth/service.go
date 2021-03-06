package auth

import (
	"context"
	"net/http"
	"regexp"
	"time"

	"github.com/kskitek/user-service/event"
	"github.com/kskitek/user-service/server"
	"github.com/kskitek/user-service/tracing"
	"github.com/kskitek/user-service/user"
	"github.com/sirupsen/logrus"
)

const (
	AuthTopic = "user-service.v1.auth"
)

type Service interface {
	Login(context.Context, string, string) (string, *server.ApiError)
}

type service struct {
	userDao       user.Dao
	authenticator Authenticator
	notifier      event.Notifier
}

func (a *service) Login(ctx context.Context, name string, password string) (string, *server.ApiError) {
	spanFinish := tracing.SetUpTraceWithTags(ctx, "dao", tracing.Tags{"username": name, "op": "matchPassword"})
	if !validateName(name) {
		msg := "invalid username"
		logrus.WithField("username", name).Error(msg)
		return "", &server.ApiError{Message: "Invalid username or password", StatusCode: http.StatusNotAcceptable}
	}
	matching, err := a.userDao.MatchPassword(name, password)
	if err != nil {
		spanFinish()
		return "", &server.ApiError{Message: "Error when checking password: " + err.Error(), StatusCode: http.StatusInternalServerError}
	}
	if !matching {
		spanFinish()
		return "", &server.ApiError{Message: "Invalid username or password", StatusCode: http.StatusNotFound}
	}
	spanFinish()

	spanFinish = tracing.SetUpTraceWithTags(ctx, "getAuthToken", tracing.Tags{"username": name})
	token, err := a.authenticator.GetToken(name, nil)
	if err != nil {
		spanFinish()
		return "", &server.ApiError{Message: err.Error(), StatusCode: http.StatusInternalServerError}
	}
	spanFinish()

	n := event.Notification{CorrelationId: tracing.ContextToString(ctx), When: time.Now(), Token: token, Payload: name}
	err = a.notifier.Notify(AuthTopic+".login", n)
	if err != nil {
		logrus.WithError(err).Error("unable to notify about login")
	}
	return token, nil
}

const (
	namePattern     = "[a-zA-Z0-9\\-\\_\\.\\+]{5,20}"
	userNamePattern = "^" + namePattern + "$"
)

var (
	userNameRegexp = regexp.MustCompile(userNamePattern)
)

func validateName(name string) bool {
	return userNameRegexp.MatchString(name)
}
