package jobs

import (
	"encoding/json"
	"fmt"
	"time"
)

func Process(b []byte) error {
	var j Job
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}
	switch j.Type {
	case "email":
		time.Sleep(200 * time.Millisecond)
		fmt.Println("processed email", j.ID)
	case "image":
		time.Sleep(300 * time.Millisecond)
		fmt.Println("processed image", j.ID)
	default:
		time.Sleep(100 * time.Millisecond)
		fmt.Println("processed generic", j.ID)
	}
	return nil
}
