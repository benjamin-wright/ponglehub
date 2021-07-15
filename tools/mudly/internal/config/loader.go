package config

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type DevEnvLoader struct {
	Compose *map[string]interface{} `yaml:"compose"`
}

type PipelineLoader struct {
	Name  string        `yaml:"name"`
	Steps []interface{} `yaml:"-"`
}

func isCommandStep(n *yaml.Node) bool {
	for _, child := range n.Content {
		if child.Value == "cmd" {
			return true
		}
	}

	return false
}

func (p *PipelineLoader) UnmarshalYAML(n *yaml.Node) error {
	type tmpLoader struct {
		Name  string      `yaml:"name"`
		Steps []yaml.Node `yaml:"steps"`
	}

	obj := &tmpLoader{}
	if err := n.Decode(obj); err != nil {
		return err
	}

	p.Name = obj.Name

	for _, stepNode := range obj.Steps {
		if isCommandStep(&stepNode) {
			step := CommandStep{}
			err := stepNode.Decode(&step)
			if err != nil {
				return err
			}

			p.Steps = append(p.Steps, step)
		}
	}

	return nil
}

type ArtefactLoader struct {
	Name     string      `yaml:"name"`
	Pipeline interface{} `yaml:"-"`
}

func (a *ArtefactLoader) UnmarshalYAML(n *yaml.Node) error {
	type tmpLoader struct {
		Name     string    `yaml:"name"`
		Pipeline yaml.Node `yaml:"pipeline"`
	}

	obj := &tmpLoader{}
	if err := n.Decode(obj); err != nil {
		return err
	}

	a.Name = obj.Name

	if obj.Pipeline.Kind == yaml.ScalarNode {
		a.Pipeline = obj.Pipeline.Value
	} else {
		pipeline := PipelineLoader{}
		if err := obj.Pipeline.Decode(&pipeline); err != nil {
			return err
		}
		a.Pipeline = pipeline
	}

	return nil
}

type ConfigLoader struct {
	DevEnv    *DevEnv          `yaml:"devEnv"`
	Artefacts []ArtefactLoader `yaml:"artefacts"`
	Pipelines []PipelineLoader `yaml:"pipelines"`
}

func LoadConfig(data []byte) (*Config, error) {
	var conf ConfigLoader

	err := yaml.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}

	result := Config{}

	if conf.DevEnv != nil {
		result.DevEnv = &DevEnv{
			Compose: conf.DevEnv.Compose,
		}
	}

	for _, artefact := range conf.Artefacts {
		if pipeline, ok := artefact.Pipeline.(PipelineLoader); ok {
			result.Artefacts = append(result.Artefacts, Artefact{
				Name: artefact.Name,
				Pipeline: Pipeline{
					Name:  pipeline.Name,
					Steps: pipeline.Steps,
				},
			})
		} else if pipeline, ok := artefact.Pipeline.(string); ok {
			result.Artefacts = append(result.Artefacts, Artefact{
				Name: artefact.Name,
				Pipeline: Pipeline{
					Name: pipeline,
				},
			})
		} else {
			logrus.Warnf("Failed to process artefact, unknown pipeline type: %+v", artefact)
		}
	}

	return &result, nil
}
