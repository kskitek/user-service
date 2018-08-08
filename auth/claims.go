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
	now := time.Now().Unix()
	if !c.VerifyExpiresAt(now, true) {
		return fmt.Errorf("token expired")
	}
	if !c.verifyUserId(true) {
		return fmt.Errorf("invalid token")
	}
	if !c.VerifyIssuedAt(now, false) {
		return fmt.Errorf("token used before issued")
	}
	if !c.VerifyIssuer(issuer, false) {
		return fmt.Errorf("unknown issuer")
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
