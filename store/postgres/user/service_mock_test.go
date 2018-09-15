package user

import (
	"strconv"
	"time"

	"github.com/kskitek/user-service/user"
)

var testSuffix = strconv.FormatInt(time.Now().Unix(), 10)

func newUser(id string, name string) *user.User {
	return &user.User{
		Id:               id,
		Name:             name + testSuffix,
		Password:         "Pwd",
		Email:            name + testSuffix + "@gmail.com",
		RegistrationDate: time.Now().UTC(),
	}
}

func UserOk() *user.User {
	return newUser("-1", "User1_")
}

func UserOk2() *user.User {
	return newUser("-1", "User2_")
}
