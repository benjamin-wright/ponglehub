package solver_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"ponglehub.co.uk/tools/mudly/internal/config"
	"ponglehub.co.uk/tools/mudly/internal/solver"
	"ponglehub.co.uk/tools/mudly/internal/target"
)

func TestSolver(t *testing.T) {
	for _, test := range []struct {
		Name     string
		Targets  []target.Target
		Configs  []config.Config
		Expected []solver.Node
	}{
		{
			Name:    "simplest",
			Targets: []target.Target{{Dir: ".", Artefact: "image"}},
			Configs: []config.Config{
				{
					Path: ".",
					Artefacts: []config.Artefact{
						{
							Name: "image",
							Pipeline: config.Pipeline{
								Steps: []interface{}{
									config.CommandStep{
										Name:    "build",
										Command: "go build -o ./bin/mudly ./cmd/mudly",
									},
								},
							},
						},
					},
				},
			},
			Expected: []solver.Node{
				{
					Path:      ".",
					Artefact:  "image",
					Step:      "build",
					Command:   "go build -o ./bin/mudly ./cmd/mudly",
					State:     solver.STATE_NONE,
					DependsOn: []*solver.Node{},
				},
			},
		},
	} {
		t.Run(test.Name, func(u *testing.T) {
			nodes, err := solver.Solve(test.Targets, test.Configs)

			assert.NoError(u, err, "didn't expect an error")

			if test.Expected != nil {
				if nodes != nil {
					assert.Equal(u, test.Expected, nodes)
				} else {
					assert.Fail(u, "expected a list of nodes")
				}
			}
		})
	}
}
