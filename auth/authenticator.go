package auth

import (
	"time"
	"github.com/dgrijalva/jwt-go"
	"fmt"
)

// Implementations are responsible for providing authentication tokens and verifying tokens from user.
type Authenticator interface {
	// Returns authentication token for given user
	//
	// When expiation time is null it will be set with default value.
	GetToken(userId string, expirationTime *time.Time) (string, error)
	Parse(string) (*AuthResult, error)
}

// Result of parsing the authentication token
//
// This struct might be just in/out struct for GetToken/Parse
type AuthResult struct {
	UserId string
}

const (
	issuer = "user-service"
)

type jwtAuthenticator struct {
	method jwt.SigningMethod
	secret string
}

func (a *jwtAuthenticator) GetToken(userId string, expirationTime *time.Time) (string, error) {
	expTime := fixExpTimeWithDefault(expirationTime)
	claims := &authClaims{
		User: userId,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    issuer,
			Subject:   userId,
		},
	}
	token := jwt.NewWithClaims(a.method, claims)

	key := []byte(a.secret)
	return token.SignedString(key)
}

func (a *jwtAuthenticator) Parse(tokenString string) (*AuthResult, error) {
	token, err := jwt.ParseWithClaims(tokenString, &authClaims{}, a.jwtKeyFunc)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*authClaims)
	if !ok {
		return nil, fmt.Errorf("unkown token claims %V", token.Claims)
	}
	if err = claims.Valid(); err != nil {
		return nil, err
	}

	return claims.toAuthResult(), nil
}

// In the future this func will be able to verify `kid` and use it together with blacklist
func (a *jwtAuthenticator) jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	secret := []byte(a.secret)
	return secret, nil
}

func fixExpTimeWithDefault(expTime *time.Time) time.Time {
	if expTime == nil {
		t := time.Now().UTC().Add(time.Hour * 24 * 7)
		return t
	} else {
		return *expTime
	}
}

func (c *authClaims) toAuthResult() *AuthResult {
	return &AuthResult{
		UserId: c.User,
	}
}
