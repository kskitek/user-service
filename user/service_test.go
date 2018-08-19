package user

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"gitlab.com/kskitek/arecar/user-service/events"
)

func newOut() Service {
	return &crud{
		dao:      NewMockDao(),
		notifier: events.NewInMemNotifier(),
	}
}

func newOutPlus() (Service, Dao, events.Notifier) {
	d := NewMockDao()
	n := events.NewInMemNotifier()
	c := &crud{
		dao:      d,
		notifier: n,
	}

	return c, d, n
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

func Test_Add_OkUser_Notifies(t *testing.T) {
	out, _, notif := newOutPlus()
	notifier := notif.(*events.MemNotifier)
	user := UserOk()

	userAdded, err := out.Add(user)

	assert.Nil(t, err)

	assert.NotEmpty(t, notifier.Events)
	lastEvent := notifier.Events[0]
	assert.Equal(t, userAdded, lastEvent.Payload)
}

//func Test_Add_NotifierFails_NothingHappens?(t *testing.T) {
//
//}
