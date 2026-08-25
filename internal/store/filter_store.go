package store

import (
	"bloomfilter/internal/model"
)

func (s *MemoryStore) CreateFilter(x *model.Filter) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.filters {
		if exist.Name == x.Name {
			return ErrConflict
		}
	}
	s.filters[x.ID] = x
	return nil
}

func (s *MemoryStore) GetFilter(id string) (*model.Filter, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	x, ok := s.filters[id]
	if !ok {
		return nil, ErrNotFound
	}
	return x, nil
}

func (s *MemoryStore) GetFilterByName(v string) (*model.Filter, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, x := range s.filters {
		if x.Name == v {
			return x, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListFilters() []*model.Filter {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Filter, 0, len(s.filters))
	for _, x := range s.filters {
		list = append(list, x)
	}
	return list
}

func (s *MemoryStore) UpdateFilter(x *model.Filter) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.filters[x.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.filters {
		if exist.ID != x.ID && exist.Name == x.Name {
			return ErrConflict
		}
	}
	s.filters[x.ID] = x
	return nil
}

func (s *MemoryStore) DeleteFilter(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.filters[id]; !ok {
		return ErrNotFound
	}
	delete(s.filters, id)
	return nil
}
