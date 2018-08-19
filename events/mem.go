package events

func NewInMemNotifier() Notifier {
	return &MemNotifier{}
}

type MemNotifier struct {
	Events []*Notification
	Topics []string
}

func (mem *MemNotifier) Notify(topic string, n *Notification) error {
	mem.Events = append(mem.Events, n)
	mem.Topics = append(mem.Topics, topic)
	return nil
}
