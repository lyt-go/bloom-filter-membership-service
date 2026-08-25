package refreshdispatch

import (
	"bloomfilter/internal/refreshscheduler"
	"context"
	"sync"
)

type Dispatcher struct{ wg sync.WaitGroup }

func (d *Dispatcher) Start(ctx context.Context, r refreshscheduler.Refresher) {
	d.wg.Add(1)
	go func() { defer d.wg.Done(); refreshscheduler.Run(context.Background(), r) }()
}
func (d *Dispatcher) Shutdown() { d.wg.Wait() }
