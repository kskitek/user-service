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

func Test_GivenUserIdThenUserIdIsPartOfToken(t *testing.T) {
	out := newOut()

	expTime := time.Now().Add(time.Hour * 1)

	userId := "1742012414"
	token, err := out.GetToken(userId, &expTime)
	assert.NoError(t, err)

	result, err := out.Parse(token)

	assert.NoError(t, err)
	assert.Equal(t, userId, result.UserId)
}

func Test_GivenNilExpirationTimeThenDefaultIsSet(t *testing.T) {
	out := newOut()

	var expTime *time.Time

	userId := "1742012414"
	_, err := out.GetToken(userId, expTime)
	assert.NoError(t, err)
}

func Test_GivenNilExpirationTimeThenDefaultIsSetToOneWeek(t *testing.T) {
	actual := fixExpTimeWithDefault(nil)
	expected := time.Now().UTC().AddDate(0, 0, 7)

	assert.Equal(t, expected.Unix()/1000, actual.Unix()/1000)
}

func Test_GivenNoUserIdThenException(t *testing.T) {
	out := newOut()

	expTime := time.Now().Add(time.Hour * 1)

	token, err := out.GetToken("", &expTime)
	assert.NoError(t, err)

	_, err = out.Parse(token)

	assert.Error(t, err)
}

func Test_GivenTokenExpiredThenException(t *testing.T) {
	out := newOut()

	expTime := time.Now().AddDate(-1, 0, 0)

	token, err := out.GetToken("1234", &expTime)
	assert.NoError(t, err)

	_, err = out.Parse(token)

	assert.Error(t, err)
}
