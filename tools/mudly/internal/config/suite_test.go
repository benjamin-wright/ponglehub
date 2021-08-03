package config_test

import (
	"errors"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"ponglehub.co.uk/tools/mudly/internal/config"
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

func (m MockFS) ReadFile(filepath string) ([]byte, error) {
	for _, file := range m.files {
		if path.Clean(file.Path) == filepath {
			return []byte(file.Content), nil
		}
	}

	return nil, errors.New("File not found")
}

func TestLoadConfig(t *testing.T) {

	for _, test := range []struct {
		Name     string
		Files    []ConfigFile
		Targets  []target.Target
		Expected []config.Config
	}{
		{
			Name: "all-in-one",
			Files: []ConfigFile{
				{
					Path: "../somedir/Mudfile",
					Content: dedent(`
                        ARTEFACT not-referenced
                          STEP do-something-else
                            COMMAND echo "ho"
                        
                        ARTEFACT other-artefact
                          DEPENDS ON +not-referenced
                          STEP do-something
                            COMMAND echo "hi"
                    `),
				},
				{
					Path: "./Mudfile",
					Content: dedent(`
                        ENV GLOBAL_VAR=value1
                        
                        ARTEFACT my-artefact
                          ENV ART_VAL=value2
                          DEPENDS ON ../somedir+other-artefact

                          STEP test
                            ENV STEP_VAR=value3
                            WATCH ./path1 ./path2
                            CONDITION echo "inline script"
                            COMMAND echo "inline command"
                        
                          STEP multiline
                            COMMAND
                              echo "multiline"
                              echo "command"
                                # random comment

                          STEP image
                            CONDITION
                              echo "multiline"
                                # indented
                              echo "script"
                            DOCKERFILE file-name
							CONTEXT ./
							TAG localhost:5000/my-image

                        ENV G_2_VAR=var2
                        
                        ARTEFACT local-pipeline
                          PIPELINE my-pipeline
                        
                        PIPELINE my-pipeline
                          ENV P_VAR=var-p
                          STEP step-1
                            COMMAND do the thing
                          STEP step-2
                            COMMAND do the other thing
                        
                        DOCKERFILE file-name
                          FROM something
                          RUN hello there
                    `),
				},
			},
			Targets: []target.Target{{Dir: "."}},
			Expected: []config.Config{
				{
					Path: ".",
					Env: map[string]string{
						"GLOBAL_VAR": "value1",
						"G_2_VAR":    "var2",
					},
					Artefacts: []config.Artefact{
						{
							Name: "my-artefact",
							Env: map[string]string{
								"ART_VAL": "value2",
							},
							DependsOn: []target.Target{
								{Dir: "../somedir", Artefact: "other-artefact"},
							},
							Steps: []config.Step{
								{
									Name: "test",
									Env: map[string]string{
										"STEP_VAR": "value3",
									},
									Condition: "echo \"inline script\"",
									Command:   "echo \"inline command\"",
									Watch: []string{
										"./path1",
										"./path2",
									},
								},
								{
									Name:    "multiline",
									Command: "echo \"multiline\"\necho \"command\"\n  # random comment",
								},
								{
									Name:       "image",
									Condition:  "echo \"multiline\"\n  # indented\necho \"script\"",
									Dockerfile: "file-name",
									Context:    "./",
									Tag:        "localhost:5000/my-image",
								},
							},
						},
						{
							Name:     "local-pipeline",
							Pipeline: "my-pipeline",
						},
					},
					Dockerfile: []config.Dockerfile{
						{
							Name:    "file-name",
							Content: "FROM something\nRUN hello there",
						},
					},
					Pipelines: []config.Pipeline{
						{
							Name: "my-pipeline",
							Env: map[string]string{
								"P_VAR": "var-p",
							},
							Steps: []config.Step{
								{Name: "step-1", Command: "do the thing"},
								{Name: "step-2", Command: "do the other thing"},
							},
						},
					},
				},
				{
					Path: "../somedir",
					Artefacts: []config.Artefact{
						{
							Name: "not-referenced",
							Steps: []config.Step{
								{
									Name:    "do-something-else",
									Command: "echo \"ho\"",
								},
							},
						},
						{
							Name: "other-artefact",
							DependsOn: []target.Target{
								{Dir: ".", Artefact: "not-referenced"},
							},
							Steps: []config.Step{
								{
									Name:    "do-something",
									Command: "echo \"hi\"",
								},
							},
						},
					},
				},
			},
		},
	} {
		t.Run(test.Name, func(u *testing.T) {
			config.SetFS(
				&MockFS{
					files: test.Files,
				},
			)

			conf, err := config.LoadConfigs(test.Targets)

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
