package activities

import (
	"context"
	"time"
)

func SimulateTimeout(ctx context.Context) error {
	time.Sleep(10 * time.Second)
	return nil
}
