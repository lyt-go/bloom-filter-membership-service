package rebuildflow

import (
	"bloomfilter/internal/rebuildclient"
	"context"
)

type Flow struct{ bound context.Context }

func (f *Flow) Rebuild(ctx context.Context, c rebuildclient.Client) error {
	if f.bound == nil {
		f.bound = ctx
	}
	if err := f.bound.Err(); err != nil {
		return err
	}
	return c.Fetch(f.bound)
}
