package probequeue

import (
	"reflect"
	"testing"
)

func TestQueuedProbeSurvivesEnvelopeReuse(t *testing.T) {
	var q Queue
	q.Enqueue("filter-alpha", []string{"edge", "login"})
	q.Enqueue("filter-beta", []string{"worker"})
	first := q.At(0)
	if first.Filter != "filter-alpha" {
		t.Fatalf("first filter changed to %q", first.Filter)
	}
	if !reflect.DeepEqual(first.Tags, []string{"edge", "login"}) {
		t.Fatalf("first tags changed to %v", first.Tags)
	}
}
