package events

import "time"

type Notification struct {
	Token string
	When time.Time
	Payload interface{}
}

type Notifier interface {
	Notify(*Notification) error
}
