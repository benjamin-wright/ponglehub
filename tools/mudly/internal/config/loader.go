package config

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type DevEnvLoader struct {
	Compose *map[string]interface{} `yaml:"compose"`
}

type PipelineLoader struct {
	Name string `yaml:"name"`
}

type ArtefactLoader struct {
	Name     string `yaml:"name"`
	Pipeline interface{}
}

type ConfigLoader struct {
	DevEnv    *DevEnv    `yaml:"devEnv"`
	Artefacts []Artefact `yaml:"artefacts"`
	Pipelines []Pipeline `yaml:"pipelines"`
}

func LoadConfig(path string) ([]Config, error) {
	yamlFile, err := ioutil.ReadFile(fmt.Sprintf("%s/mudly.yaml", path))
	if err != nil {
		log.Printf("Error loading config from file %s: %+v ", path, err)
	}

	var conf ConfigLoader

	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		return []Config{}, err
	}

	result := Config{}

	if conf.DevEnv != nil {
		result.DevEnv = &DevEnv{
			Compose: conf.DevEnv.Compose,
		}
	}

	return []Config{result}, nil
}
