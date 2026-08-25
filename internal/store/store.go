// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"bloomfilter/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	CreateFilter(x *model.Filter) error
	GetFilter(id string) (*model.Filter, error)
	GetFilterByName(v string) (*model.Filter, error)
	ListFilters() []*model.Filter
	UpdateFilter(x *model.Filter) error
	DeleteFilter(id string) error
	CreateStrategy(x *model.Strategy) error
	GetStrategy(id string) (*model.Strategy, error)
	GetStrategyByName(v string) (*model.Strategy, error)
	ListStrategies() []*model.Strategy
	UpdateStrategy(x *model.Strategy) error
	DeleteStrategy(id string) error
	CreateElement(x *model.Element) error
	GetElement(id string) (*model.Element, error)
	ListElements() []*model.Element
	DeleteElement(id string) error
	CreateProbe(x *model.Probe) error
	GetProbe(id string) (*model.Probe, error)
	ListProbes() []*model.Probe
	DeleteProbe(id string) error
	CreateGroup(x *model.Group) error
	GetGroup(id string) (*model.Group, error)
	GetGroupByName(v string) (*model.Group, error)
	ListGroups() []*model.Group
	UpdateGroup(x *model.Group) error
	DeleteGroup(id string) error
	CreateTag(x *model.Tag) error
	GetTag(id string) (*model.Tag, error)
	GetTagByName(v string) (*model.Tag, error)
	ListTags() []*model.Tag
	UpdateTag(x *model.Tag) error
	DeleteTag(id string) error
}
