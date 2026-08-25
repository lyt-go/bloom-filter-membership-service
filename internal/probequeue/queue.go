package probequeue

import "bloomfilter/internal/probeenvelope"

// Queue 保存已提交但尚未消费的探测。
// 每条探测在 Enqueue 时从 Pool 借出一个独立的 Envelope，保存其过滤器和标签；
// 该 Envelope 由队列独占，直到消费时才会归还给 Pool。
// 因此队列中每条探测保留的都是各自提交时的数据，复用也不会串数据。
type Queue struct {
	Pool  probeenvelope.Pool
	items []*probeenvelope.Envelope
}

// Enqueue 将一条探测入队，复用 Pool 中的 Envelope。
// 借出的 Envelope 由队列持有，不在此时归还——否则下一次 Get 可能复用同一对象，
// 覆盖本条已入队的数据。
func (q *Queue) Enqueue(filter string, tags []string) {
	e := q.Pool.Get()
	e.Filter = filter
	e.Tags = append(e.Tags[:0], tags...)
	q.items = append(q.items, e)
}

// At 返回第 n 条探测的只读快照。
// 返回一份独立拷贝而非池中的指针：队列内部对象由 Pool 复用，
// 拷贝时连 Tags 切片一并深拷贝，彻底断开与池中底层数组的别名，
// 调用方拿到本条数据后即使该对象后续被消费、归还、复用，也不会读到被改写的数据。
func (q *Queue) At(n int) probeenvelope.Envelope {
	return snapshot(q.items[n])
}

// Dequeue 取出并消费队首的探测，归还其 Envelope 给 Pool 以便复用。
// 消费完成后该对象不再被队列引用，归还是安全的，不会串到仍在队列中的探测。
// 返回的快照在归还前完成深拷贝，保证归还、复用不影响已取出的数据。
func (q *Queue) Dequeue() (probeenvelope.Envelope, bool) {
	if len(q.items) == 0 {
		return probeenvelope.Envelope{}, false
	}
	e := q.items[0]
	q.items[0] = nil // 释放引用，避免对象被复用后仍被队列持有
	q.items = q.items[1:]
	out := snapshot(e)
	q.Pool.Release(e)
	return out, true
}

// snapshot 生成 e 的一份独立拷贝，Tags 切片使用独立底层数组，
// 与池中对象彻底断开别名。
func snapshot(e *probeenvelope.Envelope) probeenvelope.Envelope {
	out := probeenvelope.Envelope{Filter: e.Filter}
	if len(e.Tags) > 0 {
		tags := make([]string, len(e.Tags))
		copy(tags, e.Tags)
		out.Tags = tags
	}
	return out
}

// Len 返回队列中待消费的探测数量。
func (q *Queue) Len() int { return len(q.items) }
