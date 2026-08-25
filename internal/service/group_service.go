package service

import (
	"sort"
	"time"

	"bloomfilter/internal/model"
	"bloomfilter/pkg/idgen"
)

func (s *Service) CreateGroup(input model.Group) (*model.Group, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetGroupByName(input.Name); err == nil {
		return nil, model.NewValidationError("name", "已存在")
	}
	now := time.Now()
	input.ID = idgen.Hex()
	input.CreatedAt = now
	input.UpdatedAt = now
	if err := s.store.CreateGroup(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetGroup(id string) (*model.Group, error) {
	return s.store.GetGroup(id)
}

func (s *Service) ListGroups(filter model.GroupFilter, page, size int) ([]*model.Group, int, error) {
	all := s.store.ListGroups()
	matched := make([]*model.Group, 0, len(all))
	for _, x := range all {
		if filter.Match(x) {
			matched = append(matched, x)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Group{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateGroup(id string, input model.Group) (*model.Group, error) {
	exist, err := s.store.GetGroup(id)
	if err != nil {
		return nil, err
	}
	exist.Name = input.Name
	exist.Description = input.Description
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateGroup(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

func (s *Service) DeleteGroup(id string) error {
	return s.store.DeleteGroup(id)
}
