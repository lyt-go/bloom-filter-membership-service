// Group 领域模型（过滤器分组）。
package model

import (
	"strings"
	"time"
)

type Group struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (x *Group) Validate() error {
	x.Name = strings.TrimSpace(x.Name)
	if x.Name == "" {
		return NewValidationError("name", "名称不能为空")
	}
	x.Description = strings.TrimSpace(x.Description)
	return nil
}

type GroupFilter struct {
	Keyword string
}

func (f GroupFilter) Match(x *Group) bool {
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(x.Name), k) &&
			!strings.Contains(strings.ToLower(x.Description), k) {
			return false
		}
	}
	return true
}
