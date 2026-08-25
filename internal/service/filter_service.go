package service

import (
	"sort"
	"time"

	"bloomfilter/internal/model"
	"bloomfilter/pkg/idgen"
)

func (s *Service) CreateFilter(input model.Filter) (*model.Filter, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetFilterByName(input.Name); err == nil {
		return nil, model.NewValidationError("name", "已存在")
	}
	if input.StrategyID != "" {
		if _, err := s.store.GetStrategy(input.StrategyID); err != nil {
			return nil, model.NewValidationError("strategy_id", "关联的哈希策略不存在")
		}
	}
	if input.GroupID != "" {
		if _, err := s.store.GetGroup(input.GroupID); err != nil {
			return nil, model.NewValidationError("group_id", "关联的分组不存在")
		}
	}
	now := time.Now()
	input.ID = idgen.Hex()
	input.CreatedAt = now
	input.UpdatedAt = now
	if err := s.store.CreateFilter(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetFilter(id string) (*model.Filter, error) {
	return s.store.GetFilter(id)
}

func (s *Service) ListFilters(filter model.FilterFilter, page, size int) ([]*model.Filter, int, error) {
	all := s.store.ListFilters()
	matched := make([]*model.Filter, 0, len(all))
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
		return []*model.Filter{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateFilter(id string, input model.Filter) (*model.Filter, error) {
	exist, err := s.store.GetFilter(id)
	if err != nil {
		return nil, err
	}
	exist.Name = input.Name
	exist.Description = input.Description
	exist.BitSize = input.BitSize
	exist.HashCount = input.HashCount
	exist.FalsePositiveRate = input.FalsePositiveRate
	exist.StrategyID = input.StrategyID
	exist.GroupID = input.GroupID
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	if input.StrategyID != "" {
		if _, err := s.store.GetStrategy(input.StrategyID); err != nil {
			return nil, model.NewValidationError("strategy_id", "关联的哈希策略不存在")
		}
	}
	if input.GroupID != "" {
		if _, err := s.store.GetGroup(input.GroupID); err != nil {
			return nil, model.NewValidationError("group_id", "关联的分组不存在")
		}
	}
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateFilter(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

func (s *Service) DeleteFilter(id string) error {
	return s.store.DeleteFilter(id)
}

func (s *Service) TransitionFilter(id, target string) (*model.Filter, error) {
	if !model.FilterValidStatus(target) {
		return nil, model.NewValidationError("status", "目标状态不合法")
	}
	exist, err := s.store.GetFilter(id)
	if err != nil {
		return nil, err
	}
	if !model.FilterCanTransition(exist.Status, target) {
		return nil, model.NewValidationError("status", "不允许从 "+exist.Status+" 流转到 "+target)
	}
	exist.Status = target
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateFilter(exist); err != nil {
		return nil, err
	}
	return exist, nil
}
