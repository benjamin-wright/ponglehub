package solver

import (
	"fmt"

	"ponglehub.co.uk/tools/mudly/internal/config"
	"ponglehub.co.uk/tools/mudly/internal/target"
)

type NodeState int

const (
	STATE_NONE NodeState = iota
	STATE_PENDING
	STATE_READY
	STATE_RUNNING
	STATE_ERROR
	STATE_COMPLETE
)

type Node struct {
	Path      string
	Artefact  string
	Step      interface{}
	State     NodeState
	DependsOn []*Node
}

func getArtefact(target target.Target, configs []config.Config) (*config.Config, *config.Artefact, error) {
	var cfg config.Config
	missing := true
	for _, c := range configs {
		if target.Dir == c.Path {
			cfg = c
			missing = false
			break
		}
	}

	if missing {
		return nil, nil, fmt.Errorf("couldn't find target config %s", target.Dir)
	}

	var artefact config.Artefact
	missing = true
	for _, a := range cfg.Artefacts {
		if a.Name == target.Artefact {
			artefact = a
			missing = false
			break
		}
	}

	if missing {
		return nil, nil, fmt.Errorf("couldn't find target artefact %s+%s", target.Dir, target.Artefact)
	}

	return &cfg, &artefact, nil
}

func collectTargets(targets []target.Target, configs []config.Config) ([]target.Target, error) {
	running := true
	for running {
		newTargets := []target.Target{}

		for _, target := range targets {
			_, artefact, err := getArtefact(target, configs)
			if err != nil {
				return nil, err
			}

			for _, dep := range artefact.Dependencies {
				rebased := dep.Rebase(target)

				missing := true

				for _, existing := range targets {
					if rebased.Dir == existing.Dir {
						missing = false
						break
					}
				}

				if missing {
					newTargets = append(newTargets, rebased)
				}
			}
		}

		if len(newTargets) == 0 {
			running = false
			continue
		}

		targets = append(targets, newTargets...)
	}

	output := []target.Target{}

	for _, target := range targets {
		missing := true

		for _, existing := range output {
			if target.Dir == existing.Dir {
				missing = false
				break
			}
		}

		if missing {
			output = append(output, target)
		}
	}

	return output, nil
}

func createNodes(targets []target.Target, configs []config.Config) (NodeList, error) {
	nodes := NodeList{list: []nodeListElement{}}

	for _, target := range targets {
		cfg, artefact, err := getArtefact(target, configs)
		if err != nil {
			return nodes, err
		}

		pipeline := &artefact.Pipeline
		for _, step := range pipeline.Steps {
			newNode := Node{
				Path:      cfg.Path,
				Artefact:  artefact.Name,
				Step:      step,
				State:     STATE_NONE,
				DependsOn: []*Node{},
			}

			nodes.AddNode(cfg.Path, artefact.Name, newNode)
		}
	}

	return nodes, nil
}

func Solve(targets []target.Target, configs []config.Config) ([]Node, error) {
	targets, err := collectTargets(targets, configs)
	if err != nil {
		return nil, err
	}

	nodes, err := createNodes(targets, configs)
	if err != nil {
		return nil, err
	}

	return nodes.GetList(), nil
}
