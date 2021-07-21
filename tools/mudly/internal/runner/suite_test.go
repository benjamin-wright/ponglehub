package runner_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"ponglehub.co.uk/tools/mudly/internal/runner"
	"ponglehub.co.uk/tools/mudly/internal/solver"
)

func TestRun(t *testing.T) {
	for _, test := range []struct {
		Name  string
		Nodes []*solver.Node
	}{
		{Name: "test"},
	} {
		t.Run(test.Name, func(u *testing.T) {
			err := runner.Run(test.Nodes)

			assert.NoError(u, err)
		})
	}
}
