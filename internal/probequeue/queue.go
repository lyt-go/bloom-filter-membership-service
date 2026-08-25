package probequeue

import "bloomfilter/internal/probeenvelope"

type Queue struct {
	Pool  probeenvelope.Pool
	items []*probeenvelope.Envelope
}

func (q *Queue) Enqueue(filter string, tags []string) {
	e := q.Pool.Get()
	e.Filter = filter
	e.Tags = append(e.Tags[:0], tags...)
	q.items = append(q.items, e)
	q.Pool.Release(e)
}
func (q *Queue) At(n int) *probeenvelope.Envelope { return q.items[n] }
