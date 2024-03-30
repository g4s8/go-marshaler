package marshaler

import "context"

// MapKV is a simple key-value storage implementation based on a map of strings.
type MapKV map[string]string

func (m MapKV) Get(ctx context.Context, key string) (Value, error) {
	val, ok := m[key]
	if !ok {
		return NullValue, nil
	}
	return NewStringValue(val), nil
}
