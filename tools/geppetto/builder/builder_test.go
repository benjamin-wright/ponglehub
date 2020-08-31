package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"ponglehub.co.uk/geppetto/types"
)

func TestBuilderBuild(t *testing.T) {
	t.Run("npm modules", func(t *testing.T) {
		worker := &mockWorker{}
		b := Builder{
			worker: worker,
		}

		err := b.Build([]types.Repo{
			{Name: "repo1", RepoType: types.Node},
			{Name: "repo2", RepoType: types.Node},
		})

		assert.Nil(t, err)
	})

	t.Run("unknown repo type", func(t *testing.T) {
		worker := &mockWorker{}
		b := Builder{
			worker: worker,
		}

		err := b.Build([]types.Repo{
			{Name: "repo1", RepoType: types.Node},
			{Name: "repo2"},
		})

		assert.Error(t, err)
	})
}
