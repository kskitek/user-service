package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/kskitek/user-service/event"
	"github.com/kskitek/user-service/user"
)

func NewService(dao user.Dao, notifier event.Notifier) Service {
	return &service{
		authenticator: NewAuthenticator(),
		userDao:       dao,
		notifier:      notifier,
	}
}

func NewAuthenticator() Authenticator {
	return &jwtAuthenticator{
		method: jwt.SigningMethodHS256,
		secret: "a2P5dgR2ya", // TODO Token from env?
	}
}
