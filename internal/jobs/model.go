package jobs

import "time"

type Job struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Payload    map[string]interface{} `json:"payload"`
	EnqueuedAt time.Time              `json:"enqueued_at"`
}
