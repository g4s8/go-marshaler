package marshaler

import (
	"errors"
	"testing"
	"time"
)

type scanner struct {
	value string
}

type errScanner struct{}

var (
	_ Scanner = (*scanner)(nil)
	_ Scanner = (*errScanner)(nil)
)

func (b *scanner) Scan(v any) error {
	if v, ok := v.(string); ok {
		b.value = v
		return nil
	}

	return errors.New("unexpected type")
}

func (b *errScanner) Scan(v any) error {
	return errors.New("scan error")
}

func TestNilValue(t *testing.T) {
	var s string
	if err := NullValue.UnmarshalTo(&s, ValueUnmarshalOpts{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s != "" {
		t.Fatalf("expected empty string, got %q", s)
	}
}

func TestStringValue(t *testing.T) {
	t.Run("Scanner", func(t *testing.T) {
		var target scanner
		const value = "hello"
		v := NewStringValue(value)
		if err := v.UnmarshalTo(&target, ValueUnmarshalOpts{}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if target.value != value {
			t.Fatalf("expected %q, got %q", value, target.value)
		}
	})
	t.Run("Duration", func(t *testing.T) {
		var target time.Duration
		const value = "1s"
		v := NewStringValue(value)
		if err := v.UnmarshalTo(&target, ValueUnmarshalOpts{}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if target != time.Second {
			t.Fatalf("expected %v, got %v", time.Second, target)
		}
	})
	t.Run("Time", func(t *testing.T) {
		var target time.Time
		const value = "2021-01-01T00:00:00Z"
		v := NewStringValue(value)
		if err := v.UnmarshalTo(&target, ValueUnmarshalOpts{}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if target.Format(time.RFC3339) != value {
			t.Fatalf("expected %q, got %q", value, target.Format(time.RFC3339))
		}
	})
	t.Run("String", func(t *testing.T) {
		var target string
		const value = "hello"
		v := NewStringValue(value)
		if err := v.UnmarshalTo(&target, ValueUnmarshalOpts{}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if target != value {
			t.Fatalf("expected %q, got %q", value, target)
		}
	})
	t.Run("Int", func(t *testing.T) {
		var target int
		const value = "42"
		v := NewStringValue(value)
		if err := v.UnmarshalTo(&target, ValueUnmarshalOpts{}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if target != 42 {
			t.Fatalf("expected %d, got %d", 42, target)
		}
	})
	t.Run("Int8", func(t *testing.T) {
		var target int8
		const value = "42"
		v := NewStringValue(value)
		if err := v.UnmarshalTo(&target, ValueUnmarshalOpts{}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if target != 42 {
			t.Fatalf("expected %d, got %d", 42, target)
		}
	})
	t.Run("Int16", func(t *testing.T) {
		var target int16
		const value = "42"
		v := NewStringValue(value)
		if err := v.UnmarshalTo(&target, ValueUnmarshalOpts{}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if target != 42 {
			t.Fatalf("expected %d, got %d", 42, target)
		}
	})
	t.Run("UInt", func(t *testing.T) {
		var target uint
		const value = "42"
		v := NewStringValue(value)
		if err := v.UnmarshalTo(&target, ValueUnmarshalOpts{}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if target != 42 {
			t.Fatalf("expected %d, got %d", 42, target)
		}
	})
	t.Run("UInt8", func(t *testing.T) {
		var target uint8
		const value = "42"
		v := NewStringValue(value)
		if err := v.UnmarshalTo(&target, ValueUnmarshalOpts{}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if target != 42 {
			t.Fatalf("expected %d, got %d", 42, target)
		}
	})
	t.Run("UInt16", func(t *testing.T) {
		var target uint16
		const value = "42"
		v := NewStringValue(value)
		if err := v.UnmarshalTo(&target, ValueUnmarshalOpts{}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if target != 42 {
			t.Fatalf("expected %d, got %d", 42, target)
		}
	})
	t.Run("UInt32", func(t *testing.T) {
		var target uint32
		const value = "42"
		v := NewStringValue(value)
		if err := v.UnmarshalTo(&target, ValueUnmarshalOpts{}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if target != 42 {
			t.Fatalf("expected %d, got %d", 42, target)
		}
	})
	t.Run("UInt64", func(t *testing.T) {
		var target uint64
		const value = "42"
		v := NewStringValue(value)
		if err := v.UnmarshalTo(&target, ValueUnmarshalOpts{}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if target != 42 {
			t.Fatalf("expected %d, got %d", 42, target)
		}
	})
	t.Run("Float32", func(t *testing.T) {
		var target float32
		const value = "42.42"
		v := NewStringValue(value)
		if err := v.UnmarshalTo(&target, ValueUnmarshalOpts{}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if target != 42.42 {
			t.Fatalf("expected %f, got %f", 42.42, target)
		}
	})
	t.Run("Float64", func(t *testing.T) {
		var target float64
		const value = "42.42"
		v := NewStringValue(value)
		if err := v.UnmarshalTo(&target, ValueUnmarshalOpts{}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if target != 42.42 {
			t.Fatalf("expected %f, got %f", 42.42, target)
		}
	})
	t.Run("Bool", func(t *testing.T) {
		var target bool
		const value = "true"
		v := NewStringValue(value)
		if err := v.UnmarshalTo(&target, ValueUnmarshalOpts{}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !target {
			t.Fatalf("expected true, got false")
		}
	})
	t.Run("Slice", func(t *testing.T) {
		var target []string
		const value = "hello,world"
		v := NewStringValue(value)
		if err := v.UnmarshalTo(&target, ValueUnmarshalOpts{SliceSep: ","}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(target) != 2 {
			t.Fatalf("expected 2 elements, got %d", len(target))
		}
		if target[0] != "hello" {
			t.Fatalf("expected %q, got %q", "hello", target[0])
		}
		if target[1] != "world" {
			t.Fatalf("expected %q, got %q", "world", target[1])
		}
		t.Run("NoSep", func(t *testing.T) {
			var target []string
			const value = "hello"
			v := NewStringValue(value)
			if err := v.UnmarshalTo(&target, ValueUnmarshalOpts{}); err == nil {
				t.Fatalf("expected error, got nil")
			}
		})
	})
	t.Run("Invalid", func(t *testing.T) {
		t.Run("Scanner", func(t *testing.T) {
			var target errScanner
			const value = "hello"
			v := NewStringValue(value)
			if err := v.UnmarshalTo(&target, ValueUnmarshalOpts{}); err == nil {
				t.Fatalf("expected error, got nil")
			}
		})
		t.Run("Duration", func(t *testing.T) {
			var target time.Duration
			const value = "invalid"
			v := NewStringValue(value)
			if err := v.UnmarshalTo(&target, ValueUnmarshalOpts{}); err == nil {
				t.Fatalf("expected error, got nil")
			}
		})
		t.Run("Time", func(t *testing.T) {
			var target time.Time
			const value = "invalid"
			v := NewStringValue(value)
			if err := v.UnmarshalTo(&target, ValueUnmarshalOpts{}); err == nil {
				t.Fatalf("expected error, got nil")
			}
		})
		t.Run("Int", func(t *testing.T) {
			var target int
			const value = "invalid"
			v := NewStringValue(value)
			if err := v.UnmarshalTo(&target, ValueUnmarshalOpts{}); err == nil {
				t.Fatalf("expected error, got nil")
			}
		})
		t.Run("Int8", func(t *testing.T) {
			var target int8
			const value = "invalid"
			v := NewStringValue(value)
			if err := v.UnmarshalTo(&target, ValueUnmarshalOpts{}); err == nil {
				t.Fatalf("expected error, got nil")
			}
		})
		t.Run("Int16", func(t *testing.T) {
			var target int16
			const value = "invalid"
			v := NewStringValue(value)
			if err := v.UnmarshalTo(&target, ValueUnmarshalOpts{}); err == nil {
				t.Fatalf("expected error, got nil")
			}
		})
		t.Run("Int32", func(t *testing.T) {
			var target int32
			const value = "invalid"
			v := NewStringValue(value)
			if err := v.UnmarshalTo(&target, ValueUnmarshalOpts{}); err == nil {
				t.Fatalf("expected error, got nil")
			}
		})
		t.Run("Int64", func(t *testing.T) {
			var target int64
			const value = "invalid"
			v := NewStringValue(value)
			if err := v.UnmarshalTo(&target, ValueUnmarshalOpts{}); err == nil {
				t.Fatalf("expected error, got nil")
			}
		})
		t.Run("UInt", func(t *testing.T) {
			var target uint
			const value = "invalid"
			v := NewStringValue(value)
			if err := v.UnmarshalTo(&target, ValueUnmarshalOpts{}); err == nil {
				t.Fatalf("expected error, got nil")
			}
		})
		t.Run("UInt8", func(t *testing.T) {
			var target uint8
			const value = "invalid"
			v := NewStringValue(value)
			if err := v.UnmarshalTo(&target, ValueUnmarshalOpts{}); err == nil {
				t.Fatalf("expected error, got nil")
			}
		})
		t.Run("UInt16", func(t *testing.T) {
			var target uint16
			const value = "invalid"
			v := NewStringValue(value)
			if err := v.UnmarshalTo(&target, ValueUnmarshalOpts{}); err == nil {
				t.Fatalf("expected error, got nil")
			}
		})
		t.Run("UInt32", func(t *testing.T) {
			var target uint32
			const value = "invalid"
			v := NewStringValue(value)
			if err := v.UnmarshalTo(&target, ValueUnmarshalOpts{}); err == nil {
				t.Fatalf("expected error, got nil")
			}
		})
		t.Run("UInt64", func(t *testing.T) {
			var target uint64
			const value = "invalid"
			v := NewStringValue(value)
			if err := v.UnmarshalTo(&target, ValueUnmarshalOpts{}); err == nil {
				t.Fatalf("expected error, got nil")
			}
		})
		t.Run("Float32", func(t *testing.T) {
			var target float32
			const value = "invalid"
			v := NewStringValue(value)
			if err := v.UnmarshalTo(&target, ValueUnmarshalOpts{}); err == nil {
				t.Fatalf("expected error, got nil")
			}
		})
		t.Run("Float64", func(t *testing.T) {
			var target float64
			const value = "invalid"
			v := NewStringValue(value)
			if err := v.UnmarshalTo(&target, ValueUnmarshalOpts{}); err == nil {
				t.Fatalf("expected error, got nil")
			}
		})
		t.Run("Bool", func(t *testing.T) {
			var target bool
			const value = "invalid"
			v := NewStringValue(value)
			if err := v.UnmarshalTo(&target, ValueUnmarshalOpts{}); err == nil {
				t.Fatalf("expected error, got nil")
			}
		})
	})
	t.Run("Unsupported", func(t *testing.T) {
		var target interface{}
		const value = "hello"
		v := NewStringValue(value)
		if err := v.UnmarshalTo(&target, ValueUnmarshalOpts{}); err == nil {
			t.Fatalf("expected error, got nil")
		}
	})
}
