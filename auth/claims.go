package auth

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"fmt"
)

type authClaims struct {
	User string `json:"user"`
	*jwt.StandardClaims
}

func (c *authClaims) Valid() error {
	if !c.VerifyExpiresAt(time.Now().Unix(), true) {
		return fmt.Errorf("token expired")
	}
	if !c.verifyUserId(true) {
		return fmt.Errorf("invalid token")
	}
	return nil
}
func (c *authClaims) verifyUserId(req bool) bool {
	if c.User == "" && req {
		return false
	} else {
		return true
	}
}
