package event

import (
	"encoding/json"
	"time"
)

type Notification struct {
	CorrelationId string      `json:"correlationId,omitempty"`
	Token         string      `json:"token,omitempty"`
	When          time.Time   `json:"time,omitempty"`
	Event         string      `json:"event,omitempty"`
	Payload       interface{} `json:"payload"`
}

type Notifier interface {
	Notify(topic string, n Notification) error
	AddListener(topic string, n Listener) error
	AddListenerPattern(topicPattern string, n Listener) error
}

type Listener = func(Notification)

func (n Notification) MarshalBinary() ([]byte, error) {
	return json.Marshal(n)
}

func (n *Notification) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, n)
}

func (n Notification) String() string {
	bytes, err := n.MarshalBinary()
	if err != nil {
		return ""
	}
	return string(bytes)
}
