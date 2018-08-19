package events

func NewInMemNotifier() Notifier {
	return &MemNotifier{}
}

type MemNotifier struct {
	Events []*Notification
}

func (mem *MemNotifier) Notify(n *Notification) error {
	mem.Events = append(mem.Events, n)
	return nil
}
