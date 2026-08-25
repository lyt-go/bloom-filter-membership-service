package rebuildclient

import (
	"context"
	"time"
)

type Client struct{ Delay time.Duration }

func (c Client) Fetch(ctx context.Context) error {
	select {
	case <-time.After(c.Delay):
		return nil
	case <-context.Background().Done():
		return ctx.Err()
	}
}
