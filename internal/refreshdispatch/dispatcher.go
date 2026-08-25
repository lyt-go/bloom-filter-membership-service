package refreshdispatch

import (
	"context"
	"errors"
	"sync"
	"time"

	"bloomfilter/internal/refreshscheduler"
)

// ErrShutdownTimeout is returned by Shutdown when a refresh goroutine did not
// stop within the allotted timeout. The goroutine is left to exit on its own.
var ErrShutdownTimeout = errors.New("refresh dispatcher shutdown timed out")

// shutdownTimeout caps how long Shutdown waits for an in-flight refresh to
// stop after canceling its context. It is short so the close path never gets
// stuck waiting on a runaway retry loop.
const shutdownTimeout = 5 * time.Second

// Dispatcher runs a Refresher on a background goroutine that can be stopped
// both by canceling the context passed to Start and by calling Shutdown.
type Dispatcher struct {
	wg     sync.WaitGroup
	cancel context.CancelFunc // nil until Start attaches a derived context
}

// Start launches a goroutine that runs r until it succeeds or ctx is
// canceled. The caller's ctx is honored (canceling it stops the loop), and
// Shutdown can stop it independently. Safe to call once per Dispatcher.
func (d *Dispatcher) Start(ctx context.Context, r refreshscheduler.Refresher) {
	// Derive a child context so Shutdown can cancel the loop even when the
	// caller's ctx has no cancel handle of its own.
	ctx, cancel := context.WithCancel(ctx)
	d.cancel = cancel

	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		_ = refreshscheduler.Run(ctx, r)
	}()
}

// Shutdown signals the running refresh loop to stop and waits for it to exit.
// It cancels the context (so a retry loop stops promptly) and bounds the wait
// by shutdownTimeout, returning ErrShutdownTimeout if the goroutine is still
// running when the deadline elapses.
func (d *Dispatcher) Shutdown() error {
	if d.cancel != nil {
		d.cancel()
	}

	done := make(chan struct{})
	go func() { d.wg.Wait(); close(done) }()

	select {
	case <-done:
		return nil
	case <-time.After(shutdownTimeout):
		return ErrShutdownTimeout
	}
}
