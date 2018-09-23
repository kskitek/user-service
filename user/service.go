package user

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/kskitek/user-service/event"
	"github.com/kskitek/user-service/server"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
)

const (
	CrudBaseTopic = "user-service.v1.crud"
)

type User struct {
	Id               string    `json:"id,omitempty"`
	Name             string    `json:"name"`
	Email            string    `json:"email"`
	Password         string    `json:"password,omitempty"`
	RegistrationDate time.Time `json:"registrationDate,omitempty"`
}

type Service interface {
	Get(context.Context, int64) (*User, *server.ApiError)
	Add(context.Context, *User) (*User, *server.ApiError)
	Delete(context.Context, int64) *server.ApiError
}

func NewService(dao Dao, notifier event.Notifier) Service {
	return &crud{
		dao:      dao,
		notifier: notifier,
	}
}

type crud struct {
	dao      Dao
	notifier event.Notifier
}

func (uc *crud) Get(ctx context.Context, id int64) (*User, *server.ApiError) {
	if id <= 0 {
		return nil, &server.ApiError{Message: "Id required", StatusCode: http.StatusBadRequest}
	}
	user, err := uc.dao.GetById(ctx, id)
	if err != nil {
		return nil, &server.ApiError{Message: "Cannot read user: " + err.Error(), StatusCode: http.StatusInternalServerError}
	}
	if user == nil {
		return nil, &server.ApiError{Message: "User not found", StatusCode: http.StatusNotFound}
	}

	user.Password = ""
	return user, nil
}

func (uc *crud) Add(ctx context.Context, user *User) (*User, *server.ApiError) {
	if user == nil {
		return nil, &server.ApiError{Message: "User details required", StatusCode: http.StatusUnprocessableEntity}
	}
	err := uc.checkIfExists(ctx, user)
	if err != nil {
		return nil, err
	}

	apiErr := validateAddUserPayload(user)
	if apiErr != nil {
		return nil, apiErr
	}

	tags := tags{"op": "add"}
	newUser, err := uc.add(ctx, user, tags)
	if err != nil {
		return nil, err
	}

	newUser.Password = ""
	n := event.Notification{CorrelationId: contextToString(ctx), Payload: newUser}
	uc.notify(ctx, CrudBaseTopic+".add", n, tags)
	return newUser, nil
}

func (uc *crud) Delete(ctx context.Context, id int64) *server.ApiError {
	if id == 0 {
		return &server.ApiError{Message: "Id required", StatusCode: http.StatusBadRequest}
	}
	tags := tags{"op": "delete"}
	err := uc.delete(ctx, id, tags)
	if err != nil {
		return err
	}
	n := event.Notification{CorrelationId: contextToString(ctx), Payload: id}
	uc.notify(ctx, CrudBaseTopic+".delete", n, tags)

	return nil
}

func (uc *crud) checkIfExists(ctx context.Context, user *User) *server.ApiError {
	defer setUpTraceWithTags(ctx, "dao", tags{"op": "exists"})()
	exists, err := uc.dao.Exists(user)
	if err != nil {
		return &server.ApiError{Message: "Cannot save user: " + err.Error(), StatusCode: http.StatusInternalServerError}
	}
	if exists {
		return &server.ApiError{Message: "User already exists.", StatusCode: http.StatusConflict}
	}
	return nil
}

func (uc *crud) add(ctx context.Context, user *User, tags tags) (*User, *server.ApiError) {
	defer setUpTraceWithTags(ctx, "dao", tags)()
	newUser, err := uc.dao.Add(user)
	if err != nil {
		fmt.Println(err)
		return nil, &server.ApiError{Message: "Cannot add user", StatusCode: http.StatusUnprocessableEntity}
	}
	return newUser, nil
}

func (uc *crud) notify(ctx context.Context, topic string, n event.Notification, tags tags) {
	defer setUpTraceWithTags(ctx, "notification", tags)()
	err := uc.notifier.Notify(topic, n)
	if err != nil {
		logrus.WithError(err).
			WithFields(logrus.Fields{"notification": n, "topic": topic}).
			Error("error when notifying")
	}
}
func (uc *crud) delete(ctx context.Context, id int64, tags tags) *server.ApiError {
	defer setUpTraceWithTags(ctx, "dao", tags)()
	err := uc.dao.Delete(id)
	if err != nil {
		return &server.ApiError{Message: "Cannot delete user: " + err.Error(), StatusCode: http.StatusInternalServerError}
	}
	return nil
}

func validateAddUserPayload(user *User) *server.ApiError {
	if !validateEmail(user.Email) {
		return &server.ApiError{Message: "Invalid email address", StatusCode: http.StatusUnprocessableEntity}
	}
	if user.Name == "" {
		return &server.ApiError{Message: "Name cannot be empty", StatusCode: http.StatusUnprocessableEntity}
	}
	if user.Password == "" {
		return &server.ApiError{Message: "Password cannot be empty", StatusCode: http.StatusUnprocessableEntity}
	}

	return nil
}

func validateEmail(email string) bool {
	// TODO email pattern
	if email == "" {
		return false
	} else {
		return true
	}
}

type tags map[string]string

func setUpTraceWithTags(ctx context.Context, opName string, tags map[string]string) (deferFunc func()) {
	span, _ := opentracing.StartSpanFromContext(ctx, opName)
	for k, v := range tags {
		span.SetTag(k, v)
	}
	return func() {
		span.Finish()
	}
}

func contextToString(ctx context.Context) string {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return ""
	}
	j, ok :=span.Context().(jaeger.SpanContext)
	if !ok {
		return ""
	}
	return j.String()
}
