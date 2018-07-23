package user

import (
	"strconv"
	"sync/atomic"
	"time"
)

type Dao interface {
	GetUser(int64) (*User, error)
	UserExists(*User) (bool, error)
	AddUser(*User) (*User, error)
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
	currId     int64
}

func (d *InMemDao) GetUser(id int64) (*User, error) {
	idStr := strconv.FormatInt(id, 10)
	user := d.mem[idStr]
	return user, nil
}

func (d *InMemDao) UserExists(user *User) (bool, error) {
	_, exists := d.memByEmail[user.Email]
	return exists, nil
}

func (d *InMemDao) AddUser(user *User) (*User, error) {
	newId := atomic.AddInt64(&d.currId, 1)
	user.Id = strconv.FormatInt(newId, 10)
	user.RegistrationDate = time.Now().UTC()

	d.mem[user.Id] = user
	d.memByEmail[user.Email] = user
	return user, nil
}
