package filterbuild

import (
	"bloomfilter/internal/filtercache"
	"fmt"
	"strings"
)

type Builder struct{ Cache *filtercache.Cache }

func (b *Builder) Build(id string, parts []string) (draft *filtercache.Draft, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("build failed: %v", r)
		}
	}()
	draft = &filtercache.Draft{ID: id}
	b.Cache.Put(draft)
	for _, part := range parts {
		kv := strings.Split(part, "=")
		if len(kv) != 2 {
			panic("bad segment")
		}
		if kv[0] == "strategy" {
			draft.Strategy = kv[1]
		}
	}
	return draft, nil
}
