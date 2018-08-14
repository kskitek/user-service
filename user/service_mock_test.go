package user

import (
	"time"
	"strconv"
	"sync/atomic"
	"fmt"
)

var UserOkId = int64(1)
var UserOk = &User{
	Id:               "1",
	Name:             "User1",
	Password:         "Pwd",
	Email:            "user1@gmail.com",
	RegistrationDate: time.Now(),
}

var UserErrorId = int64(10)
var UserError = &User{
	Id:               "10",
	Name:             "User10",
	Password:         "Pwd",
	Email:            "user10@gmail.com",
	RegistrationDate: time.Now(),
}

func NewMockDao() Dao {
	dao := &MockDao{
		currId:    0,
		mem:       make(map[string]*User),
		memByName: make(map[string]*User),
	}
	prepareMockData(dao)

	return dao
}
func prepareMockData(dao *MockDao) {
	dao.Add(UserOk)
}

type MockDao struct {
	mem       map[string]*User
	memByName map[string]*User
	currId    int64
}

func (d *MockDao) GetById(id int64) (*User, error) {
	if id == UserErrorId {
		return nil, fmt.Errorf("test error")
	}
	idStr := strconv.FormatInt(id, 10)
	user := d.mem[idStr]
	return user, nil
}

func (d *MockDao) GetByName(name string) (*User, error) {
	if name == UserError.Name {
		return nil, fmt.Errorf("test error")
	}
	user := d.memByName[name]
	return user, nil
}

func (d *MockDao) MatchPassword(userName string, password string) (bool, error) {
	if userName == UserError.Name {
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
	if user == UserError {
		return false, fmt.Errorf("test error")
	}
	_, exists := d.memByName[user.Name]
	return exists, nil
}

func (d *MockDao) Add(user *User) (*User, error) {
	if user == UserError {
		return nil, fmt.Errorf("test error")
	}
	newId := atomic.AddInt64(&d.currId, 1)
	user.Id = strconv.FormatInt(newId, 10)
	user.RegistrationDate = time.Now().UTC()

	d.mem[user.Id] = user
	d.memByName[user.Name] = user
	return user, nil
}

func (d *MockDao) Delete(id int64) error {
	if id == UserErrorId {
		return fmt.Errorf("test error")
	}
	idStr := strconv.FormatInt(id, 10)
	user, ok := d.mem[idStr]
	if ok {
		delete(d.mem, idStr)
		delete(d.memByName, user.Name)
	}
	return nil
}
