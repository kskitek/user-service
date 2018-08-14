package user

import (
	"strconv"
	"sync/atomic"
	"time"
)

type Dao interface {
	GetById(int64) (*User, error)
	GetByName(string) (*User, error)
	MatchPassword(userName string, password string) (bool, error)
	Exists(*User) (bool, error)
	Add(*User) (*User, error)
	Delete(int64) error
}

type MongoDao struct {
}

func (d *MongoDao) GetUser(int64) (*User, error) {
	return nil, nil
}

func (d *MongoDao) UserExists(*User) (bool, error) {
	return false, nil
}

func (d *MongoDao) AddUser(user *User) (*User, error) {
	return user, nil
}

type InMemDao struct {
	mem        map[string]*User
	memByEmail map[string]*User
	memByName  map[string]*User
	currId     int64
}

func (d *InMemDao) GetById(id int64) (*User, error) {
	idStr := strconv.FormatInt(id, 10)
	user := d.mem[idStr]
	return user, nil
}

func (d *InMemDao) GetByName(name string) (*User, error) {
	user := d.memByName[name]
	return user, nil
}

func (d *InMemDao) MatchPassword(userName string, password string) (bool, error) {
	user, err := d.GetByName(userName)
	if err != nil {
		return false, err
	}

	pwdMatching := user != nil && user.Password == password
	return pwdMatching, nil
}

func (d *InMemDao) Exists(user *User) (bool, error) {
	_, exists := d.memByEmail[user.Email]
	return exists, nil
}

func (d *InMemDao) Add(user *User) (*User, error) {
	newId := atomic.AddInt64(&d.currId, 1)
	user.Id = strconv.FormatInt(newId, 10)
	user.RegistrationDate = time.Now().UTC()

	d.mem[user.Id] = user
	d.memByName[user.Name] = user
	d.memByEmail[user.Email] = user
	return user, nil
}

func (d *InMemDao) Delete(id int64) error {
	idStr := strconv.FormatInt(id, 10)
	user, ok := d.mem[idStr]
	if ok {
		delete(d.mem, idStr)
		delete(d.memByName, user.Name)
		delete(d.memByEmail, user.Email)
	}
	return nil
}
