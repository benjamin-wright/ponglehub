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

func TestFileStructToConfig(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		fs := FileStruct{
			Node: []RepoStruct{
				{Name: "my-repo", Path: "my-path", Dependencies: []string{"one", "two", "three"}},
			},
			Go: []RepoStruct{
				{Name: "my-other", Path: "my-second", Dependencies: []string{"a", "b", "c"}},
			},
		}

		actual, err := fs.toConfig()
		expected := Config{
			Repos: []Repo{
				{Name: "my-repo", Path: "my-path", RepoType: Node, DependsOn: []string{"one", "two", "three"}},
				{Name: "my-other", Path: "my-second", RepoType: Go, DependsOn: []string{"a", "b", "c"}},
			},
		}

		assert.Nil(t, err, "Expected error to be nil, got %+v")
		assert.Equal(t, expected, actual)
	})

	t.Run("Duplicate name", func(t *testing.T) {
		fs := FileStruct{
			Node: []RepoStruct{
				{Name: "my-repo", Path: "my-path", Dependencies: []string{"one", "two", "three"}},
			},
			Go: []RepoStruct{
				{Name: "my-other", Path: "my-second", Dependencies: []string{"a", "b", "c"}},
				{Name: "my-repo", Path: "my-third", Dependencies: []string{"one"}},
			},
		}

		_, err := fs.toConfig()

		assert.Error(t, err)
	})
}
