package user

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func Test_WhenGetUserGivenEmptyIdThenError(t *testing.T) {
	out := &crud{NewMockDao()}
	var id int64

	_, apiError := out.Get(id)

	assert.NotNil(t, apiError)
}

func Test_WhenGetUserGivenErrorInDaoThenError(t *testing.T) {
	out := &crud{NewMockDao()}

	_, apiError := out.Get(UserErrorId)

	assert.NotNil(t, apiError)
}

func Test_WhenGetUserGivenNoUserForIdThenError(t *testing.T) {
	out := &crud{NewMockDao()}
	notExistingId := int64(100100)

	_, apiError := out.Get(notExistingId)

	assert.NotNil(t, apiError)
}

func Test_WhenAddGivenNilUserThenError(t *testing.T) {
	out := &crud{NewMockDao()}
	var user *User

	_, apiError := out.Add(user)

	assert.NotNil(t, apiError)
}

func Test_WhenAddGivenErrorInDaoThenError(t *testing.T) {
	out := &crud{NewMockDao()}

	_, apiError := out.Add(UserError())

	assert.NotNil(t, apiError)
}

func Test_WhenAddGivenUserExistsThenError(t *testing.T) {
	out := &crud{NewMockDao()}

	_, apiError := out.Add(UserExists())

	assert.Error(t, apiError)
}

func Test_WhenAddGivenOkThenReturnedPasswordIsEmpty(t *testing.T) {
	out := &crud{NewMockDao()}

	user, apiError := out.Add(UserOk())

	assert.Nil(t, apiError)

	assert.Equal(t, "", user.Password)
}

func Test_WhenGetUserGivenUserHasPasswordThenReturnedPasswordIsEmpty(t *testing.T) {
	out := &crud{NewMockDao()}

	user, apiError := out.Get(UserExistsId)

	assert.Nil(t, apiError)

	assert.Equal(t, "", user.Password)
}
