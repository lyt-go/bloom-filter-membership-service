// Tag 领域模型（过滤器标签）。
package model

import (
	"strings"
	"time"
)

type Tag struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (x *Tag) Validate() error {
	x.Name = strings.TrimSpace(x.Name)
	if x.Name == "" {
		return NewValidationError("name", "名称不能为空")
	}
	x.Color = strings.TrimSpace(x.Color)
	if x.Color == "" {
		x.Color = "#888888"
	}
	return nil
}

type TagFilter struct {
	Keyword string
}

func (f TagFilter) Match(x *Tag) bool {
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(x.Name), k) &&
			!strings.Contains(strings.ToLower(x.Color), k) {
			return false
		}
	}
	return true
}
