package store

import (
	"sync"

	"bloomfilter/internal/model"
)

type MemoryStore struct {
	mu         sync.RWMutex
	filters    map[string]*model.Filter
	strategies map[string]*model.Strategy
	elements   map[string]*model.Element
	probes     map[string]*model.Probe
	groups     map[string]*model.Group
	tags       map[string]*model.Tag
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		filters:    make(map[string]*model.Filter),
		strategies: make(map[string]*model.Strategy),
		elements:   make(map[string]*model.Element),
		probes:     make(map[string]*model.Probe),
		groups:     make(map[string]*model.Group),
		tags:       make(map[string]*model.Tag),
	}
}

var _ Store = (*MemoryStore)(nil)
