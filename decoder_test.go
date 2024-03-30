package marshaler

import (
	"context"
	"testing"
	"time"
)

type kvStub struct {
	data map[string]string
}

func newKVStub() *kvStub {
	return &kvStub{data: make(map[string]string)}
}

func (k *kvStub) Set(key, value string) error {
	k.data[key] = value
	return nil
}

func (k *kvStub) With(key, value string) *kvStub {
	k.data[key] = value
	return k
}

func (k *kvStub) Get(ctx context.Context, key string) (Value, error) {
	val, ok := k.data[key]
	if !ok {
		return NullValue, nil
	}
	return NewStringValue(val), nil
}

func TestDecoder(t *testing.T) {
	t.Run("Decode", func(t *testing.T) {
		target := struct {
			Foo string `kv:"foo"`
			Bar int    `kv:"bar"`
			Baz struct {
				Qwe string `kv:"qwe"`
				Asd bool   `kv:"asd"`
				Zxc *struct {
					Num int           `kv:"num"`
					Dur time.Duration `kv:"dur"`
					Bin *scanner      `kv:"bin"`
				} `kv:"zxc"`
			} `kv:"baz"`
		}{}
		kv := newKVStub().
			With("foo", "bar").
			With("bar", "42").
			With("baz/qwe", "qwe").
			With("baz/asd", "true").
			With("baz/zxc/num", "42").
			With("baz/zxc/dur", "1s").
			With("baz/zxc/bin", "hello")
		dec, err := NewDecoder(kv)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := dec.Decode(&target); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if target.Foo != "bar" {
			t.Fatalf("expected %q, got %q", "bar", target.Foo)
		}
		if target.Bar != 42 {
			t.Fatalf("expected %d, got %d", 42, target.Bar)
		}
		if target.Baz.Qwe != "qwe" {
			t.Fatalf("expected %q, got %q", "qwe", target.Baz.Qwe)
		}
		if !target.Baz.Asd {
			t.Fatalf("expected %t, got %t", true, target.Baz.Asd)
		}
		if target.Baz.Zxc == nil {
			t.Fatalf("expected non-nil, got nil")
		}
		if target.Baz.Zxc.Num != 42 {
			t.Fatalf("expected %d, got %d", 42, target.Baz.Zxc.Num)
		}
		if target.Baz.Zxc.Dur != time.Second {
			t.Fatalf("expected %v, got %v", time.Second, target.Baz.Zxc.Dur)
		}
		if target.Baz.Zxc.Bin.value != "hello" {
			t.Fatalf("expected %q, got %q", "hello", target.Baz.Zxc.Bin.value)
		}
	})
}
