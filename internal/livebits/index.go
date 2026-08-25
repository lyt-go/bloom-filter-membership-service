package livebits

import "sync"

type Index struct {
	mu   sync.RWMutex
	data map[string]map[string]int
}

func New() *Index { return &Index{data: map[string]map[string]int{"filter-a": {"initial": 1}}} }
func (i *Index) Put(filter, group string, n int) {
	i.mu.Lock()
	if i.data[filter] == nil {
		i.data[filter] = map[string]int{}
	}
	i.data[filter][group] = n
	i.mu.Unlock()
}
func (i *Index) Snapshot() map[string]map[string]int {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.data
}
