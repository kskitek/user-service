// +build it

package redis

import (
	"github.com/stretchr/testify/assert"
	"github.com/kskitek/user-service/event"
	"testing"
	"time"
)

func Test_NewNotifierIsCreatedWithoutError(t *testing.T) {
	_, err := NewNotifier()
	assert.NoError(t, err)
}

func Test_AddListener_NewMessage_AddedListenerIsNotified(t *testing.T) {
	out, err := NewNotifier()
	assert.NoError(t, err)

	topic := "test-topic"
	c := make(chan event.Notification)
	l := func(n event.Notification) {
		c <- n
	}
	expected := event.Notification{Payload: "test", Token: "t", When: time.Now()}

	err = out.AddListener(topic, l)
	assert.NoError(t, err)
	out.Notify(topic, expected)

	actual := <-c
	assert.Equal(t, expected.When.Unix(), actual.When.Unix())
	assert.Equal(t, expected.Token, actual.Token)
	assert.Equal(t, expected.Payload, actual.Payload)
}
