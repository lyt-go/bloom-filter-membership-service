// Filter 领域模型（布隆过滤器实例）。
package model

import (
	"strings"
	"time"
)

const (
	FilterStatusActive   = "active"
	FilterStatusArchived = "archived"
)

type Filter struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	BitSize           int       `json:"bit_size"`
	HashCount         int       `json:"hash_count"`
	FalsePositiveRate float64   `json:"false_positive_rate"`
	StrategyID        string    `json:"strategy_id"`
	GroupID           string    `json:"group_id"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (x *Filter) Validate() error {
	x.Name = strings.TrimSpace(x.Name)
	if x.Name == "" {
		return NewValidationError("name", "名称不能为空")
	}
	x.Description = strings.TrimSpace(x.Description)
	if x.BitSize <= 0 {
		return NewValidationError("bit_size", "位数组大小必须为正数")
	}
	if x.HashCount <= 0 {
		return NewValidationError("hash_count", "哈希函数个数必须为正数")
	}
	if x.FalsePositiveRate < 0 || x.FalsePositiveRate >= 1 {
		return NewValidationError("false_positive_rate", "误判率必须在 [0,1) 之间")
	}
	if x.Status == "" {
		x.Status = FilterStatusActive
	}
	if !FilterValidStatus(x.Status) {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

func FilterValidStatus(s string) bool {
	switch s {
	case FilterStatusActive, FilterStatusArchived:
		return true
	default:
		return false
	}
}

var filterTransitions = map[string]map[string]bool{
	FilterStatusActive:   {FilterStatusArchived: true},
	FilterStatusArchived: {FilterStatusActive: true},
}

func FilterCanTransition(from, to string) bool {
	if m, ok := filterTransitions[from]; ok {
		return m[to]
	}
	return false
}

type FilterFilter struct {
	Status     string
	StrategyID string
	GroupID    string
	Keyword    string
}

func (f FilterFilter) Match(x *Filter) bool {
	if f.Status != "" && x.Status != f.Status {
		return false
	}
	if f.StrategyID != "" && x.StrategyID != f.StrategyID {
		return false
	}
	if f.GroupID != "" && x.GroupID != f.GroupID {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(x.Name), k) &&
			!strings.Contains(strings.ToLower(x.Description), k) {
			return false
		}
	}
	return true
}
