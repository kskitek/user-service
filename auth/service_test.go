package auth

import (
	"testing"
	"time"

	"github.com/kskitek/user-service/event"
	"github.com/kskitek/user-service/event/mem"
	"github.com/stretchr/testify/assert"
)

var out = &service{
	userDao:       NewDaoMock(),
	authenticator: NewAuthMock(),
	notifier:      mem.NewNotifier(),
}

func Test_Login_DaoError_Error(t *testing.T) {
	_, apiError := out.Login(UserErrorName, "")

	assert.NotNil(t, apiError)
}

func Test_Login_PasswordMatchesInDao_ReturnToken(t *testing.T) {
	token, apiError := out.Login(UserOkName, UserOkPassword)

	assert.Nil(t, apiError)
	assert.Equal(t, UserOkToken, token)
}

func Test_Login_PasswordNotMatchesInDao_Error(t *testing.T) {
	_, apiError := out.Login(UserOkName, "WrongPassword")

	assert.NotNil(t, apiError)
}

func Test_Login_ErrorInAuthenticator_Error(t *testing.T) {
	_, apiError := out.Login(UserErrorAuthName, UserOkPassword)

	assert.NotNil(t, apiError)
}

func Test_Login_LoggsIn_Notifies(t *testing.T) {
	c := prepareNotificationTest(AuthTopic + ".login")

	token, _ := out.Login(UserOkName, UserOkPassword)
	n := waitForNotification(t, c)

	assert.Equal(t, token, n.Token)
	assert.Equal(t, UserOkName, n.Payload)
	assert.NotNil(t, n.When)
}

func prepareNotificationTest(topic string) chan event.Notification {
	c := make(chan event.Notification)
	f := func(n event.Notification) {
		c <- n
	}
	out.notifier.AddListener(topic, f)
	return c
}

func waitForNotification(t *testing.T, c chan event.Notification) (notification event.Notification) {
	t.Helper()
	select {
	case notification = <-c:
		return
	case <-time.NewTimer(time.Second).C:
		t.Fatal("Notification timeout")
		return
	}
}
