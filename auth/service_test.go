package auth

import (
	"context"
	"strconv"
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

func Test_Login_LogsIn_Notifies(t *testing.T) {
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

var validationCases = []struct {
	name           string
	expectedResult bool
}{
	{"Name1", true},
	{"", false},
	{"Name1\"", false},
	{"Name1'", false},
	{"Name1;", false},
	{"Name1#", false},
	{"11111111102222222220", true},
	{"111", false},
	{"111111111022222222203", false},
}

func Test_ValidateUserPayload_(t *testing.T) {
	for i, c := range validationCases {
		tf := func(t *testing.T) {
			result := validateName(c.name)

			t.Log("Payload:", c.name)
			assert.Equal(t, c.expectedResult, result)
		}

		t.Run(t.Name()+strconv.Itoa(i), tf)
	}
}
