package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"ponglehub.co.uk/geppetto/config"
)

func TestCircularDependencies(t *testing.T) {
	t.Run("No repos", func(t *testing.T) {
		b := Builder{
			plan: planner{},
			cfg:  &config.Config{},
		}

		collisions, ok := b.hasCircularDependencies()
		assert.Equal(t, true, ok)
		assert.Equal(t, []string{}, collisions)
	})

	t.Run("No collisions", func(t *testing.T) {
		b := Builder{
			plan: planner{},
			cfg: &config.Config{
				Repos: []config.Repo{
					config.Repo{
						Name:      "repo1",
						DependsOn: []string{"repo3"},
					},
					config.Repo{
						Name:      "repo2",
						DependsOn: []string{},
					},
					config.Repo{
						Name:      "repo3",
						DependsOn: []string{"repo2"},
					},
				},
			},
		}

		collisions, ok := b.hasCircularDependencies()
		assert.Equal(t, true, ok)
		assert.Equal(t, []string{}, collisions)
	})

	t.Run("Two way collision", func(t *testing.T) {
		b := Builder{
			plan: planner{},
			cfg: &config.Config{
				Repos: []config.Repo{
					config.Repo{
						Name:      "repo1",
						DependsOn: []string{"repo3"},
					},
					config.Repo{
						Name:      "repo2",
						DependsOn: []string{},
					},
					config.Repo{
						Name:      "repo3",
						DependsOn: []string{"repo1"},
					},
				},
			},
		}

		collisions, ok := b.hasCircularDependencies()
		assert.Equal(t, false, ok)
		assert.Equal(t, []string{"repo1", "repo3"}, collisions)
	})

	t.Run("Three way collision", func(t *testing.T) {
		b := Builder{
			plan: planner{},
			cfg: &config.Config{
				Repos: []config.Repo{
					config.Repo{
						Name:      "repo1",
						DependsOn: []string{"repo3"},
					},
					config.Repo{
						Name:      "repo2",
						DependsOn: []string{"repo1"},
					},
					config.Repo{
						Name:      "repo3",
						DependsOn: []string{"repo2"},
					},
				},
			},
		}

		collisions, ok := b.hasCircularDependencies()
		assert.Equal(t, false, ok)
		assert.Equal(t, []string{"repo1", "repo2", "repo3"}, collisions)
	})
}
