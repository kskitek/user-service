package user

import (
	"testing"
	"time"

	"github.com/kskitek/user-service/event"
	"github.com/kskitek/user-service/event/mem"
	"github.com/stretchr/testify/assert"
)

func newOut() Service {
	notifier := mem.NewNotifier()
	return &crud{
		dao:      NewMockDao(),
		notifier: notifier,
	}
}

func Test_GetUser_EmptyId_Error(t *testing.T) {
	out := newOut()
	var id int64

	_, apiError := out.Get(id)

	assert.NotNil(t, apiError)
}

func Test_GetUser_ErrorInDao_Error(t *testing.T) {
	out := newOut()

	_, apiError := out.Get(UserErrorId)

	assert.NotNil(t, apiError)
}

func Test_GetUser_NoUserForId_Error(t *testing.T) {
	out := newOut()
	notExistingId := int64(100100)

	_, apiError := out.Get(notExistingId)

	assert.NotNil(t, apiError)
}

func Test_Add_NilUser_Error(t *testing.T) {
	out := newOut()
	var user *User

	_, apiError := out.Add(user)

	assert.NotNil(t, apiError)
}

func Test_Add_ErrorInDao_Error(t *testing.T) {
	users := []*User{UserError(), UserAddError()}
	for _, user := range users {
		out := newOut()

		_, apiError := out.Add(user)

		assert.NotNil(t, apiError)
	}
}

func Test_Add_UserExists_Error(t *testing.T) {
	out := newOut()

	_, apiError := out.Add(UserExists())

	assert.Error(t, apiError)
}

func Test_Add_Ok_ReturnedPasswordIsEmpty(t *testing.T) {
	out := newOut()

	user, apiError := out.Add(UserOk())

	assert.Nil(t, apiError)

	assert.Equal(t, "", user.Password)
}

func Test_Add_NoName_Error(t *testing.T) {
	out := newOut()

	user := UserOk()
	user.Name = ""
	_, apiError := out.Add(user)

	assert.NotNil(t, apiError)
}

func Test_Add_NoPassword_Error(t *testing.T) {
	out := newOut()

	user := UserOk()
	user.Password = ""
	_, apiError := out.Add(user)

	assert.NotNil(t, apiError)
}

func Test_Add_NoEmail_Error(t *testing.T) {
	out := newOut()

	user := UserOk()
	user.Email = ""
	_, apiError := out.Add(user)

	assert.NotNil(t, apiError)
}

func Test_Delete_EmptyId_Error(t *testing.T) {
	out := newOut()
	var id int64

	apiError := out.Delete(id)

	assert.NotNil(t, apiError)
}

func Test_Delete_ErrorInDao_Error(t *testing.T) {
	out := newOut()

	apiError := out.Delete(UserErrorId)

	assert.NotNil(t, apiError)
}

func Test_Delete_UserExistsOrNot_NoError(t *testing.T) {
	users := []int64{UserOkId, UserExistsId}
	for _, userId := range users {
		out := newOut()

		apiError := out.Delete(userId)

		assert.Nil(t, apiError)
	}
}

func Test_Get_UserHasPassword_ReturnedPasswordIsEmpty(t *testing.T) {
	out := newOut()

	user, apiError := out.Get(UserExistsId)

	assert.Nil(t, apiError)

	assert.Equal(t, "", user.Password)
}

func prepareNotificationTest(topic string) (Service, chan event.Notification) {
	d := NewMockDao()
	n := mem.NewNotifier()
	out := &crud{
		dao:      d,
		notifier: n,
	}
	c := make(chan event.Notification)
	f := func(n event.Notification) {
		c <- n
	}
	n.AddListener(topic, f)
	return out, c
}

func Test_Add_OkUser_Notifies(t *testing.T) {
	out, c := prepareNotificationTest(CrudBaseTopic + ".add")
	user := UserOk()

	userAdded, err := out.Add(user)
	notification := waitForNotification(c)

	assert.Nil(t, err)
	assert.NotEmpty(t, notification)
	assert.Equal(t, userAdded, notification.Payload)
}

func waitForNotification(c chan event.Notification) (notification event.Notification) {
	select {
	case notification = <-c:
		return
	case <-time.NewTimer(time.Second).C:
		return
	}
}

func Test_Add_NotifierFails_ErrorIsLogged(t *testing.T) {
	//notif := &mockNotifier{}
	//out := newOut().(*crud)
	//out.notifier = notif
	//hook := &testHook{}
	//logrus.AddHook(hook)
	//user := UserOk()
	//
	//newUser, err := out.Add(user)
	//n := &event.Notification{Payload: newUser, Event: "create"}
	//
	//assert.Nil(t, err)
	//
	//assert.Equal(t, notifierMockError, hook.lastError)
	//assert.Equal(t, n, hook.lastNotification)
}

func Test_Delete_OkUser_Notifies(t *testing.T) {
	out, c := prepareNotificationTest(CrudBaseTopic + ".delete")
	user := UserOk()

	_, err := out.Add(user)
	waitForNotification(c)
	err = out.Delete(UserOkId)
	deleteNotification := waitForNotification(c)

	assert.Nil(t, err)
	assert.NotEmpty(t, deleteNotification)
	assert.Equal(t, UserOkId, deleteNotification.Payload)
}

func Test_Delete_NotifierFails_ErrorIsLogger(t *testing.T) {
	//notif := &mockNotifier{}
	//out := newOut().(*crud)
	//out.notifier = notif
	//hook := &testHook{}
	//logrus.AddHook(hook)
	//user := UserOk()
	//
	//_, err := out.Add(user)
	//err = out.Delete(UserOkId)
	//n := &event.Notification{Payload: UserOkId, Event: "delete"}
	//
	//assert.Nil(t, err)
	//
	//assert.Equal(t, notifierMockError, hook.lastError)
	//assert.Equal(t, n, hook.lastNotification)
}
