package store

import (
	"bloomfilter/internal/model"
)

func (s *MemoryStore) CreateElement(x *model.Element) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.elements[x.ID] = x
	return nil
}

func (s *MemoryStore) GetElement(id string) (*model.Element, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	x, ok := s.elements[id]
	if !ok {
		return nil, ErrNotFound
	}
	return x, nil
}

func (s *MemoryStore) ListElements() []*model.Element {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Element, 0, len(s.elements))
	for _, x := range s.elements {
		list = append(list, x)
	}
	return list
}

func (s *MemoryStore) DeleteElement(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.elements[id]; !ok {
		return ErrNotFound
	}
	delete(s.elements, id)
	return nil
}
