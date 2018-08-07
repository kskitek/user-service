package auth

import (
	"time"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"errors"
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

type jwtAuthenticator struct {
	method jwt.SigningMethod
	secret string
}

type authClaims struct {
	User string `json:"user"`
	*jwt.StandardClaims
}

func NewAuthenticator() Authenticator {
	return &jwtAuthenticator{
		method: jwt.SigningMethodHS256,
		secret: "a2P5dgR2ya", // TODO Token from env?
	}
}

func (a *jwtAuthenticator) GetToken(userId string, expTime *time.Time) (string, error) {
	fixExpTimeWithDefault(expTime)
	claims := authClaims{
		User: userId,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: (*expTime).Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   userId,
		},
	}
	token := jwt.NewWithClaims(a.method, claims)

	key := []byte(a.secret)
	return token.SignedString(key)
}

func (a *jwtAuthenticator) verify(token string) error {
	segments := strings.Split(token, ".")
	if len(segments) != 3 {
		return errors.New("")
	}
	key := []byte(a.secret)
	return a.method.Verify(segments[0], segments[1], key)
}

func (a *jwtAuthenticator) Parse(tokenString string) (*AuthResult, error) {
	//err := a.verify(tokenString)
	//if err != nil {
	//	return nil, err
	//}

	token, err := jwt.ParseWithClaims(tokenString, &authClaims{}, a.jwtKeyFunc)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*authClaims)
	if !ok {
		return nil, fmt.Errorf("unkown token claims %V", token.Claims)
	}

	return claims.toAuthResult(), nil
}

// In the future this func will be able to verify `kid` and use it together with blacklist
func (a *jwtAuthenticator) jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	secret := []byte(a.secret)
	return secret, nil
}

func fixExpTimeWithDefault(expTime *time.Time) {
	if expTime == nil {
		t := time.Now()
		expTime = &t
		expTime.Add(time.Hour * 24 * 7)
	}
}

func (c *authClaims) toAuthResult() *AuthResult {
	return &AuthResult{
		UserId: c.User,
	}
}
