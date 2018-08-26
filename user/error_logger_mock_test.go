package user

import (
	"github.com/kskitek/user-service/event"
	"github.com/sirupsen/logrus"
)

type testHook struct {
	lastError        error
	lastNotification *event.Notification
}

func (*testHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.ErrorLevel}
}

func (h *testHook) Fire(e *logrus.Entry) error {
	h.lastError = e.Data["error"].(error)
	n, ok := e.Data["notification"].(*event.Notification)
	if ok {
		h.lastNotification = n
	} else {
		h.lastNotification = nil
	}
	return nil
}
