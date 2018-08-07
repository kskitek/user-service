package auth

import (
	"time"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"errors"
)

// Implementations are responsible for providing authentication tokens and verifying tokens from user.
type Authenticator interface {
	// Returns authentication token for given user
	//
	// When expiation time is null it will be set with default value.
	GetToken(userId string, expirationTime *time.Time) (string, error)
	Verify(string) error
}

type jwtAuthenticator struct {
	method jwt.SigningMethod
	secret string
}

func NewAuthenticator() Authenticator {
	return &jwtAuthenticator{
		method: jwt.SigningMethodHS256,
		secret: "a2P5dgR2ya", // TODO Token from env?
	}
}

func (a *jwtAuthenticator) GetToken(user string, expTime *time.Time) (string, error) {
	fixExpTimeWithDefault(expTime)
	claims := &jwt.StandardClaims{
		ExpiresAt: (*expTime).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   user,
		Issuer: "me",
	}
	token := jwt.NewWithClaims(a.method, claims)

	key := []byte(user)
	return token.SignedString(key)
}

func (a *jwtAuthenticator) Verify(key string) error {
	segments := strings.Split(key, ".")
	if len(segments) != 3 {
		return errors.New("")
	}
	return a.method.Verify(segments[0], segments[1], segments[2])
}

func fixExpTimeWithDefault(expTime *time.Time) {
	if expTime == nil {
		t := time.Now()
		expTime = &t
		expTime.Add(time.Hour * 24 * 7)
	}
}
