package config

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"ponglehub.co.uk/tools/mudly/internal/steps"
	"ponglehub.co.uk/tools/mudly/internal/target"
)

func dedent(file string) string {
	expanded := strings.Replace(file, "\t", "    ", -1)

	lines := strings.Split(expanded, "\n")

	whitespace := 0
	for idx, char := range lines[1] {
		if char != ' ' {
			whitespace = idx
			break
		}
	}

	trimmed := []string{}
	for _, line := range lines {
		if len(line) > whitespace {
			trimmed = append(trimmed, line[whitespace:])
		} else {
			trimmed = append(trimmed, "")
		}
	}

	return strings.Join(trimmed, "\n")
}

type ConfigFile struct {
	Path    string
	Content string
}

type MockFS struct {
	files []ConfigFile
}

func (m MockFS) ReadFile(path string) ([]byte, error) {
	for _, file := range m.files {
		if file.Path == path {
			return []byte(file.Content), nil
		}
	}

	return nil, errors.New("File not found")
}

func TestLoader(t *testing.T) {
	for _, test := range []struct {
		Name     string
		Files    []ConfigFile
		Targets  []target.Target
		Expected []Config
	}{
		{
			Name: "all-in-one",
			Files: []ConfigFile{
				{
					Path: "./mudly.yaml",
					Content: dedent(`
						env:
						  GLOBAL_VAR: value4
						artefacts:
						- name: mudly
						  env:
						    ARTEFACT_VAR: value3
						  pipeline:
							env:
							  SHARED_VAR: value2
						    steps:
						    - name: build
						      cmd: go build -o=bin/mudly ./cmd/mudly
							  env:
							  	SOME_VAR: value
						    - name: image
						      ignore: [ "**", "!bin/mudly" ]
						      context: ./bin
						      dockerfile: ../../dockerfiles/golang
					`),
				},
			},
			Targets: []target.Target{{Dir: "."}},
			Expected: []Config{
				{
					Path: ".",
					Env: map[string]string{
						"GLOBAL_VAR": "value4",
					},
					Artefacts: []Artefact{
						{
							Name:         "mudly",
							Dependencies: []target.Target{},
							Env: map[string]string{
								"ARTEFACT_VAR": "value3",
							},
							Pipeline: Pipeline{
								Env: map[string]string{"SHARED_VAR": "value2"},
								Steps: []Runnable{
									steps.CommandStep{Name: "build", Command: "go build -o=bin/mudly ./cmd/mudly", Env: map[string]string{"SOME_VAR": "value"}},
									steps.DockerStep{Name: "image", Context: "./bin", Dockerfile: "../../dockerfiles/golang"},
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "external pipelines",
			Files: []ConfigFile{
				{
					Path: "./mudly.yaml",
					Content: dedent(`
						artefacts:
						- name: mudly
						  pipeline: external
						pipelines:
						- name: external
						  steps:
						  - name: build
						    cmd: go build -o=bin/${ARTEFACT} ./cmd/${ARTEFACT}
						  - name: image
						    ignore: [ "**", "!bin/${ARTEFACT}" ]
						    context: ./bin
						    dockerfile: ../../dockerfiles/golang
					`),
				},
			},
			Targets: []target.Target{{Dir: "."}},
			Expected: []Config{
				{
					Path: ".",
					Artefacts: []Artefact{
						{
							Name:         "mudly",
							Dependencies: []target.Target{},
							Pipeline: Pipeline{
								Name: "external",
								Steps: []Runnable{
									steps.CommandStep{Name: "build", Command: "go build -o=bin/${ARTEFACT} ./cmd/${ARTEFACT}"},
									steps.DockerStep{Name: "image", Context: "./bin", Dockerfile: "../../dockerfiles/golang"},
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "linked dependencies",
			Files: []ConfigFile{
				{
					Path: "subdir1/mudly.yaml",
					Content: dedent(`
						artefacts:
						- name: mudly1
						  dependencies:
						  - ../subdir2+mudly2
						  pipeline:
						    steps:
						    - name: build
						      cmd: go build -o=bin/mudly ./cmd/mudly
						    - name: image
						      ignore: [ "**", "!bin/mudly" ]
						      context: ./bin
						      dockerfile: ../../dockerfiles/golang
					`),
				},
				{
					Path: "subdir2/mudly.yaml",
					Content: dedent(`
						artefacts:
						- name: mudly2
						  pipeline:
						    steps:
						    - name: build
						      cmd: go build -o=bin/mudly2 ./cmd/mudly2
						    - name: image
						      ignore: [ "**", "!bin/mudly2" ]
						      context: ./bin
						      dockerfile: ../../dockerfiles/golang
					`),
				},
			},
			Targets: []target.Target{{Dir: "subdir1"}},
			Expected: []Config{
				{
					Path: "subdir1",
					Artefacts: []Artefact{
						{
							Name:         "mudly1",
							Dependencies: []target.Target{{Dir: "../subdir2", Artefact: "mudly2"}},
							Pipeline: Pipeline{
								Steps: []Runnable{
									steps.CommandStep{Name: "build", Command: "go build -o=bin/mudly ./cmd/mudly"},
									steps.DockerStep{Name: "image", Context: "./bin", Dockerfile: "../../dockerfiles/golang"},
								},
							},
						},
					},
				},
				{
					Path: "subdir2",
					Artefacts: []Artefact{
						{
							Name:         "mudly2",
							Dependencies: []target.Target{},
							Pipeline: Pipeline{
								Steps: []Runnable{
									steps.CommandStep{Name: "build", Command: "go build -o=bin/mudly2 ./cmd/mudly2"},
									steps.DockerStep{Name: "image", Context: "./bin", Dockerfile: "../../dockerfiles/golang"},
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "self-reference",
			Files: []ConfigFile{
				{
					Path: "subdir1/mudly.yaml",
					Content: dedent(`
						artefacts:
						- name: mudly1
						  dependencies:
						  - +mudly2
						  pipeline: external
						- name: mudly2
						  pipeline: external
						pipelines:
						- name: external
						  steps:
						  - name: build
						    cmd: go build -o=bin/${ARTEFACT} ./cmd/${ARTEFACT}
						  - name: image
						    ignore: [ "**", "!bin/${ARTEFACT}" ]
						    context: ./bin
						    dockerfile: ../../dockerfiles/golang
					`),
				},
			},
			Targets: []target.Target{{Dir: "subdir1"}},
			Expected: []Config{
				{
					Path: "subdir1",
					Artefacts: []Artefact{
						{
							Name:         "mudly1",
							Dependencies: []target.Target{{Dir: ".", Artefact: "mudly2"}},
							Pipeline: Pipeline{
								Name: "external",
								Steps: []Runnable{
									steps.CommandStep{Name: "build", Command: "go build -o=bin/${ARTEFACT} ./cmd/${ARTEFACT}"},
									steps.DockerStep{Name: "image", Context: "./bin", Dockerfile: "../../dockerfiles/golang"},
								},
							},
						},
						{
							Name:         "mudly2",
							Dependencies: []target.Target{},
							Pipeline: Pipeline{
								Name: "external",
								Steps: []Runnable{
									steps.CommandStep{Name: "build", Command: "go build -o=bin/${ARTEFACT} ./cmd/${ARTEFACT}"},
									steps.DockerStep{Name: "image", Context: "./bin", Dockerfile: "../../dockerfiles/golang"},
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "multi-target",
			Files: []ConfigFile{
				{
					Path: "subdir1/mudly.yaml",
					Content: dedent(`
						artefacts:
						- name: mudly1
						  pipeline:
						    steps:
						    - name: build
						      cmd: go build -o=bin/mudly ./cmd/mudly
						    - name: image
						      ignore: [ "**", "!bin/mudly" ]
						      context: ./bin
						      dockerfile: ../../dockerfiles/golang
					`),
				},
				{
					Path: "subdir2/mudly.yaml",
					Content: dedent(`
						artefacts:
						- name: mudly2
						  pipeline:
						    steps:
						    - name: build
						      cmd: go build -o=bin/mudly2 ./cmd/mudly2
						    - name: image
						      ignore: [ "**", "!bin/mudly2" ]
						      context: ./bin
						      dockerfile: ../../dockerfiles/golang
					`),
				},
			},
			Targets: []target.Target{{Dir: "subdir1"}, {Dir: "subdir2"}},
			Expected: []Config{
				{
					Path: "subdir1",
					Artefacts: []Artefact{
						{
							Name:         "mudly1",
							Dependencies: []target.Target{},
							Pipeline: Pipeline{
								Steps: []Runnable{
									steps.CommandStep{Name: "build", Command: "go build -o=bin/mudly ./cmd/mudly"},
									steps.DockerStep{Name: "image", Context: "./bin", Dockerfile: "../../dockerfiles/golang"},
								},
							},
						},
					},
				},
				{
					Path: "subdir2",
					Artefacts: []Artefact{
						{
							Name:         "mudly2",
							Dependencies: []target.Target{},
							Pipeline: Pipeline{
								Steps: []Runnable{
									steps.CommandStep{Name: "build", Command: "go build -o=bin/mudly2 ./cmd/mudly2"},
									steps.DockerStep{Name: "image", Context: "./bin", Dockerfile: "../../dockerfiles/golang"},
								},
							},
						},
					},
				},
			},
		},
	} {
		t.Run(test.Name, func(u *testing.T) {
			conf, err := LoadConfig(&LoadConfigOptions{
				Targets: test.Targets,
				FS:      &MockFS{files: test.Files},
			})

			assert.NoError(u, err, "didn't expect an error")

			if test.Expected != nil {
				if conf != nil {
					assert.Equal(u, test.Expected, conf)
				} else {
					assert.Fail(u, "expected a config response")
				}
			}
		})
	}
}
