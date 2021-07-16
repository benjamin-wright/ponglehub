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

func Solve(targets []target.Target, configs []config.Config) ([]Node, error) {
	nodes := []Node{}

	for _, target := range targets {
		cfg, artefact, err := getArtefact(target, configs)
		if err != nil {
			return nil, err
		}

		pipeline := &artefact.Pipeline
		var previous *Node
		for _, step := range pipeline.Steps {
			dependsOn := []*Node{}
			if previous != nil {
				dependsOn = append(dependsOn, previous)
			}

			newNode := Node{
				Path:      cfg.Path,
				Artefact:  artefact.Name,
				Step:      step,
				State:     STATE_NONE,
				DependsOn: dependsOn,
			}

			previous = &newNode
			nodes = append(nodes, newNode)
		}
	}

	return nodes, nil
}
