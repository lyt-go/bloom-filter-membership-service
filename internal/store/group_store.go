package store

import (
	"bloomfilter/internal/model"
)

func (s *MemoryStore) CreateGroup(x *model.Group) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.groups {
		if exist.Name == x.Name {
			return ErrConflict
		}
	}
	s.groups[x.ID] = x
	return nil
}

func (s *MemoryStore) GetGroup(id string) (*model.Group, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	x, ok := s.groups[id]
	if !ok {
		return nil, ErrNotFound
	}
	return x, nil
}

func (s *MemoryStore) GetGroupByName(v string) (*model.Group, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, x := range s.groups {
		if x.Name == v {
			return x, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListGroups() []*model.Group {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Group, 0, len(s.groups))
	for _, x := range s.groups {
		list = append(list, x)
	}
	return list
}

func (s *MemoryStore) UpdateGroup(x *model.Group) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.groups[x.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.groups {
		if exist.ID != x.ID && exist.Name == x.Name {
			return ErrConflict
		}
	}
	s.groups[x.ID] = x
	return nil
}

func (s *MemoryStore) DeleteGroup(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.groups[id]; !ok {
		return ErrNotFound
	}
	delete(s.groups, id)
	return nil
}
