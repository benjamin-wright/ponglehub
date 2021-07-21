package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"

	"gopkg.in/yaml.v3"
	"ponglehub.co.uk/tools/mudly/internal/steps"
	"ponglehub.co.uk/tools/mudly/internal/target"
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
		Name  string            `yaml:"name"`
		Steps []yaml.Node       `yaml:"steps"`
		Env   map[string]string `yaml:"env"`
	}

	obj := &tmpLoader{}
	if err := n.Decode(obj); err != nil {
		return err
	}

	p.Name = obj.Name
	p.Env = obj.Env

	for _, stepNode := range obj.Steps {
		if isCommandStep(&stepNode) {
			step := steps.CommandStep{}
			if err := stepNode.Decode(&step); err != nil {
				return err
			}
			p.Steps = append(p.Steps, step)
		} else if isDockerStep(&stepNode) {
			step := steps.DockerStep{}
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
	Name         string            `yaml:"name"`
	Pipeline     interface{}       `yaml:"-"`
	Dependencies []string          `yaml:"dependencies"`
	Env          map[string]string `yaml:"env"`
}

func (a *ArtefactLoader) UnmarshalYAML(n *yaml.Node) error {
	type tmpLoader struct {
		Name         string            `yaml:"name"`
		Pipeline     yaml.Node         `yaml:"pipeline"`
		Dependencies []string          `yaml:"dependencies"`
		Env          map[string]string `yaml:"env"`
	}

	obj := &tmpLoader{}
	if err := n.Decode(obj); err != nil {
		return err
	}

	a.Name = obj.Name
	a.Dependencies = obj.Dependencies
	a.Env = obj.Env

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
	DevEnv    *DevEnv           `yaml:"devEnv"`
	Artefacts []ArtefactLoader  `yaml:"artefacts"`
	Pipelines []Pipeline        `yaml:"pipelines"`
	Env       map[string]string `yaml:"env"`
}

type FileSystem interface {
	ReadFile(path string) ([]byte, error)
}

type DefaultFS struct{}

func (fs DefaultFS) ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

type LoadConfigOptions struct {
	Targets []target.Target
	FS      FileSystem
}

func loadConfigFromFile(filepath string, filesystem FileSystem) (*Config, error) {
	var loader ConfigLoader
	var fs FileSystem
	if filesystem != nil {
		fs = filesystem
	} else {
		fs = DefaultFS{}
	}

	data, err := fs.ReadFile(fmt.Sprintf("%s/mudly.yaml", filepath))
	if err != nil {
		log.Printf("Error loading config from file %s: %+v ", filepath, err)
	}

	err = yaml.Unmarshal(data, &loader)
	if err != nil {
		return nil, err
	}

	config := Config{
		Path: path.Clean(filepath),
		Env:  loader.Env,
	}

	if loader.DevEnv != nil {
		config.DevEnv = &DevEnv{
			Compose: loader.DevEnv.Compose,
		}
	}

	for _, artefact := range loader.Artefacts {
		dependencies := []target.Target{}

		for _, targetString := range artefact.Dependencies {
			dependency, err := target.ParseTarget(targetString)
			if err != nil {
				return nil, fmt.Errorf("failed parsing dependency target: %+v", err)
			}

			dependencies = append(dependencies, *dependency)
		}

		var resolvedPipeline Pipeline

		switch pipeline := artefact.Pipeline.(type) {
		case Pipeline:
			resolvedPipeline = pipeline
		case string:
			missing := true

			for _, external := range loader.Pipelines {
				if external.Name == pipeline {
					resolvedPipeline = external
					missing = false
					break
				}
			}

			if missing {
				return nil, fmt.Errorf("failed to resolve pipeline %s", pipeline)
			}
		default:
			return nil, fmt.Errorf("failed to process artefact, unknown pipeline type: %+v", artefact)
		}

		config.Artefacts = append(config.Artefacts, Artefact{
			Name:         artefact.Name,
			Pipeline:     resolvedPipeline,
			Dependencies: dependencies,
			Env:          artefact.Env,
		})
	}

	return &config, nil
}

func getDependencyTargets(config *Config) []target.Target {
	targets := []target.Target{}

	// Resolve dependency configs and add to the list
	for _, artefact := range config.Artefacts {
		for _, dependency := range artefact.Dependencies {
			targets = append(targets, target.Target{
				Dir:      path.Clean(fmt.Sprintf("%s/%s", config.Path, dependency.Dir)),
				Artefact: dependency.Artefact,
			})
		}
	}

	return targets
}

func getNewDependencies(configs []Config) []target.Target {
	targets := []target.Target{}

	for _, config := range configs {
		targets = append(targets, getDependencyTargets(&config)...)
	}

	newTargets := []target.Target{}
	for _, target := range targets {
		isNew := true

		for _, config := range configs {
			if config.Path == target.Dir {
				isNew = false
				break
			}
		}

		if isNew {
			newTargets = append(newTargets, target)
		}
	}

	return newTargets
}

func dedupConfigs(configs []Config) []Config {
	result := []Config{}

	for _, config := range configs {
		add := true
		for _, existing := range result {
			if config.Path == existing.Path {
				add = false
				break
			}
		}

		if add {
			result = append(result, config)
		}
	}

	return result
}

func LoadConfig(options *LoadConfigOptions) ([]Config, error) {
	configs := []Config{}
	targets := options.Targets
	running := true

	for running {
		if len(targets) == 0 {
			running = false
			continue
		}

		for _, target := range targets {
			config, err := loadConfigFromFile(target.Dir, options.FS)
			if err != nil {
				return nil, err
			}

			configs = append(configs, *config)
		}

		targets = getNewDependencies(configs)
	}

	return dedupConfigs(configs), nil
}
