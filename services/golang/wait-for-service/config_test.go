package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func env(url string, interval string, timeout string) {
	os.Unsetenv("WAIT_FOR_URL")
	os.Unsetenv("WAIT_FOR_INTERVAL")
	os.Unsetenv("WAIT_FOR_TIMEOUT")

	if url != "" {
		os.Setenv("WAIT_FOR_URL", url)
	}

	if interval != "" {
		os.Setenv("WAIT_FOR_INTERVAL", interval)
	}

	if timeout != "" {
		os.Setenv("WAIT_FOR_TIMEOUT", timeout)
	}
}

func TestLoadConfig(t *testing.T) {
	for _, test := range []struct {
		name     string
		url      string
		interval string
		timeout  string
	}{
		{name: "no url"},
		{name: "interval is a string", url: "url", interval: "bob"},
		{name: "interval is a float", url: "url", interval: "1.5"},
		{name: "interval is negative", url: "url", interval: "-3"},
		{name: "interval is zero", url: "url", interval: "0"},
		{name: "timeout is a string", url: "url", timeout: "bob"},
		{name: "timeout is a float", url: "url", timeout: "1.5"},
		{name: "timeout is negative", url: "url", timeout: "-3"},
		{name: "timeout is zero", url: "url", timeout: "0"},
	} {
		t.Run(test.name, func(t *testing.T) {
			env(test.url, test.interval, test.timeout)

			_, err := newConfig()

			assert.Error(t, err)
		})
	}

	expect := func(url string, interval int, timeout int) config {
		return config{url: url, interval: interval, timeout: timeout}
	}

	for _, test := range []struct {
		name     string
		url      string
		interval string
		timeout  string
		expected config
	}{
		{name: "defaults", url: "my-url", expected: expect("my-url", 2, 60)},
		{name: "set timeout", url: "my-url", timeout: "15", expected: expect("my-url", 2, 15)},
		{name: "set interval", url: "my-url", interval: "10", expected: expect("my-url", 10, 60)},
		{name: "set both", url: "my-url", interval: "12", timeout: "36", expected: expect("my-url", 12, 36)},
	} {
		t.Run(test.name, func(t *testing.T) {
			env(test.url, test.interval, test.timeout)

			cfg, err := newConfig()

			assert.Nil(t, err)
			assert.Equal(t, test.expected, cfg)
		})
	}
}
