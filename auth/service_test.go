package auth

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

var out = &service{
	userDao:       NewDaoMock(),
	authenticator: NewAuthMock(),
}

func Test_Login_DaoError_Error(t *testing.T) {
	_, apiError := out.Login(UserErrorName, "")

	assert.NotNil(t, apiError)
}

func Test_Login_PasswordMatchesInDao_ReturnToken(t *testing.T) {
	token, apiError := out.Login(UserOkName, UserOkPassword)

	assert.Nil(t, apiError)
	assert.Equal(t, UserOkToken, token)
}

func Test_Login_PasswordNotMatchesInDao_Error(t *testing.T) {
	_, apiError := out.Login(UserOkName, "WrongPassword")

	assert.NotNil(t, apiError)
}

func Test_Login_ErrorInAuthenticator_Error(t *testing.T) {
	_, apiError := out.Login(UserErrorAuthName, UserOkPassword)

	assert.NotNil(t, apiError)
}
