package config

type DevEnv struct {
	Compose *map[string]interface{}
}

type Pipeline struct {
	Name  string
	Steps []interface{}
}

type CommandStep struct {
	Name    string `yaml:"name"`
	Command string `yaml:"cmd"`
}

type Artefact struct {
	Name     string
	Pipeline Pipeline
}

type Config struct {
	DevEnv    *DevEnv
	Path      string
	Artefacts []Artefact
}
