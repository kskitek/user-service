package mem

import (
	"github.com/kskitek/user-service/event"
	"github.com/sirupsen/logrus"
)

func NewNotifier() event.Notifier {
	return &Mem{make(map[string][]event.Listener)}
}

type Mem struct {
	Listeners map[string][]event.Listener
}

func (m *Mem) Notify(topic string, n event.Notification) error {
	logrus.WithFields(logrus.Fields{"t": topic, "n": n}).Debug("Notifying")
	listeners := m.Listeners[topic]
	for _, l := range listeners {
		go l(n)
	}

	return nil
}

func (m *Mem) AddListener(topic string, l event.Listener) error {
	listeners, found := m.Listeners[topic]
	if found {
		m.Listeners[topic] = append(listeners, l)
	} else {
		m.Listeners[topic] = []event.Listener{l}
	}

	return nil
}
