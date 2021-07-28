package config_new

import "ponglehub.co.uk/tools/mudly/internal/target"

type Step struct {
	Name       string
	Watch      []string
	Command    string
	Condition  string
	Dockerfile string
}

type Pipeline struct {
	Name  string
	Steps []Step
	Env   map[string]string
}

type Artefact struct {
	Name      string
	Pipeline  Pipeline
	DependsOn []target.Target
	Env       map[string]string
}

type Config struct {
	Path      string
	Artefacts []Artefact
	Env       map[string]string
}
