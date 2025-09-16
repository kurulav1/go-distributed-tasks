package jobs

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Store(ctx context.Context, db *pgxpool.Pool, j Job) error {
	b, err := json.Marshal(j.Payload)
	if err != nil {
		return err
	}
	_, err = db.Exec(ctx, `
insert into jobs(id,user_id,type,payload,enqueued_at)
values($1,$2,$3,$4,$5)
`, j.ID, j.UserID, j.Type, b, j.EnqueuedAt)
	return err
}

func List(ctx context.Context, db *pgxpool.Pool, userID string, limit int) ([]Job, error) {
	rows, err := db.Query(ctx, `
select id,user_id,type,payload,enqueued_at
from jobs
where user_id=$1
order by enqueued_at desc
limit $2
`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []Job
	for rows.Next() {
		var j Job
		var pb []byte
		var t time.Time
		if err := rows.Scan(&j.ID, &j.UserID, &j.Type, &pb, &t); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(pb, &j.Payload); err != nil {
			return nil, err
		}
		j.EnqueuedAt = t
		res = append(res, j)
	}
	return res, rows.Err()
}
