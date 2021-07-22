package steps

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type CommandStep struct {
	Name      string            `yaml:"name"`
	Condition string            `yaml:"condition"`
	Command   string            `yaml:"cmd"`
	Env       map[string]string `yaml:"env"`
}

func (c CommandStep) Run(artefact string, env map[string]string) bool {
	if c.Condition != "" {
		test := runShellCommand(&shellCommand{
			artefact: artefact,
			step:     fmt.Sprintf("%s (test)", c.Name),
			command:  "/bin/bash",
			args:     []string{"-c", c.Condition},
			env:      c.Env,
			test:     true,
		})

		if test {
			logrus.Infof("%s[%s (test)]: Skipping step", artefact, c.Name)
			return true
		}
	}

	return runShellCommand(&shellCommand{
		artefact: artefact,
		step:     c.Name,
		command:  "/bin/bash",
		args:     []string{"-c", c.Command},
		env:      c.Env,
	})
}

func (c CommandStep) String() string {
	return c.Name
}
