# go-marshaler/consul

The `consul` submodule of `go-marshaler` extends the base library
to support decoding configuration data directly from HashiCorp Consul KV stores into Go structs,
including complex and nested structures.

## Getting Started

### Installation

Install the `consul` submodule:

```sh
go get github.com/g4s8/go-marshaler/consul
```

### Usage

Decode a configuration stored in Consul KV into a Go struct, showcasing nested structures:

```go
package main

import (
	"fmt"

	"github.com/g4s8/go-marshaler/consul"
	capi "github.com/hashicorp/consul/api"
)

type Config struct {
	Host   string `consul:"host"`
	Port   int    `consul:"port"`
	Path   string `consul:"path"`
	Logger struct {
		Level   string `consul:"level"`
		Output  string `consul:"output"`
		Enabled bool   `consul:"enabled"`
	} `consul:"logger"`
}

func (cfg Config) String() string {
	return fmt.Sprintf("Config{Host: %s, Port: %d, Path: %s, Logger: {Level: %s, Output: %s, Enabled: %t}}",
		cfg.Host, cfg.Port, cfg.Path, cfg.Logger.Level, cfg.Logger.Output, cfg.Logger.Enabled)
}

func main() {
	const consulAddr = "http://localhost:8500"
	cli, err := capi.NewClient(&capi.Config{Address: consulAddr})
	if err != nil {
		panic(err)
	}
	var cfg Config
	if err := consul.Unmarshal(cli, &cfg); err != nil {
		panic(err)
	}
	fmt.Printf("config: %s\n", cfg)
}
```

## Advanced Usage

The `consul` submodule provides a tailored API for integrating with Consul KV stores,
using the `consul` field tag to map configuration values. It offers flexibility and customizability
through several configuration options and a distinct public API, allowing for more precise control over the unmarshalling process.

### Configuration Options

To customize the behavior of the Consul decoder, the following options are available:

 - `WithSliceSeparator(separator string)`: Sets the separator for slice values. The default separator is `/`.
 - `WithPrefix(prefix string)`: Sets a global prefix for all keys during the decoding process.
 This is useful for scoping keys within a certain namespace in Consul KV.

These options can be passed to the decoder constructor to modify its behavior.

### Public API

The `consul` submodule exposes an API for decoding data from Consul:

 - `NewDecoder(cli *capi.Client, opts ...DecoderOption) *Decoder`: Creates a new Consul decoder.
 This decoder provides a `Decode(any) error` method for unmarshalling data into a provided struct.
 - `Unmarshal(cli *capi.Client, v any) error`: A convenience function to unmarshal data from the Consul client
 into the provided variable `v`.
 - `UnmarshalContext(ctx context.Context, cli *capi.Client, v any) error`: Similar to `Unmarshal` but allows passing a
 `context.Context` for deadline or cancellation.
 - `UnmarshalDefault(v any) error`: Unmarshals data using a default Consul client,
 which is configured through environment variables.
 - `UnmarshalDefaultContext(ctx context.Context, v any) error`: The same as `UnmarshalDefault`
 but with support for `context.Context`.

