package user

import (
	"time"
	"fmt"
)

var UserOkId = int64(1)

func UserOk() *User {
	return &User{
		Id:               "-1",
		Name:             "User1",
		Password:         "Pwd",
		Email:            "user1@gmail.com",
		RegistrationDate: time.Now(),
	}
}

func UserOk2() *User {
	return &User{
		Id:               "-1",
		Name:             "User2",
		Password:         "Pwd",
		Email:            "user2@gmail.com",
		RegistrationDate: time.Now(),
	}
}

var UserExistsId = int64(20)

func UserExists() *User {
	return &User{
		Id:               "20",
		Name:             "User20",
		Password:         "Pwd",
		Email:            "user20@gmail.com",
		RegistrationDate: time.Now(),
	}
}

var UserErrorId = int64(10)

func UserError() *User {
	return &User{
		Id:               "10",
		Name:             "User10",
		Password:         "Pwd",
		Email:            "user10@gmail.com",
		RegistrationDate: time.Now(),
	}
}

func UserAddError() *User {
	return &User{
		Id:               "11",
		Name:             "User11",
		Password:         "Pwd",
		Email:            "user11@gmail.com",
		RegistrationDate: time.Now(),
	}
}

func NewMockDao() Dao {
	return &MockDao{}
}

type MockDao struct {
}

func (d *MockDao) GetById(id int64) (*User, error) {
	if id == UserErrorId {
		return nil, fmt.Errorf("test error")
	}
	if id == UserExistsId {
		return UserExists(), nil
	}
	return nil, nil
}

func (d *MockDao) GetByName(name string) (*User, error) {
	if name == UserError().Name {
		return nil, fmt.Errorf("test error")
	}
	return UserOk(), nil
}

func (d *MockDao) MatchPassword(userName string, password string) (bool, error) {
	if userName == UserError().Name {
		return false, fmt.Errorf("test error")
	}
	user, err := d.GetByName(userName)
	if err != nil {
		return false, err
	}

	pwdMatching := user != nil && user.Password == password
	return pwdMatching, nil
}

func (d *MockDao) Exists(user *User) (bool, error) {
	if user.Id == UserError().Id {
		return false, fmt.Errorf("test error")
	}
	exists := user.Id == UserExists().Id
	return exists, nil
}

func (d *MockDao) Add(user *User) (*User, error) {
	if user.Id == UserError().Id || user.Id == UserAddError().Id {
		return nil, fmt.Errorf("test error")
	}
	return UserOk(), nil
}

func (d *MockDao) Delete(id int64) error {
	if id == UserErrorId {
		return fmt.Errorf("test error")
	}
	return nil
}
