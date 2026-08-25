package refreshdispatch

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

type failingRefresher struct{ calls atomic.Int64 }

func (r *failingRefresher) Refresh() error { r.calls.Add(1); return errors.New("offline") }
func TestCancelStopsRefreshRetriesAndShutdown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	r := &failingRefresher{}
	var d Dispatcher
	d.Start(ctx, r)
	deadline := time.After(200 * time.Millisecond)
	for r.calls.Load() < 2 {
		select {
		case <-deadline:
			t.Fatal("refresh did not start")
		default:
			time.Sleep(time.Millisecond)
		}
	}
	cancel()
	before := r.calls.Load()
	time.Sleep(40 * time.Millisecond)
	if after := r.calls.Load(); after != before {
		t.Fatalf("calls grew after cancel: %d to %d", before, after)
	}
	done := make(chan struct{})
	go func() { d.Shutdown(); close(done) }()
	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("shutdown timed out")
	}
}
