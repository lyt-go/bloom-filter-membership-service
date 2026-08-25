package membershipbatch

import "bloomfilter/internal/hashdecode"

type Archive struct {
	decoder hashdecode.Decoder
	batches map[string][]string
	sent    [][]string
}

func New() *Archive { return &Archive{batches: map[string][]string{}} }
func (a *Archive) Ingest(id, raw string) {
	values := a.decoder.Decode(raw)
	a.batches[id] = values
	a.sent = append(a.sent, values)
}
func (a *Archive) Snapshot(id string) []string { return a.batches[id] }
func (a *Archive) Sent(n int) []string         { return a.sent[n] }
