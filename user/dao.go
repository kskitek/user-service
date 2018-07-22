package user

type dao interface {
	AddUser(*User) (*User, error)
}

type mongoDao struct {
}

func (d *mongoDao) AddUser(user *User) (*User, error) {
	return nil, nil
}
