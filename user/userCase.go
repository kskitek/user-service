package user

import "time"

type User struct {
	Id               string    `json:"id,omitempty"`
	Name             string    `json:"name"`
	Email            string    `json:"email"`
	Password         string    `json:"password"`
	RegistrationDate time.Time `json:"registrationDate, omitempty"`
}

type UseCase interface {
	AddUser(*User) (*User, error)
}

type uc struct {
	dao dao
}

func (uc *uc) AddUser(user *User) (*User, error) {
	return &User{}, nil
}
