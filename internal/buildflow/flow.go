package buildflow

import (
	"bloomfilter/internal/buildstate"
	"fmt"
	"sync"
)

type Flow struct {
	Store   *buildstate.Store
	mu      sync.Mutex
	effects map[string]bool
	Calls   int
}

func New(store *buildstate.Store) *Flow { return &Flow{Store: store, effects: map[string]bool{}} }
func (f *Flow) effect(id string, generation int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	key := fmt.Sprintf("%s-%d", id, generation)
	if !f.effects[key] {
		f.effects[key] = true
		f.Calls++
	}
}
func (f *Flow) Retry(id string) {
	f.Store.Save(id, buildstate.State{Status: "ready", Generation: 2})
	f.effect(id, 2)
}
func (f *Flow) Callback(id string, generation int) {
	f.Store.Save(id, buildstate.State{Status: "building", Generation: generation})
	f.effect(id, generation)
}
