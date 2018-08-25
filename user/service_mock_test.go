package user

import (
	"time"
	"fmt"
	"strconv"
)

var testSuffix = strconv.FormatInt(time.Now().Unix(), 10)

var UserOkId = int64(1)

func newUser(id string, name string) *User {
	return &User{
		Id:               id,
		Name:             name + testSuffix,
		Password:         "Pwd",
		Email:            name + testSuffix + "@gmail.com",
		RegistrationDate: time.Now().UTC(),
	}
}

func UserOk() *User {
	return newUser("-1", "User1_")
}

func UserOk2() *User {
	return newUser("-1", "User2_")
}

var UserExistsId = int64(20)

func UserExists() *User {
	return newUser("20", "User20_")
}

var UserErrorId = int64(10)

func UserError() *User {
	return newUser("10", "User10_")
}

func UserAddError() *User {
	return newUser("11", "User11_")
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
