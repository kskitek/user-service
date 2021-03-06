package user

import (
	"context"
	"strconv"
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

	_, apiError := out.Get(context.TODO(), id)

	assert.NotNil(t, apiError)
}

func Test_GetUser_ErrorInDao_Error(t *testing.T) {
	out := newOut()

	_, apiError := out.Get(context.TODO(), UserErrorId)

	assert.NotNil(t, apiError)
}

func Test_GetUser_NoUserForId_Error(t *testing.T) {
	out := newOut()
	notExistingId := int64(100100)

	_, apiError := out.Get(context.TODO(), notExistingId)

	assert.NotNil(t, apiError)
}

func Test_Add_NilUser_Error(t *testing.T) {
	out := newOut()
	var user *User

	_, apiError := out.Add(context.TODO(), user)

	assert.NotNil(t, apiError)
}

func Test_Add_ErrorInDao_Error(t *testing.T) {
	users := []*User{UserError(), UserAddError()}
	for _, user := range users {
		out := newOut()

		_, apiError := out.Add(context.TODO(), user)

		assert.NotNil(t, apiError)
	}
}

func Test_Add_UserExists_Error(t *testing.T) {
	out := newOut()

	_, apiError := out.Add(context.TODO(), UserExists())

	assert.Error(t, apiError)
}

func Test_Add_Ok_ReturnedPasswordIsEmpty(t *testing.T) {
	out := newOut()

	user, apiError := out.Add(context.TODO(), UserOk())

	assert.Nil(t, apiError)

	assert.Equal(t, "", user.Password)
}

func Test_Add_NoName_Error(t *testing.T) {
	out := newOut()

	user := UserOk()
	user.Name = ""
	_, apiError := out.Add(context.TODO(), user)

	assert.NotNil(t, apiError)
}

func Test_Add_NoPassword_Error(t *testing.T) {
	out := newOut()

	user := UserOk()
	user.Password = ""
	_, apiError := out.Add(context.TODO(), user)

	assert.NotNil(t, apiError)
}

func Test_Add_NoEmail_Error(t *testing.T) {
	out := newOut()

	user := UserOk()
	user.Email = ""
	_, apiError := out.Add(context.TODO(), user)

	assert.NotNil(t, apiError)
}

func Test_Delete_EmptyId_Error(t *testing.T) {
	out := newOut()
	var id int64

	apiError := out.Delete(context.TODO(), id)

	assert.NotNil(t, apiError)
}

func Test_Delete_ErrorInDao_Error(t *testing.T) {
	out := newOut()

	apiError := out.Delete(context.TODO(), UserErrorId)

	assert.NotNil(t, apiError)
}

func Test_Delete_UserExistsOrNot_NoError(t *testing.T) {
	users := []int64{UserOkId, UserExistsId}
	for _, userId := range users {
		out := newOut()

		apiError := out.Delete(context.TODO(), userId)

		assert.Nil(t, apiError)
	}
}

func Test_Get_UserHasPassword_ReturnedPasswordIsEmpty(t *testing.T) {
	out := newOut()

	user, apiError := out.Get(context.TODO(), UserExistsId)

	assert.Nil(t, apiError)

	assert.Equal(t, "", user.Password)
}

func prepareNotificationTest(t *testing.T, topic string) (Service, chan event.Notification) {
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
	err := n.AddListener(topic, f)
	assert.NoError(t, err)
	return out, c
}

func Test_Add_OkUser_Notifies(t *testing.T) {
	out, c := prepareNotificationTest(t, CrudBaseTopic+".add")
	user := UserOk()

	userAdded, err := out.Add(context.TODO(), user)
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

func Test_Delete_OkUser_Notifies(t *testing.T) {
	out, c := prepareNotificationTest(t, CrudBaseTopic+".delete")
	user := UserOk()

	_, err := out.Add(context.TODO(), user)
	assert.Nil(t, err)

	waitForNotification(c)
	err = out.Delete(context.TODO(), UserOkId)
	deleteNotification := waitForNotification(c)

	assert.Nil(t, err)
	assert.NotEmpty(t, deleteNotification)
	assert.Equal(t, UserOkId, deleteNotification.Payload)
}

var validationCases = []struct {
	in             User
	expectedResult bool
}{
	{User{Name: "Name1", Password: "Pwd", Email: "name1@email.com"}, true},
	{User{Name: "", Password: "Pwd", Email: "name1@email.com"}, false},
	{User{Name: "Name1", Password: "", Email: "name1@email.com"}, false},
	{User{Name: "Name1", Password: "Pwd", Email: ""}, false},
	{User{Name: "Name1\"", Password: "Pwd", Email: "name1@email.com"}, false},
	{User{Name: "Name1'", Password: "Pwd", Email: "name1@email.com"}, false},
	{User{Name: "Name1;", Password: "Pwd", Email: "name1@email.com"}, false},
	{User{Name: "Name1#", Password: "Pwd", Email: "name1@email.com"}, false},
	{User{Name: "Name1", Password: "Pwd", Email: "name1@email.com\""}, false},
	{User{Name: "Name1", Password: "Pwd", Email: "name1@email.com'"}, false},
	{User{Name: "Name1", Password: "Pwd", Email: "name1@email.com;"}, false},
	{User{Name: "Name1", Password: "Pwd", Email: "name1@email.com#"}, false},
	{User{Name: "11111111102222222220", Password: "Pwd", Email: "11111111102222222220@1111111110222222.com"}, true},
	{User{Name: "111", Password: "Pwd", Email: "11111111102222222220@11111111102222222220.com"}, false},
	{User{Name: "111111111022222222203", Password: "Pwd", Email: "11111111102222222220@1111111110222222.com"}, false},
	{User{Name: "11111111102222222222", Password: "Pwd", Email: "11111111102222222220@1111111110222222.com3"}, false},
}

func Test_ValidateUserPayload_(t *testing.T) {
	for i, c := range validationCases {
		tf := func(t *testing.T) {
			result := validateAddUserPayload(&c.in)

			t.Log("Payload:", c.in)
			t.Log("Error:", result)
			assert.Equal(t, c.expectedResult, result == nil)
		}

		t.Run(t.Name()+strconv.Itoa(i), tf)
	}
}
