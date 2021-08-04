package config

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
	Env   map[string]string
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
	Command    string
	Watch      []string
	Dockerfile string
	Context    string
	Tag        string
	WaitFor    []string
}
