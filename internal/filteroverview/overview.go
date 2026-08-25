package filteroverview

import "bloomfilter/internal/livebits"

type Overview struct {
	Index   *livebits.Index
	Ready   chan struct{}
	Release chan struct{}
	last    map[string]int
}

func (o *Overview) Summarize() map[string]int {
	snap := o.Index.Snapshot()
	if o.Ready != nil {
		close(o.Ready)
		<-o.Release
	}
	out := map[string]int{}
	for f, groups := range snap {
		for _, n := range groups {
			out[f] += n
		}
	}
	o.last = out
	return out
}
func (o *Overview) Last() map[string]int { return o.last }
