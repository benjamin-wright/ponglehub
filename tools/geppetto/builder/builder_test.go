package builder

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"ponglehub.co.uk/geppetto/types"
// )

// func calls(channel chan call) []string {
// 	result := []string{}
// 	close(channel)

// 	for c := range channel {
// 		result = append(result, fmt.Sprintf("%s:%s", c.repo, c.lang))
// 	}

// 	return result
// }

// func TestBuilderBuild(t *testing.T) {
// 	t.Run("npm modules", func(t *testing.T) {
// 		progress := make(chan buildState, 5)
// 		channel, worker := makeMockWorker()
// 		b := Builder{
// 			worker: worker,
// 		}

// 		state, err := b.Build([]types.Repo{
// 			{Name: "repo1", RepoType: types.Node},
// 			{Name: "repo2", RepoType: types.Node},
// 			{Name: "repo3", RepoType: types.Node},
// 		}, progress)

// 		assert.Nil(t, err)
// 		assert.Equal(t, state.Repos, []RepoState{
// 			{repo: "repo1", state: BuiltState},
// 			{repo: "repo2", state: BuiltState},
// 			{repo: "repo3", state: BuiltState},
// 		})
// 		assert.ElementsMatch(t, calls(channel), []string{
// 			"repo1:npm",
// 			"repo2:npm",
// 			"repo3:npm",
// 		})
// 	})

// 	t.Run("with dependencies", func(t *testing.T) {
// 		progress := make(chan buildState, 5)
// 		_, worker := makeMockWorker()
// 		b := Builder{
// 			worker: worker,
// 		}

// 		state, err := b.Build([]types.Repo{
// 			{Name: "repo1", RepoType: types.Node, DependsOn: []string{"repo3"}},
// 			{Name: "repo2", RepoType: types.Node, DependsOn: []string{"repo1", "repo3"}},
// 			{Name: "repo3", RepoType: types.Node},
// 		}, progress)

// 		assert.Nil(t, err)
// 		assert.Equal(t, state.Repos, []RepoState{
// 			{repo: "repo3", state: BuiltState},
// 			{repo: "repo1", state: BuiltState},
// 			{repo: "repo2", state: BuiltState},
// 		})
// 	})

// 	t.Run("unknown repo type", func(t *testing.T) {
// 		progress := make(chan buildState, 5)
// 		_, worker := makeMockWorker()
// 		b := Builder{
// 			worker: worker,
// 		}

// 		state, err := b.Build([]types.Repo{
// 			{Name: "repo1", RepoType: types.Node},
// 			{Name: "repo2"},
// 		}, progress)

// 		assert.Error(t, err)
// 		assert.Equal(t, state.Repos, []RepoState{
// 			{repo: "repo1", state: BuildingState},
// 			{repo: "repo2", state: ErroredState},
// 		})
// 	})
// }
