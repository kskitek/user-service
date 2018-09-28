// +build it

package user

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/kskitek/user-service/user"
	"github.com/stretchr/testify/assert"
)

var out = NewPgDao()

func Test_Add_OkUser_ReturnsUserWithNewId(t *testing.T) {
	user1, err1 := out.Add(UserOk())
	user2, err2 := out.Add(UserOk2())

	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.NotEqual(t, user1.Id, UserOk().Id)
	assert.NotEqual(t, user2.Id, UserOk2().Id)
	assert.NotEqual(t, user1.Id, user2.Id)
}

func Test_Add_TheSameUserTwoTimes_Error(t *testing.T) {
	u := getTestUser(t)
	_, err1 := out.Add(u)
	_, err2 := out.Add(u)

	assert.Nil(t, err1)
	assert.NotNil(t, err2)
}

func Test_AddAndGet_ReturnsTheUser(t *testing.T) {
	u := getTestUser(t)

	newUser, err := out.Add(u)
	assert.Nil(t, err)
	newId, _ := strconv.ParseInt(newUser.Id, 10, 64)

	userByName, err := out.GetById(context.TODO(), newId)
	assert.Nil(t, err)

	assert.Equal(t, newUser, userByName)
}

func Test_Add_GetByNameReturnsTheUser(t *testing.T) {
	u := getTestUser(t)

	newUser, err := out.Add(u)
	assert.Nil(t, err)

	userByName, err := out.GetByName(u.Name)
	assert.Nil(t, err)

	assert.Equal(t, newUser, userByName)
}

func Test_Add_ExistsReturnsTrue(t *testing.T) {
	u := getTestUser(t)

	newUser, err := out.Add(u)
	assert.Nil(t, err)

	userByName, err := out.Exists(newUser)
	assert.Nil(t, err)

	assert.True(t, userByName)
}

func Test_AfterDelete_CannotGetUser(t *testing.T) {
	u := getTestUser(t)

	newUser, err := out.Add(u)
	assert.Nil(t, err)

	newId, _ := strconv.ParseInt(newUser.Id, 10, 64)
	err = out.Delete(newId)
	assert.Nil(t, err)

	userExists, err := out.Exists(newUser)
	assert.Nil(t, err)
	assert.False(t, userExists)

	userById, err := out.GetById(context.TODO(), newId)
	assert.Nil(t, err)
	assert.Nil(t, userById)

	userByName, err := out.GetByName(u.Name)
	assert.Nil(t, err)
	assert.Nil(t, userByName)
}

func Test_MatchPassword_UserNotExists_False(t *testing.T) {
	u := getTestUser(t)

	matching, err := out.MatchPassword(u.Name, u.Password)
	assert.Nil(t, err)

	assert.False(t, matching)
}

func Test_MatchPassword_UserExists_PasswordIsCompared(t *testing.T) {
	u := getTestUser(t)

	out.Add(u)
	matching, err := out.MatchPassword(u.Name, "wrongPwd")
	assert.Nil(t, err)
	assert.False(t, matching)

	matching, err = out.MatchPassword(u.Name, u.Password)
	assert.Nil(t, err)
	assert.True(t, matching)
}

func Test_Add_Password_PasswordIsSavedHashed(t *testing.T) {
	u := getTestUser(t)
	u.Password = "pwd"
	pwdHash := "oRWenfNnDVSdBFJFMmKfVHfOt97sm0XkfowAlQbsssg="
	origPwd := u.Password

	out.Add(u)
	matching, err := out.MatchPassword(u.Name, origPwd)
	assert.Nil(t, err)
	assert.False(t, matching)

	matching, err = out.MatchPassword(u.Name, pwdHash)
	assert.Nil(t, err)
	assert.True(t, matching)
}

func getTestUser(t *testing.T) *user.User {
	regT := time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
	return &user.User{
		Name:             t.Name() + "_" + testSuffix,
		Email:            t.Name() + "_" + testSuffix + "@gmail.com",
		Password:         "pwd",
		RegistrationDate: regT,
	}
}
