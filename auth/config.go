package auth

import (
	"gitlab.com/kskitek/arecar/user-service/http_boundary"
	"gitlab.com/kskitek/arecar/user-service/user"
	"github.com/dgrijalva/jwt-go"
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
