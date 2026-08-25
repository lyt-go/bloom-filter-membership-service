package buildflow

import (
	"bloomfilter/internal/buildstate"
	"testing"
)

func TestRetrySuccessRejectsStaleBuildAndDuplicateEffect(t *testing.T) {
	store := buildstate.New()
	flow := New(store)
	flow.Retry("filter-a")
	flow.Callback("filter-a", 1)
	state := store.Get("filter-a")
	if state.Status != "ready" || state.Generation != 2 {
		t.Fatalf("stale callback changed state: %+v", state)
	}
	if flow.Calls != 1 {
		t.Fatalf("external build calls=%d", flow.Calls)
	}
}
