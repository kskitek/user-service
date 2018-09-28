package user

import "context"

type Dao interface {
	GetById(context.Context, int64) (*User, error)
	GetByName(string) (*User, error)
	MatchPassword(userName string, password string) (bool, error)
	Exists(*User) (bool, error)
	Add(*User) (*User, error)
	Delete(int64) error
}
