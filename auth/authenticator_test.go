package auth

import (
	"testing"
	"github.com/dgrijalva/jwt-go"
	"time"
	"github.com/stretchr/testify/assert"
)

func newOut() Authenticator {
	return &jwtAuthenticator{
		method: jwt.SigningMethodHS256,
		secret: "abc",
	}
}

func Test(t *testing.T) {
	out := newOut()

	expTime := time.Now().Add(time.Hour * 1)

	token, err := out.GetToken("1742012414", &expTime)

	assert.NoError(t, err)
	assert.Equal(t, "", token)
}
