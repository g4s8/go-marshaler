# go-marshaler

`go-marshaler` is a Go library designed to simplify the unmarshalling and decoding
of arbitrary key-value data into Go structs,
with support for nested structures.
This makes it particularly useful for complex configurations.
Leveraging struct tags for custom field mappings,
it provides a flexible way to work with configuration data stored in various formats.


## Implementations

Here is the list of implementations:
 - [consul](./consul)

## Getting Started

### Installation

Install `go-marshaler` by running:

```sh
go get github.com/g4s8/go-marshaler
```

### Usage

A simple example to decode a map of key-value pairs into a Go struct, demonstrating nested structure support:

```go
package main

import (
    "log"
    "github.com/g4s8/go-marshaler"
)

// Define your configuration struct
type Config struct {
    Host    *string       `kv:"host"`
    Port    int           `kv:"port"`
    Timeout time.Duration `kv:"timeout"`
    Logger  *LoggerConfig `kv:"logger"`
    // Add other fields...
}

type LoggerConfig struct {
    Level  string `kv:"level"`
    Output string `kv:"output"`
}

func main() {
    kv := marshaler.MapKV{
        "host":          "localhost",
        "port":          "8080",
        "timeout":       "5s",
        "logger/level":  "info",
        "logger/output": "stdout",
        // Populate with your key-value pairs
    }
    decoder := marshaler.NewDecoder(kv, marshaler.WithSeparator("/"), marshaler.WithTag("kv"))
    var cfg Config
    if err := decoder.Decode(&cfg); err != nil {
        log.Fatalf("error decoding: %v", err)
    }
    log.Printf("decoded config: %+v", &cfg)
}
```

## Configuration Options

`go-marshaler` allows customizing the decoding process through various options:

- `WithSeparator(string)`: Specifies the separator for nested keys.
- `WithSliceSeparator(string)`: Specifies the separator for slice values.
- `WithTag(string)`: Sets the struct tag to use for key mapping.
- `WithPrefix(string)`: Adds a prefix to all keys during decoding.

Refer to the API documentation for more details on how to use these options.
