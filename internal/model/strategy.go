// Strategy 领域模型（哈希策略）。
package model

import (
	"strings"
	"time"
)

const (
	StrategyStatusEnabled  = "enabled"
	StrategyStatusDisabled = "disabled"

	StrategyAlgoFNV1a = "fnv1a"
	StrategyAlgoFNV1  = "fnv1"
	StrategyAlgoDJB2  = "djb2"
	StrategyAlgoSDBM  = "sdbm"
)

type Strategy struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Algorithm   string    `json:"algorithm"`
	Seed        int64     `json:"seed"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (x *Strategy) Validate() error {
	x.Name = strings.TrimSpace(x.Name)
	if x.Name == "" {
		return NewValidationError("name", "名称不能为空")
	}
	x.Algorithm = strings.TrimSpace(x.Algorithm)
	if !StrategyValidAlgorithm(x.Algorithm) {
		return NewValidationError("algorithm", "算法不合法")
	}
	x.Description = strings.TrimSpace(x.Description)
	if x.Status == "" {
		x.Status = StrategyStatusEnabled
	}
	if !StrategyValidStatus(x.Status) {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

func StrategyValidAlgorithm(s string) bool {
	switch s {
	case StrategyAlgoFNV1a, StrategyAlgoFNV1, StrategyAlgoDJB2, StrategyAlgoSDBM:
		return true
	default:
		return false
	}
}

func StrategyValidStatus(s string) bool {
	switch s {
	case StrategyStatusEnabled, StrategyStatusDisabled:
		return true
	default:
		return false
	}
}

var strategyTransitions = map[string]map[string]bool{
	StrategyStatusEnabled:  {StrategyStatusDisabled: true},
	StrategyStatusDisabled: {StrategyStatusEnabled: true},
}

func StrategyCanTransition(from, to string) bool {
	if m, ok := strategyTransitions[from]; ok {
		return m[to]
	}
	return false
}

type StrategyFilter struct {
	Algorithm string
	Status    string
}

func (f StrategyFilter) Match(x *Strategy) bool {
	if f.Algorithm != "" && x.Algorithm != f.Algorithm {
		return false
	}
	if f.Status != "" && x.Status != f.Status {
		return false
	}
	return true
}
