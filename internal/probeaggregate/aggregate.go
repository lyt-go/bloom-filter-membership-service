package probeaggregate

import (
	"bloomfilter/internal/segmentstream"
	"context"
)

func Collect(ctx context.Context, parts []string) ([]string, error) {
	data, _ := segmentstream.Start(parts)
	var out []string
	for {
		select {
		case v, ok := <-data:
			if !ok {
				return out, nil
			}
			out = append(out, v)
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}
