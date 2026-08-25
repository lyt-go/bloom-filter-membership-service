package service

import (
	"sort"
	"time"

	"bloomfilter/internal/model"
	"bloomfilter/pkg/idgen"
)

func (s *Service) CreateElement(input model.Element) (*model.Element, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetFilter(input.FilterID); err != nil {
		return nil, model.NewValidationError("filter_id", "关联的过滤器不存在")
	}
	now := time.Now()
	input.ID = idgen.Hex()
	input.CreatedAt = now
	if err := s.store.CreateElement(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetElement(id string) (*model.Element, error) {
	return s.store.GetElement(id)
}

func (s *Service) ListElements(filter model.ElementFilter, page, size int) ([]*model.Element, int, error) {
	all := s.store.ListElements()
	matched := make([]*model.Element, 0, len(all))
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
		return []*model.Element{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) DeleteElement(id string) error {
	return s.store.DeleteElement(id)
}
