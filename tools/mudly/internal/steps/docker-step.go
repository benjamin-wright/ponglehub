package steps

type DockerStep struct {
	Name       string `yaml:"name"`
	Dockerfile string `yaml:"dockerfile"`
	Context    string `yaml:"context"`
	Tag        string `yaml:"tag"`
}

func (d DockerStep) Run() bool {
	return true
}

func (d DockerStep) String() string {
	return d.Name
}
