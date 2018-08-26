package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/kskitek/user-service/http_boundary"
	"github.com/kskitek/user-service/user"
)

func NewHandler() http_boundary.Handler {
	return &handler{
		s: NewService(),
	}
}

func NewService() Service {
	return &service{
		authenticator: NewAuthenticator(),
		userDao:       user.NewDao(),
	}
}

func NewAuthenticator() Authenticator {
	return &jwtAuthenticator{
		method: jwt.SigningMethodHS256,
		secret: "a2P5dgR2ya", // TODO Token from env?
	}
}
