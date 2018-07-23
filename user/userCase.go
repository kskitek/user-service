package user

import (
	"time"
	"net/http"
	"fmt"
	"gitlab.com/kskitek/arecar/user-service/http_boundary"
)

type User struct {
	Id               string    `json:"id,omitempty"`
	Name             string    `json:"name"`
	Email            string    `json:"email"`
	Password         string    `json:"password,omitempty"`
	RegistrationDate time.Time `json:"registrationDate,omitempty"`
}

type UseCases interface {
	AddUser(*User) (*User, *http_boundary.ApiError)
}

type ucs struct {
	dao Dao
}

func (uc *ucs) AddUser(user *User) (*User, *http_boundary.ApiError) {
	if user == nil {
		return nil, &http_boundary.ApiError{"User details required", http.StatusUnprocessableEntity}
	}
	exists, err := uc.dao.UserExists(user)
	if err != nil {
		return nil, &http_boundary.ApiError{}
	}
	if exists {
		return nil, &http_boundary.ApiError{"User already exists.", http.StatusConflict}
	}

	apiErr := validateAddUserPayload(user)
	if apiErr != nil {
		return nil, apiErr
	}
	newUser, err := uc.dao.AddUser(user)
	if err != nil {
		fmt.Println(err)
		return nil, &http_boundary.ApiError{Message: "Cannot add user", StatusCode: http.StatusUnprocessableEntity}
	}

	return newUser, nil
}
func validateAddUserPayload(user *User) *http_boundary.ApiError {
	if !validateEmail(user.Email) {
		return &http_boundary.ApiError{Message: "Invalid email address", StatusCode: http.StatusUnprocessableEntity}
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
