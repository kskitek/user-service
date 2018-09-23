package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/kskitek/user-service/user"
)

var UserErrorName = "UserErrorName"
var UserOkName = "UserOkName"
var UserOkPassword = "UserOkPassword"
var UserOkToken = "Token"
var UserErrorAuthName = "UserErrorAuthName"

func NewDaoMock() user.Dao {
	return &mock{}
}

func NewAuthMock() Authenticator {
	return &auth{}
}

type mock struct {
}

type auth struct {
}

func (a *auth) GetToken(userId string, expirationTime *time.Time) (string, error) {
	if userId == UserErrorAuthName {
		return "", fmt.Errorf("test error")
	}
	return UserOkToken, nil
}

func (auth) Parse(string) (*Result, error) {
	panic("implement me")
}

func (mock) GetById(context.Context, int64) (*user.User, error) {
	panic("implement me")
}

func (mock) GetByName(string) (*user.User, error) {
	panic("implement me")
}

func (mock) MatchPassword(userName string, password string) (bool, error) {
	if userName == UserErrorName {
		return false, fmt.Errorf("test error")
	}
	if userName == UserOkName {
		return password == UserOkPassword, nil
	}

	return true, nil
}

func (mock) Exists(*user.User) (bool, error) {
	panic("implement me")
}

func (mock) Add(*user.User) (*user.User, error) {
	panic("implement me")
}

func (mock) Delete(int64) error {
	panic("implement me")
}
