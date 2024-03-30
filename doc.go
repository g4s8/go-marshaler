/*
Package marshaler provides a simple interface for marshaling and unmarshaling
values from a key-value storage.

The package provides a Decoder that reads and decodes values from a key-value
storage. The Decoder uses a set of options to configure the behavior of the
decoder.

The package also provides a set of interfaces and types to work with values
from a key-value storage. The Value interface represents a value from a
key-value storage that can be unmarshaled to a target type. The ValueUnmarshalOpts
type is a set of options for value unmarshaling.

Example:

	package main

	import (
		"log"
		"time"

		"github.com/g4s8/go-marshaler"
	)

	type LoggerConfig struct {
		Level  string `kv:"level"`
		Output string `kv:"output"`
	}

	type Config struct {
		Host    *string       `kv:"host"`
		Port    int           `kv:"port"`
		Debug   *bool         `kv:"debug"`
		Opt     *bool         `kv:"opt,omitempty"`
		Timeout time.Duration `kv:"timeout"`
		Logger  *LoggerConfig `kv:"logger"`
		Params  []string      `kv:"params"`
	}

	func main() {
		kv := marshaler.MapKV{
			"host":          "localhost",
			"port":          "8080",
			"logger/level":  "info",
			"logger/output": "stdout",
			"timeout":       "5s",
			"params":        "a,b,c",
		}
		decoder, err := marshaler.NewDecoder(kv, marshaler.WithSeparator("/"), marshaler.WithTag("kv"))
		if err != nil {
			log.Fatalf("error creating decoder: %v", err)
		}
		var cfg Config
		if err := decoder.Decode(&cfg); err != nil {
			log.Fatalf("error decoding: %v", err)
		}
		log.Printf("decoded config: %+v", cfg)
	}

There is also a [Scanner] interface that can be implemented by a target type to
provide custom unmarshaling logic.
*/
package marshaler
