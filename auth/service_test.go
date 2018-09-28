package auth

import (
	"context"
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
	_, apiError := out.Login(context.TODO(), UserErrorName, "")

	assert.NotNil(t, apiError)
}

func Test_Login_PasswordMatchesInDao_ReturnToken(t *testing.T) {
	token, apiError := out.Login(context.TODO(), UserOkName, UserOkPassword)

	assert.Nil(t, apiError)
	assert.Equal(t, UserOkToken, token)
}

func Test_Login_PasswordNotMatchesInDao_Error(t *testing.T) {
	_, apiError := out.Login(context.TODO(), UserOkName, "WrongPassword")

	assert.NotNil(t, apiError)
}

func Test_Login_ErrorInAuthenticator_Error(t *testing.T) {
	_, apiError := out.Login(context.TODO(), UserErrorAuthName, UserOkPassword)

	assert.NotNil(t, apiError)
}

func Test_Login_LoggsIn_Notifies(t *testing.T) {
	c := prepareNotificationTest(t, AuthTopic+".login")

	token, _ := out.Login(context.TODO(), UserOkName, UserOkPassword)
	n := waitForNotification(t, c)

	assert.Equal(t, token, n.Token)
	assert.Equal(t, UserOkName, n.Payload)
	assert.NotNil(t, n.When)
}

func prepareNotificationTest(t *testing.T, topic string) chan event.Notification {
	c := make(chan event.Notification)
	f := func(n event.Notification) {
		c <- n
	}
	err := out.notifier.AddListener(topic, f)
	assert.NoError(t, err)
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
