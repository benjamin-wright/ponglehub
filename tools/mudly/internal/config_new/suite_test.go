package config_new_test

import (
	"errors"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"ponglehub.co.uk/tools/mudly/internal/config_new"
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
		Expected []config_new.Config
	}{
		{
			Name: "all-in-one",
			Files: []ConfigFile{
				{
					Path: "./Mudfile",
					Content: dedent(`
						ENV GLOBAL_VAR=value1

						ARTEFACT my-artefact
					`),
				},
			},
			Targets: []target.Target{{Dir: "."}},
			Expected: []config_new.Config{
				{
					Path: ".",
					Env: map[string]string{
						"GLOBAL_VAR": "value1",
					},
					Artefacts: []config_new.Artefact{
						{
							Name: "my-artefact",
						},
					},
				},
			},
		},
	} {
		t.Run(test.Name, func(u *testing.T) {
			config_new.SetFS(
				&MockFS{
					files: test.Files,
				},
			)

			conf, err := config_new.LoadConfigs(test.Targets)

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
