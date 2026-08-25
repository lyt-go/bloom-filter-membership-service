package buildstate

import "sync"

type State struct {
	Status     string
	Generation int
}
type Store struct {
	mu     sync.Mutex
	states map[string]State
}

func New() *Store                            { return &Store{states: map[string]State{}} }
func (s *Store) Save(id string, state State) { s.mu.Lock(); s.states[id] = state; s.mu.Unlock() }
func (s *Store) Get(id string) State         { s.mu.Lock(); defer s.mu.Unlock(); return s.states[id] }
