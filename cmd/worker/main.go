package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/example/distrib-jobs/internal/broker"
	"github.com/example/distrib-jobs/internal/jobs"
	"github.com/nats-io/nats.go"
)

func main() {
	natsURL := getenv("NATS_URL", "nats://127.0.0.1:4222")
	stream := getenv("STREAM_NAME", "JOBS")
	workerType := getenv("WORKER_TYPE", "default")
	consumerName := getenv("CONSUMER_NAME", "worker-"+workerType)

	js, err := broker.Connect(natsURL)
	if err != nil {
		log.Fatal(err)
	}
	if err := broker.EnsureStream(js, stream, []string{"jobs.*"}); err != nil {
		log.Fatal(err)
	}

	subject := "jobs." + workerType
	ctx, cancel := signalContext()
	defer cancel()

	_, err = js.PullSubscribe(subject, consumerName, nats.BindStream(stream), nats.ManualAck())
	if err != nil && !strings.Contains(err.Error(), "already in use") {
		log.Fatal(err)
	}

	sub, err := js.PullSubscribe(subject, consumerName, nats.BindStream(stream), nats.ManualAck())
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			msgs, err := sub.Fetch(10)
			if err != nil && !strings.Contains(err.Error(), "timeout") {
				log.Println(err)
				continue
			}
			for _, m := range msgs {
				if err := jobs.Process(m.Data); err != nil {
					log.Println(err)
					_ = m.Nak()
					continue
				}
				_ = m.Ack()
			}
		}
	}
}

func getenv(k, d string) string {
	v := os.Getenv(k)
	if v == "" {
		return d
	}
	return v
}

func signalContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-ch
		cancel()
	}()
	return ctx, cancel
}
