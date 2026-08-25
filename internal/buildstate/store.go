// Package buildstate 维护过滤器构建状态。
//
// 状态写回一律经过带守卫的方法，保证两条不变式：
//   - ready 状态不能倒退：一旦 ready，迟到的 building 回调无法将其改回 building；
//   - 版本（generation）单调递增：旧版本的回调无法覆盖新版本。
package buildstate

import "sync"

type State struct {
	Status     string
	Generation int
}

type Store struct {
	mu     sync.Mutex
	states map[string]State
}

func New() *Store { return &Store{states: map[string]State{}} }

// Get 返回过滤器当前的构建状态；未构建过时返回零值。
func (s *Store) Get(id string) State {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.states[id]
}

// ApplyBuilding 记录某 generation 的构建进行中（building）。
//
// 当回调会导致状态倒退时直接丢弃并返回 false：
//   - 当前已 ready：ready 不得退回 building（无论 generation 如何）；
//   - 回调的 generation 落后于当前：旧版本不得覆盖新版本。
//
// 检查与写入在同一把锁内完成，晚到的回调无法在 Get/Save 之间钻空子。
func (s *Store) ApplyBuilding(id string, generation int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	cur := s.states[id]
	if cur.Status == "ready" || generation < cur.Generation {
		return false
	}
	s.states[id] = State{Status: "building", Generation: generation}
	return true
}

// PromoteToReady 将过滤器标记为某 generation 的 ready。
//
// 当 generation 会倒退（落后于当前）时返回 false，避免重试用旧版本覆盖更新的 ready 状态。
func (s *Store) PromoteToReady(id string, generation int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	cur := s.states[id]
	if generation < cur.Generation {
		return false
	}
	s.states[id] = State{Status: "ready", Generation: generation}
	return true
}
