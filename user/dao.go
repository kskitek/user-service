package user

type Dao interface {
	GetUser(int) (*User, error)
	UserExists(*User) (bool, error)
	AddUser(*User) (*User, error)
}

type MongoDao struct {
}

func (d *MongoDao) GetUser(int) (*User, error) {
	return nil, nil
}

func (d *MongoDao) UserExists(*User) (bool, error) {
	return false, nil
}

func (d *MongoDao) AddUser(user *User) (*User, error) {
	user.Id = "1"
	return user, nil
}
