package user

import (
	"time"
	"net/http"
	"fmt"
	"gitlab.com/kskitek/arecar/user-service/http_boundary"
	"gitlab.com/kskitek/arecar/user-service/event"
	"github.com/sirupsen/logrus"
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
	Get(int64) (*User, *http_boundary.ApiError)
	Add(*User) (*User, *http_boundary.ApiError)
	Delete(int64) *http_boundary.ApiError
}

type crud struct {
	dao      Dao
	notifier event.Notifier
}

func (uc *crud) Get(id int64) (*User, *http_boundary.ApiError) {
	if id <= 0 {
		return nil, &http_boundary.ApiError{Message: "Id required", StatusCode: http.StatusBadRequest}
	}
	user, err := uc.dao.GetById(id)
	if err != nil {
		return nil, &http_boundary.ApiError{Message: "Cannot read user: " + err.Error(), StatusCode: http.StatusInternalServerError}
	}
	if user == nil {
		return nil, &http_boundary.ApiError{Message: "User not found", StatusCode: http.StatusNotFound}
	}

	user.Password = ""
	return user, nil
}

func (uc *crud) Add(user *User) (*User, *http_boundary.ApiError) {
	if user == nil {
		return nil, &http_boundary.ApiError{Message: "User details required", StatusCode: http.StatusUnprocessableEntity}
	}
	exists, err := uc.dao.Exists(user)
	if err != nil {
		return nil, &http_boundary.ApiError{Message: "Cannot save user: " + err.Error(), StatusCode: http.StatusInternalServerError}
	}
	if exists {
		return nil, &http_boundary.ApiError{Message: "User already exists.", StatusCode: http.StatusConflict}
	}

	apiErr := validateAddUserPayload(user)
	if apiErr != nil {
		return nil, apiErr
	}
	newUser, err := uc.dao.Add(user)
	if err != nil {
		fmt.Println(err)
		return nil, &http_boundary.ApiError{Message: "Cannot add user", StatusCode: http.StatusUnprocessableEntity}
	}

	newUser.Password = ""
	n := event.Notification{Payload: newUser}
	err = uc.notifier.Notify(CrudBaseTopic+".add", n)
	if err != nil {
		logrus.WithError(err).WithField("notification", n).Error("error when notifying about new user")
	}
	return newUser, nil
}

func (uc *crud) Delete(id int64) *http_boundary.ApiError {
	if id == 0 {
		return &http_boundary.ApiError{Message: "Id required", StatusCode: http.StatusBadRequest}
	}
	err := uc.dao.Delete(id)
	if err != nil {
		return &http_boundary.ApiError{Message: "Cannot delete user: " + err.Error(), StatusCode: http.StatusInternalServerError}
	}
	n := event.Notification{Payload: id}
	err = uc.notifier.Notify(CrudBaseTopic+".delete", n)
	if err != nil {
		logrus.WithError(err).WithField("notification", n).Error("error when notifying about deleted user")
	}

	return nil
}

func validateAddUserPayload(user *User) *http_boundary.ApiError {
	if !validateEmail(user.Email) {
		return &http_boundary.ApiError{Message: "Invalid email address", StatusCode: http.StatusUnprocessableEntity}
	}
	if user.Name == "" {
		return &http_boundary.ApiError{Message: "Name cannot be empty", StatusCode: http.StatusUnprocessableEntity}
	}
	if user.Password == "" {
		return &http_boundary.ApiError{Message: "Password cannot be empty", StatusCode: http.StatusUnprocessableEntity}
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
