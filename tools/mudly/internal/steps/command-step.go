package steps

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/tools/mudly/internal/utils"
)

type CommandStep struct {
	Name      string            `yaml:"name"`
	Condition string            `yaml:"condition"`
	Command   string            `yaml:"cmd"`
	Env       map[string]string `yaml:"env"`
}

func (c CommandStep) Run(dir string, artefact string, env map[string]string) CommandResult {
	merged := utils.MergeMaps(env, c.Env)

	if c.Condition != "" {
		test := runShellCommand(&shellCommand{
			dir:      dir,
			artefact: artefact,
			step:     fmt.Sprintf("%s (test)", c.Name),
			command:  "/bin/bash",
			args:     []string{"-c", c.Condition},
			env:      merged,
			test:     true,
		})

		if !test {
			logrus.Infof("%s[%s (test)]: Skipping step", artefact, c.Name)
			return COMMAND_SKIPPED
		}
	}

	success := runShellCommand(&shellCommand{
		dir:      dir,
		artefact: artefact,
		step:     c.Name,
		command:  "/bin/bash",
		args:     []string{"-c", c.Command},
		env:      merged,
	})

	if success {
		return COMMAND_SUCCESS
	} else {
		return COMMAND_ERROR
	}
}

func (c CommandStep) String() string {
	return c.Name
}
