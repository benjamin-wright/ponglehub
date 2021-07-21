package config

import "ponglehub.co.uk/tools/mudly/internal/target"

type Runnable interface {
	Run(env map[string]string) bool
}

type DevEnv struct {
	Compose *map[string]interface{} `yaml:"compose"`
}

type Pipeline struct {
	Name  string            `yaml:"name"`
	Steps []Runnable        `yaml:"-"`
	Env   map[string]string `yaml:"env"`
}

type Artefact struct {
	Name         string            `yaml:"name"`
	Pipeline     Pipeline          `yaml:"pipeline"`
	Dependencies []target.Target   `yaml:"dependencies"`
	Env          map[string]string `yaml:"env"`
}

type Config struct {
	DevEnv    *DevEnv
	Path      string
	Artefacts []Artefact
	Env       map[string]string `yaml:"env"`
}
