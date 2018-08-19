package events

import "time"

type Notification struct {
	Token   string      `json:"token,omitempty"`
	When    time.Time   `json:"time,omitempty"`
	Event   string      `json:"event,omitempty"`
	Payload interface{} `json:",inline"`
}

type Notifier interface {
	Notify(topic string, n *Notification) error
}
