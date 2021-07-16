package config

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func isCommandStep(n *yaml.Node) bool {
	for _, child := range n.Content {
		if child.Value == "cmd" {
			return true
		}
	}

	return false
}

func isDockerStep(n *yaml.Node) bool {
	for _, child := range n.Content {
		if child.Value == "dockerfile" {
			return true
		}
	}

	return false
}

func (p *Pipeline) UnmarshalYAML(n *yaml.Node) error {
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
			if err := stepNode.Decode(&step); err != nil {
				return err
			}
			p.Steps = append(p.Steps, step)
		} else if isDockerStep(&stepNode) {
			step := DockerStep{}
			if err := stepNode.Decode(&step); err != nil {
				return err
			}
			p.Steps = append(p.Steps, step)
		} else {
			return fmt.Errorf("failed to indentify step type: %+v", stepNode)
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
		pipeline := Pipeline{}
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
	Pipelines []Pipeline       `yaml:"pipelines"`
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
		switch pipeline := artefact.Pipeline.(type) {
		case Pipeline:
			result.Artefacts = append(result.Artefacts, Artefact{
				Name:     artefact.Name,
				Pipeline: pipeline,
			})
		case string:
			result.Artefacts = append(result.Artefacts, Artefact{
				Name: artefact.Name,
				Pipeline: Pipeline{
					Name: pipeline,
				},
			})
		default:
			logrus.Warnf("Failed to process artefact, unknown pipeline type: %+v", artefact)
		}
	}

	return &result, nil
}
