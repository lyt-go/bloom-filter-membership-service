package service

import (
	"sort"
	"time"

	"bloomfilter/internal/model"
	"bloomfilter/pkg/idgen"
)

func (s *Service) CreateTag(input model.Tag) (*model.Tag, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetTagByName(input.Name); err == nil {
		return nil, model.NewValidationError("name", "已存在")
	}
	now := time.Now()
	input.ID = idgen.Hex()
	input.CreatedAt = now
	input.UpdatedAt = now
	if err := s.store.CreateTag(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetTag(id string) (*model.Tag, error) {
	return s.store.GetTag(id)
}

func (s *Service) ListTags(filter model.TagFilter, page, size int) ([]*model.Tag, int, error) {
	all := s.store.ListTags()
	matched := make([]*model.Tag, 0, len(all))
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
		return []*model.Tag{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateTag(id string, input model.Tag) (*model.Tag, error) {
	exist, err := s.store.GetTag(id)
	if err != nil {
		return nil, err
	}
	exist.Name = input.Name
	exist.Color = input.Color
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateTag(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

func (s *Service) DeleteTag(id string) error {
	return s.store.DeleteTag(id)
}
