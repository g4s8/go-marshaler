package marshaler

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// NullValue is a special value that represents a null value.
//
// It doesn't perfor any unmarshaling operatrion.
var NullValue Value = nullValue{}

type nullValue struct{}

func (nullValue) UnmarshalTo(out any, opts ValueUnmarshalOpts) error {
	return nil
}

// StringValue holds a string value, that can be unmarshaled to a target type.
//
// It supports unmarshaling to the following types:
// - string
// - int, int8, int16, int32, int64
// - uint, uint8, uint16, uint32, uint64
// - bool
// - float32, float64
// - time.Duration
// - time.Time (RFC3339)
// - []string
// - Scanner
type StringValue struct {
	value string
}

// NewStringValue creates a new string value from a string.
func NewStringValue(value string) StringValue {
	return StringValue{value: value}
}

// NewBytesValue creates a new string value from a byte slice.
func NewBytesValue(value []byte) StringValue {
	return StringValue{value: string(value)}
}

// UnmarshalTo unmarshals the string value to a target type using the provided options.
func (v StringValue) UnmarshalTo(out any, opts ValueUnmarshalOpts) error {
	if s, ok := out.(Scanner); ok {
		if err := s.Scan(v.value); err != nil {
			return fmt.Errorf("unmarshal binary: %w", err)
		}
		return nil
	}

	var parseErr error
	switch out := out.(type) {
	case *time.Duration:
		duration, err := time.ParseDuration(v.value)
		if err != nil {
			return fmt.Errorf("parse duration from %q: %w", v.value, err)
		}
		*out = duration
	case *time.Time:
		t, err := time.Parse(time.RFC3339, v.value)
		if err != nil {
			return fmt.Errorf("parse time from %q: %w", v.value, err)
		}
		*out = t

	case *string:
		*out = v.value

	case *int:
		intVal, err := strconv.Atoi(v.value)
		if err != nil {
			return fmt.Errorf("parse int from %q: %w", v.value, err)
		}
		*out = intVal
	case *int8:
		parseErr = unmarshalIntNumber(v.value, 8, out)
	case *int16:
		parseErr = unmarshalIntNumber(v.value, 16, out)
	case *int32:
		parseErr = unmarshalIntNumber(v.value, 32, out)
	case *int64:
		parseErr = unmarshalIntNumber(v.value, 64, out)
	case *uint:
		uintVal, err := strconv.ParseUint(v.value, 10, 64)
		if err != nil {
			return fmt.Errorf("parse uint from %q: %w", v.value, err)
		}
		*out = uint(uintVal)
	case *uint8:
		parseErr = unmarshalIntNumber(v.value, 8, out)
	case *uint16:
		parseErr = unmarshalIntNumber(v.value, 16, out)
	case *uint32:
		parseErr = unmarshalIntNumber(v.value, 32, out)
	case *uint64:
		parseErr = unmarshalIntNumber(v.value, 64, out)
	case *bool:
		boolVal, err := strconv.ParseBool(v.value)
		if err != nil {
			return fmt.Errorf("parse bool from %q: %w", v.value, err)
		}
		*out = boolVal
	case *float32:
		parseErr = unmarshalFloatNumber(v.value, 32, out)
	case *float64:
		parseErr = unmarshalFloatNumber(v.value, 64, out)

	case *[]string:
		if opts.SliceSep == "" {
			return fmt.Errorf("slice separator is not set")
		}
		sep := opts.SliceSep
		*out = strings.Split(v.value, sep)

	default:
		return fmt.Errorf("unsupported type %T", out)
	}

	if parseErr != nil {
		return fmt.Errorf("unmarshal %T: %w", out, parseErr)
	}

	return nil
}

type intNumber interface {
	int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64
}

type floatNumber interface {
	float32 | float64
}

func unmarshalIntNumber[T intNumber](s string, size int, out *T) error {
	val, err := strconv.ParseInt(s, 10, size)
	if err != nil {
		return fmt.Errorf("parse (u)int%d from %q: %w", size, s, err)
	}
	*out = T(val)
	return nil
}

func unmarshalFloatNumber[T floatNumber](s string, size int, out *T) error {
	val, err := strconv.ParseFloat(s, size)
	if err != nil {
		return fmt.Errorf("parse float%d from %q: %w", size, s, err)
	}
	*out = T(val)
	return nil
}
