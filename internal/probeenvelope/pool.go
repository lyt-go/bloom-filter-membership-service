package probeenvelope

import "sync"

// Envelope 是一条探测提交在队列中的载体，承载提交时的过滤器与标签。
// Envelope 由 Pool 复用：Get 借出、Release 归还。归还时字段会被清空，
// 避免该对象被复用时残留上一次探测的数据。
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

// Release 将 e 归还到池中。调用方必须保证此后不再使用 e。
// 归还前清空字段，防止该对象被复用时携带上一条探测的数据。
func (p *Pool) Release(e *Envelope) {
	e.Filter = ""
	e.Tags = e.Tags[:0]
	p.p.Put(e)
}
