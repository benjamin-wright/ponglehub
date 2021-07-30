package config_new

import "ponglehub.co.uk/tools/mudly/internal/target"

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
	DependsOn []target.Target
	Env       map[string]string
	Steps     []Step
	Pipeline  string
}

type Step struct {
	Name       string
	Env        map[string]string
	Condition  string
	Watch      []string
	Dommand    string
	Dockerfile string
}
