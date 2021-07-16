package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestLoader(t *testing.T) {
	for _, test := range []struct {
		Name     string
		Input    string
		Expected *Config
	}{
		{
			Name: "all-in-one",
			Input: dedent(`
				artefacts:
				- name: mudly
				  pipeline:
				    steps:
				    - name: build
				      cmd: go build -o=bin/mudly ./cmd/mudly
					- name: image
					  ignore: [ "**", "!bin/${APP_NAME}" ]
					  context: ./bin
					  dockerfile: ../../dockerfiles/golang
			`),
			Expected: &Config{
				Path: "",
				Artefacts: []Artefact{
					{
						Name: "mudly",
						Pipeline: Pipeline{
							Steps: []interface{}{
								CommandStep{Name: "build", Command: "go build -o=bin/mudly ./cmd/mudly"},
								DockerStep{Name: "image", Ignore: []string{"**", "!bin/${APP_NAME}"}, Context: "./bin", Dockerfile: "../../dockerfiles/golang"},
							},
						},
					},
				},
			},
		},
	} {
		t.Run(test.Name, func(u *testing.T) {
			conf, err := LoadConfig([]byte(test.Input))

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
