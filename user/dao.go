package user

import (
	"strconv"
	"sync/atomic"
	"time"
	"fmt"
	_ "github.com/lib/pq"
	"database/sql"
	"github.com/sirupsen/logrus"
	"os"
)

const defaultTimeout = "5m"

type Dao interface {
	GetById(int64) (*User, error)
	GetByName(string) (*User, error)
	MatchPassword(userName string, password string) (bool, error)
	Exists(*User) (bool, error)
	Add(*User) (*User, error)
	Delete(int64) error
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
	_, exists := d.memByName[user.Name]
	return exists, nil
}

func (d *InMemDao) Add(user *User) (*User, error) {
	newId := atomic.AddInt64(&d.currId, 1)
	user.Id = strconv.FormatInt(newId, 10)
	user.RegistrationDate = time.Now().UTC()

	_, byName := d.memByName[user.Name]
	_, byEmail := d.memByEmail[user.Email]
	if byName || byEmail {
		return nil, fmt.Errorf("user already exists")
	}

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

const (
	allUserFields = " id, name, email, password, creationDate "
	getByPattern  = "SELECT " + allUserFields + " FROM users WHERE %s = $1"
)

type pgDao struct {
	db *sql.DB
}

func (d *pgDao) GetById(id int64) (*User, error) {
	return d.getUser(fmt.Sprintf(getByPattern, "id"), id)
}

func (d *pgDao) GetByName(name string) (*User, error) {
	return d.getUser(fmt.Sprintf(getByPattern, "name"), name)
}

func (d *pgDao) getUser(query string, params ...interface{}) (*User, error) {
	u := &User{}
	err := d.db.QueryRow(query, params...).
		Scan(&u.Id, &u.Name, &u.Email, &u.Password, &u.RegistrationDate)
	u.RegistrationDate = u.RegistrationDate.In(time.UTC)

	switch err {
	case nil:
		return u, nil
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

func (d *pgDao) MatchPassword(userName string, password string) (bool, error) {
	user, err := d.getUser("SELECT "+allUserFields+" FROM users WHERE name = $1 AND password = $2", userName, password)
	return user != nil, err
}

func (d *pgDao) Exists(u *User) (bool, error) {
	user, err := d.GetByName(u.Name)
	return user != nil, err
}

func (d *pgDao) Add(u *User) (*User, error) {
	var id int64
	regDate := u.RegistrationDate.In(time.UTC)
	err := d.db.QueryRow("INSERT INTO users (name, email, password, creationDate) VALUES ($1, $2, $3, $4) RETURNING id",
		u.Name, u.Email, u.Password, regDate).Scan(&id)
	if err == nil {
		u.Id = strconv.FormatInt(id, 10)
		return u, err
	}

	return nil, err
}

func (d *pgDao) Delete(id int64) error {
	_, err := d.db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}

func NewPgDao() Dao {
	connStr := os.Getenv("DB_CONN_STR")
	timeout := getTimeout()
	db := getDb(connStr, timeout)

	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Minute * 10)
	db.SetMaxIdleConns(10 / 2)

	return &pgDao{db}
}

func getTimeout() time.Duration {
	timeoutStr := os.Getenv("DB_TIMEOUT")
	if timeoutStr == "" {
		timeoutStr = defaultTimeout
	}
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		logrus.WithField("timeoutStr", timeoutStr).Panic("cannot read database timeout string")
	}
	return timeout
}

func getDb(connStr string, timeout time.Duration) *sql.DB {
	start := time.Now()
	logrus.Info("trying to connect to database")

	for {
		db, err := sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
		}
		if err != nil {
			logrus.WithError(err).Warning("unable to connect to database")
			time.Sleep(timeout / 10)
			timePassed := time.Since(start)
			if timePassed > timeout {
				logrus.WithFields(logrus.Fields{"timeout": timeout, "connStr": connStr}).Panic("unable to connect to database within time limit")
			}
		} else {
			logrus.Info("got connection to database")
			return db
		}
	}
}
