package store

import (
	"bloomfilter/internal/model"
)

func (s *MemoryStore) CreateProbe(x *model.Probe) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.probes[x.ID] = x
	return nil
}

func (s *MemoryStore) GetProbe(id string) (*model.Probe, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	x, ok := s.probes[id]
	if !ok {
		return nil, ErrNotFound
	}
	return x, nil
}

func (s *MemoryStore) ListProbes() []*model.Probe {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Probe, 0, len(s.probes))
	for _, x := range s.probes {
		list = append(list, x)
	}
	return list
}

func (s *MemoryStore) DeleteProbe(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.probes[id]; !ok {
		return ErrNotFound
	}
	delete(s.probes, id)
	return nil
}
