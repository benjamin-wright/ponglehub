package config_new

type Dockerfile struct {
	Name    string
	Content string
}

type Config struct {
	Path       string
	Dockerfile []Dockerfile
	Artefacts  []Artefact
	Pipelines  []Pipeline
	Env        map[string]string
}

type Pipeline struct {
	Name  string
	Steps []Step
}

type Artefact struct {
	Name      string
	DependsOn []string
	Env       map[string]string
	Steps     []Step
	Pipeline  string
}

type Step struct {
	Name       string
	Env        map[string]string
	Watch      []string
	Dommand    string
	Dockerfile string
}
