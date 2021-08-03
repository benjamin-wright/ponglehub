package steps

import "ponglehub.co.uk/tools/mudly/internal/runner"

type DockerStep struct {
	Name       string `yaml:"name"`
	Dockerfile string `yaml:"dockerfile"`
	Context    string `yaml:"context"`
	Tag        string `yaml:"tag"`
}

func (d DockerStep) args() []string {
	args := []string{}
	if d.Dockerfile != "" {
		args = append(args, "-f", d.Dockerfile)
	}

	if d.Tag != "" {
		args = append(args, "-t", d.Tag)
	}

	if d.Context != "" {
		args = append(args, d.Context)
	} else {
		args = append(args, ".")
	}

	return args
}

func (d DockerStep) Run(dir string, artefact string, env map[string]string) runner.CommandResult {
	success := runShellCommand(&shellCommand{
		dir:      dir,
		artefact: artefact,
		step:     d.Name,
		command:  "docker",
		args:     d.args(),
	})

	if success {
		return runner.COMMAND_SUCCESS
	} else {
		return runner.COMMAND_ERROR
	}
}

func (d DockerStep) String() string {
	return d.Name
}
