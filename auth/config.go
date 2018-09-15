package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/kskitek/user-service/user"
)

func NewService(dao user.Dao) Service {
	return &service{
		authenticator: NewAuthenticator(),
		userDao:       dao,
	}
}

func NewAuthenticator() Authenticator {
	return &jwtAuthenticator{
		method: jwt.SigningMethodHS256,
		secret: "a2P5dgR2ya", // TODO Token from env?
	}
}
