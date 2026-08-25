// Probe 领域模型（成员检测记录）。
package model

import (
	"strings"
	"time"
)

const (
	ProbeResultPresent = "present"
	ProbeResultAbsent  = "absent"
)

type Probe struct {
	ID            string    `json:"id"`
	FilterID      string    `json:"filter_id"`
	Value         string    `json:"value"`
	Result        string    `json:"result"`
	FalsePositive bool      `json:"false_positive"`
	CreatedAt     time.Time `json:"created_at"`
}

func (x *Probe) Validate() error {
	x.FilterID = strings.TrimSpace(x.FilterID)
	if x.FilterID == "" {
		return NewValidationError("filter_id", "过滤器 ID 不能为空")
	}
	x.Value = strings.TrimSpace(x.Value)
	if x.Value == "" {
		return NewValidationError("value", "元素值不能为空")
	}
	if x.Result == "" {
		x.Result = ProbeResultAbsent
	}
	if !ProbeValidResult(x.Result) {
		return NewValidationError("result", "检测结果不合法")
	}
	return nil
}

func ProbeValidResult(s string) bool {
	switch s {
	case ProbeResultPresent, ProbeResultAbsent:
		return true
	default:
		return false
	}
}

type ProbeFilter struct {
	FilterID string
	Result   string
	Keyword  string
}

func (f ProbeFilter) Match(x *Probe) bool {
	if f.FilterID != "" && x.FilterID != f.FilterID {
		return false
	}
	if f.Result != "" && x.Result != f.Result {
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
