package solver_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"ponglehub.co.uk/tools/mudly/internal/config"
	"ponglehub.co.uk/tools/mudly/internal/solver"
	"ponglehub.co.uk/tools/mudly/internal/steps"
	"ponglehub.co.uk/tools/mudly/internal/target"
)

type testNode struct {
	Path      string
	Artefact  string
	SharedEnv map[string]string
	Step      config.Runnable
	State     solver.NodeState
	DependsOn []int
}

func convert(nodes []testNode) []*solver.Node {
	converted := []*solver.Node{}

	for id := range nodes {
		node := nodes[id]
		converted = append(converted, &solver.Node{
			Path:      node.Path,
			Artefact:  node.Artefact,
			SharedEnv: node.SharedEnv,
			Step:      node.Step,
			State:     node.State,
			DependsOn: []*solver.Node{},
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
							Pipeline: config.Pipeline{
								Steps: []config.Runnable{
									steps.CommandStep{
										Name:    "build",
										Command: "go build -o ./bin/mudly ./cmd/mudly",
									},
									steps.DockerStep{
										Name:       "docker",
										Dockerfile: "./Dockerfile",
										Context:    ".",
									},
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
					State:     solver.STATE_PENDING,
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
					State:     solver.STATE_PENDING,
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
							Pipeline: config.Pipeline{
								Steps: []config.Runnable{
									steps.CommandStep{
										Name:    "build",
										Command: "go build -o ./bin/mudly ./cmd/mudly",
									},
									steps.DockerStep{
										Name:       "docker",
										Dockerfile: "./Dockerfile",
										Context:    ".",
									},
								},
							},
						},
						{
							Name: "something",
							Pipeline: config.Pipeline{
								Steps: []config.Runnable{
									steps.CommandStep{
										Name:    "echo",
										Command: "echo \"hi\"",
									},
									steps.CommandStep{
										Name:    "build",
										Command: "whatevs",
									},
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
					State:     solver.STATE_PENDING,
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
					State:     solver.STATE_PENDING,
					DependsOn: []int{0},
				},
				{
					Path:     ".",
					Artefact: "something",
					Step: steps.CommandStep{
						Name:    "echo",
						Command: "echo \"hi\"",
					},
					State:     solver.STATE_PENDING,
					DependsOn: []int{},
				},
				{
					Path:     ".",
					Artefact: "something",
					Step: steps.CommandStep{
						Name:    "build",
						Command: "whatevs",
					},
					State:     solver.STATE_PENDING,
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
					Artefacts: []config.Artefact{
						{
							Name: "image",
							Env: map[string]string{
								"ARTEFACT_ENV": "value2",
							},
							Dependencies: []target.Target{
								{Dir: ".", Artefact: "something"},
							},
							Pipeline: config.Pipeline{
								Env: map[string]string{
									"PIPELINE_ENV": "value1",
								},
								Steps: []config.Runnable{
									steps.CommandStep{
										Name:    "build",
										Command: "go build -o ./bin/mudly ./cmd/mudly",
										Env: map[string]string{
											"STEP_ENV": "value0",
										},
									},
									steps.DockerStep{
										Name:       "docker",
										Dockerfile: "./Dockerfile",
										Context:    ".",
									},
								},
							},
						},
						{
							Name: "something",
							Pipeline: config.Pipeline{
								Steps: []config.Runnable{
									steps.CommandStep{
										Name:    "echo",
										Command: "echo \"hi\"",
									},
									steps.CommandStep{
										Name:    "build",
										Command: "whatevs",
									},
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
					State:     solver.STATE_PENDING,
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
					State:     solver.STATE_PENDING,
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
					State:     solver.STATE_PENDING,
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
					State:     solver.STATE_PENDING,
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
