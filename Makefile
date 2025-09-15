run:
	docker compose up --build

api:
	NATS_URL=${NATS_URL} HTTP_ADDR=${HTTP_ADDR} STREAM_NAME=${STREAM_NAME} STREAM_SUBJECTS=${STREAM_SUBJECTS} go run cmd/api/main.go

worker:
	NATS_URL=${NATS_URL} STREAM_NAME=${STREAM_NAME} WORKER_TYPE=${WORKER_TYPE} CONSUMER_NAME=${CONSUMER_NAME} go run cmd/worker/main.go

test:
	go test ./...
