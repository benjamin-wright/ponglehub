package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"ponglehub.co.uk/geppetto/types"
)

func TestBuilderBuild(t *testing.T) {
	t.Run("npm modules", func(t *testing.T) {
		_, worker := makeMockWorker()
		b := Builder{
			worker: worker,
		}

		state, err := b.Build([]types.Repo{
			{Name: "repo1", RepoType: types.Node},
			{Name: "repo2", RepoType: types.Node},
			{Name: "repo3", RepoType: types.Node},
		})

		assert.Nil(t, err)
		assert.Equal(t, state.repos, []repoState{
			repoState{repo: "repo1", state: BuiltState},
			repoState{repo: "repo2", state: BuiltState},
			repoState{repo: "repo3", state: BuiltState},
		})
	})

	t.Run("with dependencies", func(t *testing.T) {
		_, worker := makeMockWorker()
		b := Builder{
			worker: worker,
		}

		state, err := b.Build([]types.Repo{
			{Name: "repo1", RepoType: types.Node, DependsOn: []string{"repo3"}},
			{Name: "repo2", RepoType: types.Node, DependsOn: []string{"repo1", "repo3"}},
			{Name: "repo3", RepoType: types.Node},
		})

		assert.Nil(t, err)
		assert.Equal(t, state.repos, []repoState{
			repoState{repo: "repo3", state: BuiltState},
			repoState{repo: "repo1", state: BuiltState},
			repoState{repo: "repo2", state: BuiltState},
		})
	})

	t.Run("unknown repo type", func(t *testing.T) {
		_, worker := makeMockWorker()
		b := Builder{
			worker: worker,
		}

		state, err := b.Build([]types.Repo{
			{Name: "repo1", RepoType: types.Node},
			{Name: "repo2"},
		})

		assert.Error(t, err)
		assert.Equal(t, state.repos, []repoState{
			repoState{repo: "repo1", state: BuildingState},
			repoState{repo: "repo2", state: ErroredState},
		})
	})
}
