package probeenvelope

import "sync"

type Envelope struct {
	Filter string
	Tags   []string
}
type Pool struct{ p sync.Pool }

func (p *Pool) Get() *Envelope {
	v := p.p.Get()
	if v == nil {
		return &Envelope{}
	}
	return v.(*Envelope)
}
func (p *Pool) Release(e *Envelope) { p.p.Put(e) }
