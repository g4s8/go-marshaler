package marshaler

import (
	"encoding"
	"reflect"
	"strings"
)

func initField(f reflect.Value, t reflect.StructField) {
	// If the field is a pointer and is nil, create a new instance
	if t.Type.Kind() == reflect.Ptr && f.IsNil() {
		f.Set(reflect.New(t.Type.Elem()))
	}
}

func tryTextUnmarshaler(f reflect.Value, val string) (bool, error) {
	// check if the field implements the TextUnmarshaler interface and use it
	iface := f.Interface()
	unmarshaler, ok := iface.(encoding.TextUnmarshaler)
	if !ok {
		return false, nil
	}
	if err := unmarshaler.UnmarshalText([]byte(val)); err != nil {
		return false, err
	}
	return true, nil
}

func structField(f reflect.Value, t reflect.StructField) (reflect.Value, bool) {
	_, isScanner := f.Interface().(Scanner)
	if isScanner {
		// scanner can scan value by itself even if it's a struct.
		return f, false
	}

	switch t.Type.Kind() {
	case reflect.Ptr:
		if e := f.Elem(); e.Kind() == reflect.Struct {
			return e, true
		}
	case reflect.Struct:
		return f, true
	}
	return f, false
}

type tagSpec struct {
	key       string
	omitempty bool
}

func getTagSpec(tag string) tagSpec {
	// tag could be
	//  `kv:"myKey,omitempty"`
	//  `kv:"myKey"`

	specs := strings.Split(tag, ",")
	if len(specs) == 0 {
		panic("invalid tag, should be checked before calling getTagSpec")
	}
	var spec tagSpec
	spec.key = specs[0]
	for _, s := range specs[1:] {
		switch s {
		case "omitempty":
			spec.omitempty = true
		}
	}
	return spec
}
