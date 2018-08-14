package user

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func Test_WhenGetUserGivenEmptyIdThenError(t *testing.T) {
	out := &crud{NewMockDao()}
	var id int64

	_, apiError := out.Get(id)

	assert.Error(t, apiError)
}

func Test_WhenGetUserGivenErrorInDaoThenError(t *testing.T) {
	out := &crud{NewMockDao()}

	_, apiError := out.Get(UserErrorId)

	assert.Error(t, apiError)
}

func Test_WhenGetUserGivenNoUserForIdThenError(t *testing.T) {
	out := &crud{NewMockDao()}
	notExistingId := int64(100100)

	_, apiError := out.Get(notExistingId)

	assert.Error(t, apiError)
}

func Test_WhenGetUserGivenUserHasPasswordThenReturnedPasswordIsEmpty(t *testing.T) {
	out := &crud{NewMockDao()}

	user, apiError := out.Get(UserOkId)

	assert.Nil(t, apiError)

	assert.Equal(t, "", user.Password)
}
