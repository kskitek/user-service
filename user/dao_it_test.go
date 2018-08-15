// +build it

package user

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"strconv"
)

var out = NewDao()

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
	user := &User{
		Name:     t.Name(),
		Email:    t.Name() + "@gmail.com",
		Password: "pwd",
	}
	_, err1 := out.Add(user)
	_, err2 := out.Add(user)

	assert.Nil(t, err1)
	assert.NotNil(t, err2)
}

func Test_AddAndGet_ReturnsTheUser(t *testing.T) {
	user := getTestUser(t)

	newUser, err := out.Add(user)
	assert.Nil(t, err)
	newId, _ := strconv.ParseInt(newUser.Id, 10, 64)

	userByName, err := out.GetById(newId)
	assert.Nil(t, err)

	assert.Equal(t, newUser, userByName)
}

func Test_Add_GetByNameReturnsTheUser(t *testing.T) {
	user := getTestUser(t)

	newUser, err := out.Add(user)
	assert.Nil(t, err)

	userByName, err := out.GetByName(user.Name)
	assert.Nil(t, err)

	assert.Equal(t, newUser, userByName)
}

func Test_Add_ExistsReturnsTrue(t *testing.T) {
	user := getTestUser(t)

	newUser, err := out.Add(user)
	assert.Nil(t, err)

	userByName, err := out.Exists(newUser)
	assert.Nil(t, err)

	assert.True(t, userByName)
}

func Test_AfterDelete_CannotGetUser(t *testing.T) {
	user := getTestUser(t)

	newUser, err := out.Add(user)
	assert.Nil(t, err)

	newId, _ := strconv.ParseInt(newUser.Id, 10, 64)
	err = out.Delete(newId)
	assert.Nil(t, err)

	userExists, err := out.Exists(newUser)
	assert.Nil(t, err)
	assert.False(t, userExists)

	userById, err := out.GetById(newId)
	assert.Nil(t, err)
	assert.Nil(t, userById)

	userByName, err := out.GetByName(user.Name)
	assert.Nil(t, err)
	assert.Nil(t, userByName)
}

func Test_MatchPassword_UserNotExists_False(t *testing.T) {
	user := getTestUser(t)

	matching, err := out.MatchPassword(user.Name, "pwd")
	assert.Nil(t, err)

	assert.False(t, matching)
}

func Test_MatchPassword_UserExists_PasswordIsCompared(t *testing.T) {
	user := getTestUser(t)

	out.Add(user)
	matching, err := out.MatchPassword(user.Name, "wrongPwd")
	assert.Nil(t, err)
	assert.False(t, matching)

	matching, err = out.MatchPassword(user.Name, "pwd")
	assert.Nil(t, err)
	assert.True(t, matching)
}

// TODO check that password is saved hashed

func getTestUser(t *testing.T) *User {
	return &User{
		Name:     t.Name(),
		Email:    t.Name() + "@gmail.com",
		Password: "pwd",
	}
}
