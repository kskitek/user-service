package redis

import (
	"github.com/go-redis/redis"
	"github.com/kskitek/user-service/event"
	"github.com/sirupsen/logrus"
)

func NewNotifier() (event.Notifier, error) {
	// TODO envs
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		logrus.WithError(err).Error("error when setting up redis client")
		return nil, err
	}
	return &redisNotifier{client: client}, nil
}

type redisNotifier struct {
	client *redis.Client
}

func (r *redisNotifier) Notify(topic string, n event.Notification) error {
	cmd := r.client.Publish(topic, n)
	logrus.WithField("notification", n).Debug("Notify")
	return cmd.Err()
}

func (r *redisNotifier) AddListener(topic string, n event.Listener) error {
	pubSub := r.client.Subscribe(topic)
	_, err := pubSub.Receive()
	if err != nil {
		return err
	}
	c := pubSub.Channel()
	go channelListener(c, n, topic)

	return nil
}

func channelListener(c <-chan *redis.Message, l event.Listener, t string) {
	for m := range c {
		logrus.WithField("message", m).Debug("Notified")
		n := event.Notification{}
		err := n.UnmarshalBinary([]byte(m.Payload))
		if err != nil {
			logrus.WithError(err).WithField("message", m).
				Error("error when unmarshaling message. Did not notify listener.")
		}
		l(n)
	}
	logrus.WithField("topic", t).Warning("closing channel")
}
