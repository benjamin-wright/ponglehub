package config

type DevEnv struct {
	Compose *map[string]interface{} `yaml:"compose"`
}

type Pipeline struct {
	Name  string        `yaml:"name"`
	Steps []interface{} `yaml:"-"`
}

type CommandStep struct {
	Name    string `yaml:"name"`
	Command string `yaml:"cmd"`
}

type DockerStep struct {
	Name       string   `yaml:"name"`
	Dockerfile string   `yaml:"dockerfile"`
	Ignore     []string `yaml:"ignore"`
	Context    string   `yaml:"context"`
}

type Artefact struct {
	Name     string   `yaml:"name"`
	Pipeline Pipeline `yaml:"pipeline"`
}

type Config struct {
	DevEnv    *DevEnv
	Path      string
	Artefacts []Artefact
}
