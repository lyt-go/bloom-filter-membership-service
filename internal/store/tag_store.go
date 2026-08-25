package store

import (
	"bloomfilter/internal/model"
)

func (s *MemoryStore) CreateTag(x *model.Tag) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.tags {
		if exist.Name == x.Name {
			return ErrConflict
		}
	}
	s.tags[x.ID] = x
	return nil
}

func (s *MemoryStore) GetTag(id string) (*model.Tag, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	x, ok := s.tags[id]
	if !ok {
		return nil, ErrNotFound
	}
	return x, nil
}

func (s *MemoryStore) GetTagByName(v string) (*model.Tag, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, x := range s.tags {
		if x.Name == v {
			return x, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListTags() []*model.Tag {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Tag, 0, len(s.tags))
	for _, x := range s.tags {
		list = append(list, x)
	}
	return list
}

func (s *MemoryStore) UpdateTag(x *model.Tag) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tags[x.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.tags {
		if exist.ID != x.ID && exist.Name == x.Name {
			return ErrConflict
		}
	}
	s.tags[x.ID] = x
	return nil
}

func (s *MemoryStore) DeleteTag(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tags[id]; !ok {
		return ErrNotFound
	}
	delete(s.tags, id)
	return nil
}
