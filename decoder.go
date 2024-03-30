package marshaler

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

var defaultConfig = decoderConfig{
	separator: "/",
	sliceSep:  ",",
	tag:       "kv",
}

var (
	// Errors for decoder options validation.
	ErrEmptyKeySep   = fmt.Errorf("empty key separator")
	ErrEmptySliceSep = fmt.Errorf("empty slice separator")
	ErrorEmptyTag    = fmt.Errorf("empty tag name")
)

type decoderConfig struct {
	separator string
	sliceSep  string
	tag       string
	prefix    string
}

// DecoderOption is an option for decoder configuration.
//
// It returns an error if the option parameter is invalid.
type DecoderOption func(*decoderConfig) error

// WithSeparator sets the key separator for decoder.
//
// Default is "/".
func WithSeparator(separator string) DecoderOption {
	return func(c *decoderConfig) error {
		if separator == "" {
			return ErrEmptyKeySep
		}
		c.separator = separator
		return nil
	}
}

// WithSliceSeparator sets the separator of string values for slice fields.
//
// Default is ",".
func WithSliceSeparator(separator string) DecoderOption {
	return func(d *decoderConfig) error {
		if separator == "" {
			return ErrEmptySliceSep
		}
		d.sliceSep = separator
		return nil
	}
}

// WithTag sets the tag name for decoder.
//
// Default is "kv".
func WithTag(tag string) DecoderOption {
	return func(d *decoderConfig) error {
		if tag == "" {
			return ErrorEmptyTag
		}
		d.tag = tag
		return nil
	}
}

// WithPrefix sets the prefix for all keys.
func WithPrefix(prefix string) DecoderOption {
	return func(d *decoderConfig) error {
		d.prefix = prefix
		return nil
	}
}

// Decoder reads and decodes values from a key-value storage.
type Decoder struct {
	kv     KV
	config decoderConfig
}

// NewDecoder returns a new decoder that reads from kv.
func NewDecoder(kv KV, opts ...DecoderOption) (*Decoder, error) {
	dec := &Decoder{kv: kv}
	cfg := defaultConfig

	var errs []error
	for _, opt := range opts {
		if err := opt(&cfg); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	dec.config = cfg
	return dec, nil
}

// Decode reads values from the key-value storage and decodes them into v.
func (d *Decoder) Decode(v any) error {
	return d.DecodeContext(context.Background(), v)
}

// DecodeContext reads values from the key-value storage and decodes them into v.
// It uses the provided context for the deadline and cancellation.
func (d *Decoder) DecodeContext(ctx context.Context, v any) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return fmt.Errorf("decode target must be a non-nil pointer")
	}
	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("decode target must be a pointer to a struct")
	}

	return d.decodeStruct(ctx, d.config.prefix, val)
}

func (d *Decoder) decodeStruct(ctx context.Context, prefix string, val reflect.Value) error {
	t := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if !field.CanSet() || !field.CanAddr() {
			continue // Skip unexported and unaddressable fields
		}

		fieldType := t.Field(i)
		if err := d.decodeField(ctx, field, fieldType, prefix); err != nil {
			return fmt.Errorf("decode field %s: %w", fieldType.Name, err)
		}
	}

	return nil
}

func (d *Decoder) decodeField(ctx context.Context, f reflect.Value, t reflect.StructField, prefix string) error {
	tag := t.Tag.Get(d.config.tag)
	if tag == "" {
		return nil // Skip fields without consul tag
	}

	tagSpec := getTagSpec(tag)

	key := prefix + tagSpec.key

	initField(f, t)

	// check if the field is a struct or a pointer to a struct.
	// if so, recursively decode the struct.
	if sf, ok := structField(f, t); ok {
		if err := d.decodeStruct(ctx, key+"/", sf); err != nil {
			return fmt.Errorf("decode struct field %s: %w", t.Name, err)
		}
		return nil
	}

	value, err := d.kv.Get(ctx, key)
	if err != nil {
		return fmt.Errorf("get key %q: %w", key, err)
	}

	opts := ValueUnmarshalOpts{SliceSep: d.config.sliceSep}
	var out any
	if f.Kind() == reflect.Ptr {
		out = f.Interface()
	} else {
		out = f.Addr().Interface()
	}

	if err := value.UnmarshalTo(out, opts); err != nil {
		return fmt.Errorf("unmarshal value of %q: %w", key, err)
	}

	return nil
}

// UnmarshalContext reads values from the key-value storage and decodes them into v
// using the provided context and default decoder configuration.
func UnmarshalContext(ctx context.Context, kv KV, v any) error {
	dec, err := NewDecoder(kv)
	if err != nil {
		return err
	}
	return dec.DecodeContext(ctx, v)
}

// Unmarshal reads values from the key-value storage and decodes them into v.
// See [UnmarshalContext] for more details.
func Unmarshal(kv KV, v any) error {
	return UnmarshalContext(context.Background(), kv, v)
}
