package user

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

var userOkId = int64(1)
var userOk = &User{
	Id:               "1",
	Name:             "User1",
	Password:         "Pwd",
	Email:            "user1@gmail.com",
	RegistrationDate: time.Now(),
}

func newOut() Service {
	dao := NewDao()
	dao.Add(userOk)
	return &crud{dao}
}

func Test_WhenGetUserGivenUserHasPasswordThenReturnedPasswordIsEmpty(t *testing.T) {
	out := newOut()

	user, apiError := out.GetUser(userOkId)

	assert.Nil(t, apiError)

	assert.Equal(t, "", user.Password)
}
