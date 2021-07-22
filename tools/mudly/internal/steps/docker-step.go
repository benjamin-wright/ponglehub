package steps

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

func (d DockerStep) Run(artefact string, env map[string]string) bool {
	return runShellCommand(&shellCommand{
		artefact: artefact,
		step:     d.Name,
		command:  "docker",
		args:     d.args(),
	})
}

func (d DockerStep) String() string {
	return d.Name
}
