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

	userId := "1742012414"
	token, err := out.GetToken(userId, &expTime)
	assert.NoError(t, err)

	result, err := out.Parse(token)

	assert.NoError(t, err)
	assert.Equal(t, userId, result.UserId)
}

/* TODO
-validate checks expiration date
- validate checks presence of userId
- getToken sets `exp` when expTime is nil
- fixExpTime sets a 24H * 7
 */
