package consul

import (
	"os"
	"strconv"
	"testing"
	"time"

	capi "github.com/hashicorp/consul/api"
)

type testTarget struct {
	Host   string `consul:"host"`
	Port   int    `consul:"port"`
	Debug  bool   `consul:"debug"`
	Logger struct {
		Level  string `consul:"level"`
		Output string `consul:"output"`
	} `consul:"logger"`
	Timeout time.Duration `consul:"timeout"`
}

func TestDecoder(t *testing.T) {
	consulAddr := os.Getenv("TEST_CONSUL_ADDR")
	if consulAddr == "" {
		t.Skip("TEST_CONSUL_ADDR is not set")
	}

	cli, err := capi.NewClient(&capi.Config{Address: consulAddr})
	if err != nil {
		t.Fatalf("error creating consul client: %v", err)
	}

	const (
		host      = "localhost"
		port      = 8080
		debug     = true
		logLevel  = "info"
		logOutput = "stdout"
		timeout   = 5 * time.Second
	)

	kv := cli.KV()
	pairs := []*capi.KVPair{
		{Key: "host", Value: []byte(host)},
		{Key: "port", Value: []byte(strconv.Itoa(port))},
		{Key: "debug", Value: []byte(strconv.FormatBool(debug))},
		{Key: "logger/level", Value: []byte(logLevel)},
		{Key: "logger/output", Value: []byte(logOutput)},
		{Key: "timeout", Value: []byte(timeout.String())},
	}
	for _, pair := range pairs {
		_, err := kv.Put(pair, nil)
		if err != nil {
			t.Fatalf("error putting key %q: %v", pair.Key, err)
		}
		t.Logf("put key %q", pair.Key)
		t.Cleanup(func() {
			_, err := kv.Delete(pair.Key, nil)
			if err != nil {
				t.Logf("error deleting key %q: %v", pair.Key, err)
			}
			t.Logf("deleted key %q", pair.Key)
		})
	}

	var target testTarget
	t.Logf("unmarshaling config from consul")
	if err := Unmarshal(cli, &target); err != nil {
		t.Fatalf("error unmarshaling: %v", err)
	}

	t.Logf("checking config")
	if target.Host != host {
		t.Errorf("unexpected host: %q, expected %q", target.Host, host)
	}
	if target.Port != port {
		t.Errorf("unexpected port: %d, expected %d", target.Port, port)
	}
	if target.Debug != debug {
		t.Errorf("unexpected debug: %t, expected %t", target.Debug, debug)
	}
	if target.Logger.Level != logLevel {
		t.Errorf("unexpected log level: %q, expected %q", target.Logger.Level, logLevel)
	}
	if target.Logger.Output != logOutput {
		t.Errorf("unexpected log output: %q, expected %q", target.Logger.Output, logOutput)
	}
	if target.Timeout != timeout {
		t.Errorf("unexpected timeout: %v, expected %v", target.Timeout, timeout)
	}
}
