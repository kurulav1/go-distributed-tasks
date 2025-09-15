package broker

import (
	"time"

	"github.com/nats-io/nats.go"
)

func Connect(url string) (nats.JetStreamContext, error) {
	nc, err := nats.Connect(url, nats.Name("distrib-jobs"))
	if err != nil {
		return nil, err
	}
	return nc.JetStream(nats.PublishAsyncMaxPending(256), nats.MaxWait(5*time.Second))
}

func EnsureStream(js nats.JetStreamContext, name string, subjects []string) error {
	_, err := js.StreamInfo(name)
	if err == nil {
		return nil
	}
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     name,
		Subjects: subjects,
		Replicas: 1,
		Storage:  nats.FileStorage,
	})
	return err
}
