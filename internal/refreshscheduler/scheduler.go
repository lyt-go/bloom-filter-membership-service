package refreshscheduler

import (
	"context"
	"time"
)

// ErrCanceled indicates Run returned because its context was canceled.
var ErrCanceled = context.Canceled

// Refresher performs a single refresh attempt. A nil error means the refresh
// succeeded and no further attempts are needed.
type Refresher interface{ Refresh() error }

// Run drives r.Refresh() until it either succeeds or ctx is canceled. It
// honors cancellation between attempts (and during the backoff sleep), so a
// canceled request stops retrying promptly instead of looping forever.
func Run(ctx context.Context, r Refresher) error {
	for {
		// Bail out before doing any work if we were already canceled; this is
		// what stops the retry loop once a refresh request is canceled.
		if err := ctx.Err(); err != nil {
			return err
		}

		if r.Refresh() == nil {
			return nil
		}

		// Back off, but wake up immediately if the context is canceled so the
		// loop exits without waiting the full sleep.
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(15 * time.Millisecond):
		}
	}
}
