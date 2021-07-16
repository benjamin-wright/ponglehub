package solver

import (
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
	Step      string
	Command   string
	State     NodeState
	DependsOn []*Node
}

func Solve(targets []target.Target, configs []config.Config) ([]Node, error) {

	return nil, nil
}
