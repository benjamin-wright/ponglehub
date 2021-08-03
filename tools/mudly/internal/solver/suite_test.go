package solver_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"ponglehub.co.uk/tools/mudly/internal/config"
	"ponglehub.co.uk/tools/mudly/internal/runner"
	"ponglehub.co.uk/tools/mudly/internal/solver"
	"ponglehub.co.uk/tools/mudly/internal/steps"
	"ponglehub.co.uk/tools/mudly/internal/target"
)

type testNode struct {
	Path      string
	Artefact  string
	SharedEnv map[string]string
	Step      runner.Runnable
	State     runner.NodeState
	DependsOn []int
}

func convert(nodes []testNode) []*runner.Node {
	converted := []*runner.Node{}

	for id := range nodes {
		node := nodes[id]
		converted = append(converted, &runner.Node{
			Path:      node.Path,
			Artefact:  node.Artefact,
			SharedEnv: node.SharedEnv,
			Step:      node.Step,
			State:     node.State,
			DependsOn: []*runner.Node{},
		})
	}

	for idx, node := range nodes {
		for _, dep := range node.DependsOn {
			converted[idx].DependsOn = append(converted[idx].DependsOn, converted[dep])
		}
	}

	return converted
}

func TestSolver(t *testing.T) {
	for _, test := range []struct {
		Name     string
		Targets  []target.Target
		Configs  []config.Config
		Expected []testNode
	}{
		{
			Name:    "simple",
			Targets: []target.Target{{Dir: ".", Artefact: "image"}},
			Configs: []config.Config{
				{
					Path: ".",
					Artefacts: []config.Artefact{
						{
							Name: "image",
							Steps: []config.Step{
								{
									Name:    "build",
									Command: "go build -o ./bin/mudly ./cmd/mudly",
								},
								{
									Name:       "docker",
									Dockerfile: "./Dockerfile",
									Context:    ".",
								},
							},
						},
					},
				},
			},
			Expected: []testNode{
				{
					Path:     ".",
					Artefact: "image",
					Step: steps.CommandStep{
						Name:    "build",
						Command: "go build -o ./bin/mudly ./cmd/mudly",
					},
					State:     runner.STATE_PENDING,
					DependsOn: []int{},
				},
				{
					Path:     ".",
					Artefact: "image",
					Step: steps.DockerStep{
						Name:       "docker",
						Dockerfile: "./Dockerfile",
						Context:    ".",
					},
					State:     runner.STATE_PENDING,
					DependsOn: []int{0},
				},
			},
		},
		{
			Name:    "parallel builds",
			Targets: []target.Target{{Dir: ".", Artefact: "image"}, {Dir: ".", Artefact: "something"}},
			Configs: []config.Config{
				{
					Path: ".",
					Artefacts: []config.Artefact{
						{
							Name: "image",
							Steps: []config.Step{
								{
									Name:    "build",
									Command: "go build -o ./bin/mudly ./cmd/mudly",
								},
								{
									Name:       "docker",
									Dockerfile: "./Dockerfile",
									Context:    ".",
								},
							},
						},
						{
							Name: "something",
							Steps: []config.Step{
								{
									Name:    "echo",
									Command: "echo \"hi\"",
								},
								{
									Name:    "build",
									Command: "whatevs",
								},
							},
						},
					},
				},
			},
			Expected: []testNode{
				{
					Path:     ".",
					Artefact: "image",
					Step: steps.CommandStep{
						Name:    "build",
						Command: "go build -o ./bin/mudly ./cmd/mudly",
					},
					State:     runner.STATE_PENDING,
					DependsOn: []int{},
				},
				{
					Path:     ".",
					Artefact: "image",
					Step: steps.DockerStep{
						Name:       "docker",
						Dockerfile: "./Dockerfile",
						Context:    ".",
					},
					State:     runner.STATE_PENDING,
					DependsOn: []int{0},
				},
				{
					Path:     ".",
					Artefact: "something",
					Step: steps.CommandStep{
						Name:    "echo",
						Command: "echo \"hi\"",
					},
					State:     runner.STATE_PENDING,
					DependsOn: []int{},
				},
				{
					Path:     ".",
					Artefact: "something",
					Step: steps.CommandStep{
						Name:    "build",
						Command: "whatevs",
					},
					State:     runner.STATE_PENDING,
					DependsOn: []int{2},
				},
			},
		},
		{
			Name:    "linked builds",
			Targets: []target.Target{{Dir: ".", Artefact: "image"}},
			Configs: []config.Config{
				{
					Path: ".",
					Env: map[string]string{
						"GLOBAL_ENV": "value3",
					},
					Pipelines: []config.Pipeline{
						{
							Name: "my-pipeline",
							Env: map[string]string{
								"PIPELINE_ENV": "value1",
							},
							Steps: []config.Step{
								{
									Name:    "build",
									Command: "go build -o ./bin/mudly ./cmd/mudly",
									Env: map[string]string{
										"STEP_ENV": "value0",
									},
								},
								{
									Name:       "docker",
									Dockerfile: "./Dockerfile",
									Context:    ".",
								},
							},
						},
					},
					Artefacts: []config.Artefact{
						{
							Name: "image",
							Env: map[string]string{
								"ARTEFACT_ENV": "value2",
							},
							DependsOn: []target.Target{
								{Dir: ".", Artefact: "something"},
							},
							Pipeline: "my-pipeline",
						},
						{
							Name: "something",
							Steps: []config.Step{
								{
									Name:    "echo",
									Command: "echo \"hi\"",
								},
								{
									Name:    "build",
									Command: "whatevs",
								},
							},
						},
					},
				},
			},
			Expected: []testNode{
				{
					Path:     ".",
					Artefact: "image",
					SharedEnv: map[string]string{
						"PIPELINE_ENV": "value1",
						"ARTEFACT_ENV": "value2",
						"GLOBAL_ENV":   "value3",
					},
					Step: steps.CommandStep{
						Name:    "build",
						Command: "go build -o ./bin/mudly ./cmd/mudly",
						Env: map[string]string{
							"STEP_ENV": "value0",
						},
					},
					State:     runner.STATE_PENDING,
					DependsOn: []int{3},
				},
				{
					Path:     ".",
					Artefact: "image",
					SharedEnv: map[string]string{
						"PIPELINE_ENV": "value1",
						"ARTEFACT_ENV": "value2",
						"GLOBAL_ENV":   "value3",
					},
					Step: steps.DockerStep{
						Name:       "docker",
						Dockerfile: "./Dockerfile",
						Context:    ".",
					},
					State:     runner.STATE_PENDING,
					DependsOn: []int{0},
				},
				{
					Path:     ".",
					Artefact: "something",
					SharedEnv: map[string]string{
						"GLOBAL_ENV": "value3",
					},
					Step: steps.CommandStep{
						Name:    "echo",
						Command: "echo \"hi\"",
					},
					State:     runner.STATE_PENDING,
					DependsOn: []int{},
				},
				{
					Path:     ".",
					Artefact: "something",
					SharedEnv: map[string]string{
						"GLOBAL_ENV": "value3",
					},
					Step: steps.CommandStep{
						Name:    "build",
						Command: "whatevs",
					},
					State:     runner.STATE_PENDING,
					DependsOn: []int{2},
				},
			},
		},
	} {
		t.Run(test.Name, func(u *testing.T) {
			nodes, err := solver.Solve(test.Targets, test.Configs)

			assert.NoError(u, err, "didn't expect an error")

			if test.Expected != nil {
				if nodes != nil {
					assert.Equal(u, convert(test.Expected), nodes)
				} else {
					assert.Fail(u, "expected a list of nodes")
				}
			}
		})
	}
}
