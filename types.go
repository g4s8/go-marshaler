package marshaler

import (
	"context"
)

// KV is a key-value storage API.
//
// It provides a method to get a value by key.
type KV interface {
	Get(ctx context.Context, key string) (Value, error)
}

// Value is a kv-storage value, that can be unmarshaled to a target type.
//
// KV implementation can use provided Value implementation to unmarshal,
// like [StringValue](#StringValue) or implement custom Value type.
type Value interface {
	// UnmarshalTo unmarshal value to a target type or return an error
	// if it's not possible.
	UnmarshalTo(out any, opts ValueUnmarshalOpts) error
}

// Scanner is an interface for scanning values.
//
// It can scan a value by itself. The scan method could be called by [Value.UnmarshalTo](#Value.UnmarshalTo)
// implementation to scan a value to a target type.
//
// Example:
//
//	type Custom struct {
//		Foo string
//		Bar int
//	}
//
//	func (c *Custom) Scan(v any) error {
//		val, ok := v.(string)
//		if !ok {
//			return errors.New("unexpected type")
//		}
//		parts := strings.Split(val, ",")
//		c.Foo = parts[0]
//		i, err := strconv.Atoi(parts[1])
//		if err != nil {
//			return err
//		}
//		c.Bar = i
//		return nil
//	}
type Scanner interface {
	// Scan scans a value to a target type.
	Scan(v any) error
}

// ValueUnmarshalOpts is a set of options for value unmarshaling.
type ValueUnmarshalOpts struct {
	// SliceSep is a separator for slice values.
	SliceSep string
}
