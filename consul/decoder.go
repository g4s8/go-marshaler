package consul

import (
	"context"
	"fmt"

	"github.com/g4s8/go-marshaler"
	capi "github.com/hashicorp/consul/api"
)

type decoderConfig struct {
	sliceSep string
	prefix   string
}

type DecoderOption func(*decoderConfig)

func WithSliceSeparator(separator string) DecoderOption {
	return func(d *decoderConfig) {
		d.sliceSep = separator
	}
}

func WithPrefix(prefix string) DecoderOption {
	return func(d *decoderConfig) {
		d.prefix = prefix
	}
}

type Decoder struct {
	dec *marshaler.Decoder
}

func NewDecoder(cli *capi.Client, opts ...DecoderOption) (*Decoder, error) {
	kv := &consulKV{ckv: cli.KV()}
	var cfg decoderConfig
	for _, opt := range opts {
		opt(&cfg)
	}
	decOpts := make([]marshaler.DecoderOption, 0)
	decOpts = append(decOpts, marshaler.WithSeparator("/"))
	decOpts = append(decOpts, marshaler.WithTag("consul"))
	if cfg.sliceSep != "" {
		decOpts = append(decOpts, marshaler.WithSliceSeparator(cfg.sliceSep))
	}
	if cfg.prefix != "" {
		decOpts = append(decOpts, marshaler.WithPrefix(cfg.prefix))
	}
	dec, err := marshaler.NewDecoder(kv, decOpts...)
	if err != nil {
		return nil, fmt.Errorf("create decoder: %w", err)
	}
	return &Decoder{dec: dec}, nil
}

func (d *Decoder) Decode(v any) error {
	return d.dec.Decode(v)
}

func (d *Decoder) DecodeContext(ctx context.Context, v any) error {
	return d.dec.DecodeContext(ctx, v)
}

func Unmarshal(cli *capi.Client, v any) error {
	return UnmarshalContext(context.Background(), cli, v)
}

func UnmarshalContext(ctx context.Context, cli *capi.Client, v any) error {
	dec, err := NewDecoder(cli)
	if err != nil {
		return fmt.Errorf("create decoder: %w", err)
	}
	return dec.DecodeContext(ctx, v)
}

func UnmarshalDefault(v any) error {
	return UnmarshalDefaultContext(context.Background(), v)
}

func UnmarshalDefaultContext(ctx context.Context, v any) error {
	cli, err := capi.NewClient(capi.DefaultConfig())
	if err != nil {
		return fmt.Errorf("create consul client: %w", err)
	}
	return UnmarshalContext(ctx, cli, v)
}
