package consul

import (
	"context"
	"fmt"

	"github.com/g4s8/go-marshaler"
	capi "github.com/hashicorp/consul/api"
)

var _ marshaler.KV = (*consulKV)(nil)

type consulKV struct {
	ckv *capi.KV
}

func (kv *consulKV) Get(ctx context.Context, key string) (marshaler.Value, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	pair, _, err := kv.ckv.Get(key, nil)
	if err != nil {
		return nil, fmt.Errorf("get key %q: %w", key, err)
	}
	if pair == nil {
		return marshaler.NullValue, nil
	}

	return marshaler.NewBytesValue(pair.Value), nil
}
