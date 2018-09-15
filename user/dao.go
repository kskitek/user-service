package user

type Dao interface {
	GetById(int64) (*User, error)
	GetByName(string) (*User, error)
	MatchPassword(userName string, password string) (bool, error)
	Exists(*User) (bool, error)
	Add(*User) (*User, error)
	Delete(int64) error
}
