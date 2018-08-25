package event

import "time"

type Notification struct {
	//CorrelationId string `json:"correlationId,omitempty"`
	Token   string      `json:"token,omitempty"`
	When    time.Time   `json:"time,omitempty"`
	Payload interface{} `json:",inline"`
}

type Notifier interface {
	Notify(topic string, n Notification) error
	AddListener(topic string, n Listener) error
}

type Listener = func(Notification)
