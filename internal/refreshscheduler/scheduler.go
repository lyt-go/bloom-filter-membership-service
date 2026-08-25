package refreshscheduler

import (
	"context"
	"time"
)

type Refresher interface{ Refresh() error }

func Run(ctx context.Context, r Refresher) {
	for {
		if r.Refresh() == nil {
			return
		}
		time.Sleep(15 * time.Millisecond)
	}
}
