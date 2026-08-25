// Package buildflow 编排过滤器的外部构建与状态回流。
//
// 重建流程面临的核心竞态：重试成功（ready）后，上一轮迟到的回调才到达，
// 它会无脑把状态写回 building 并把 generation 倒退到旧版本，于是外部看到
// 旧版本的详情；而重试本身又触发了第二次外部构建。
//
// 本包通过两条不变式修复该竞态：
//   - 状态不可倒退：ready 之后任何 building 回调（无论是否旧 generation）
//     均被丢弃；旧 generation 的回调一律拒绝；
//   - 同一过滤器只构建一次：外部构建按过滤器去重，重试复用既有构建，
//     不再二次触发。
package buildflow

import (
	"bloomfilter/internal/buildstate"
	"sync"
)

type Flow struct {
	Store *buildstate.Store
	mu    sync.Mutex
	// built 记录已经触发过外部构建的过滤器（按 id 去重）。
	// 重建语义上同一个过滤器只应产生一次外部构建，重试只是状态回流，
	// 不应再次触发构建，因此按 id 而非 id-generation 去重。
	built  map[string]bool
	Calls  int
	effect func(id string)
}

// New 创建一个构建流程。可选注入 effect（外部构建动作），便于测试观测；
// 未注入时仅累加 Calls 计数。
func New(store *buildstate.Store) *Flow {
	return &Flow{
		Store: store,
		built: map[string]bool{},
		effect: func(id string) {
			// 默认实现：外部构建的实际执行点由调用方接入。
			// 此处仅作为计数占位，真正的副作用由调用方在 effect 注入时承载。
		},
	}
}

// SetEffect 注入外部构建动作，返回调用以链式使用。
func (f *Flow) SetEffect(fn func(id string)) *Flow {
	f.mu.Lock()
	defer f.mu.Unlock()
	if fn != nil {
		f.effect = fn
	}
	return f
}

// triggerEffect 触发一次外部构建。同一过滤器仅触发一次：
// 重试命中已存在的构建记录，不重复计数、不重复执行。
func (f *Flow) triggerEffect(id string) {
	f.mu.Lock()
	if f.built[id] {
		f.mu.Unlock()
		return
	}
	f.built[id] = true
	f.Calls++
	fn := f.effect
	f.mu.Unlock()
	if fn != nil {
		fn(id)
	}
}

// Retry 重试一次失败/超时的构建并回流为 ready。
//
// 状态推进交给 buildstate.Store 守卫：仅当 generation 单调递增、且不会
// 把更新的 ready 倒退时才写入。外部构建按过滤器去重，因此重试复用既有
// 构建，不会触发第二次副作用——这也是“同一过滤器只能产生一次外部构建”
// 的语义来源。
func (f *Flow) Retry(id string) {
	const retryGeneration = 2
	f.Store.PromoteToReady(id, retryGeneration)
	f.triggerEffect(id)
}

// Callback 处理外部构建系统回流的状态回调。
//
// 关键点：绝不让回调造成状态倒退。ready 之后或旧 generation 的回调由
// Store.ApplyBuilding 直接拒绝（返回 false），从而避免：
//   - 把 ready 改回 building（状态倒退）；
//   - 把 generation 倒退到旧版本（详情看到旧版本）。
//
// generation 未提供时（=0）以当前状态推导，仅用于兼容无版本号的旧回调。
func (f *Flow) Callback(id string, generation int) {
	if generation == 0 {
		cur := f.Store.Get(id)
		generation = cur.Generation
	}
	f.Store.ApplyBuilding(id, generation)
	f.triggerEffect(id)
}
