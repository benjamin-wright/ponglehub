package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigToString(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		cfg := &Config{
			Repos: []Repo{
				{Name: "repo1", RepoType: Node},
				{Name: "repo2", RepoType: Go},
			},
		}

		assert.Equal(t, "[repo1: Node], [repo2: Go]", cfg.String())
	})

	t.Run("Nil", func(t *testing.T) {
		var cfg *Config = nil

		assert.Equal(t, "Nil config", cfg.String())
	})

	t.Run("Empty", func(t *testing.T) {
		cfg := &Config{
			Repos: []Repo{},
		}

		assert.Equal(t, "Empty", cfg.String())
	})
}
