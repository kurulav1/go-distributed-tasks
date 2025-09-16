package main

import (
	"log"
	"net/http"
	"os"

	"github.com/example/distrib-jobs/internal/auth"
	"github.com/example/distrib-jobs/internal/broker"
	"github.com/example/distrib-jobs/internal/db"
	"github.com/example/distrib-jobs/internal/httpapi"
)

func main() {
	natsURL := getenv("NATS_URL", "nats://127.0.0.1:4222")
	stream := getenv("STREAM_NAME", "JOBS")
	subjects := getenv("STREAM_SUBJECTS", "jobs.*")
	addr := getenv("HTTP_ADDR", ":8080")
	dsn := getenv("DB_DSN", "postgresql://root@127.0.0.1:26257/app?sslmode=disable")
	jwtSecret := getenv("JWT_SECRET", "devsecret")

	js, err := broker.Connect(natsURL)
	if err != nil {
		log.Fatal(err)
	}
	if err := broker.EnsureStream(js, stream, []string{subjects}); err != nil {
		log.Fatal(err)
	}

	pool, err := db.Open(dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Migrate(pool); err != nil {
		log.Fatal(err)
	}

	authSvc := auth.New(jwtSecret)
	h := httpapi.NewHandler(js, stream, pool, authSvc)
	srv := &http.Server{Addr: addr, Handler: h}
	log.Printf("api listening on %s", addr)
	log.Fatal(srv.ListenAndServe())
}

func getenv(k, d string) string {
	v := os.Getenv(k)
	if v == "" {
		return d
	}
	return v
}
