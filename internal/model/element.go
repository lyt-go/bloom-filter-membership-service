// Element 领域模型（已登记元素）。
package model

import (
	"strings"
	"time"
)

type Element struct {
	ID        string    `json:"id"`
	FilterID  string    `json:"filter_id"`
	Value     string    `json:"value"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
}

func (x *Element) Validate() error {
	x.FilterID = strings.TrimSpace(x.FilterID)
	if x.FilterID == "" {
		return NewValidationError("filter_id", "过滤器 ID 不能为空")
	}
	x.Value = strings.TrimSpace(x.Value)
	if x.Value == "" {
		return NewValidationError("value", "元素值不能为空")
	}
	x.Category = strings.TrimSpace(x.Category)
	return nil
}

type ElementFilter struct {
	FilterID string
	Category string
	Keyword  string
}

func (f ElementFilter) Match(x *Element) bool {
	if f.FilterID != "" && x.FilterID != f.FilterID {
		return false
	}
	if f.Category != "" && x.Category != f.Category {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(x.Value), k) {
			return false
		}
	}
	return true
}
