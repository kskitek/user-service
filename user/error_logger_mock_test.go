package user

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/kskitek/arecar/user-service/events"
)

type testHook struct {
	lastError        error
	lastNotification *events.Notification
}

func (*testHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.ErrorLevel}
}

func (h *testHook) Fire(e *logrus.Entry) error {
	h.lastError = e.Data["error"].(error)
	n, ok := e.Data["notification"].(*events.Notification)
	if ok {
		h.lastNotification = n
	} else {
		h.lastNotification = nil
	}
	return nil
}
