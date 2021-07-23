package runner

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/tools/mudly/internal/solver"
	"ponglehub.co.uk/tools/mudly/internal/steps"
)

type runResult struct {
	node    *solver.Node
	success bool
}

func Run(nodes []*solver.Node) (err error) {
	numRunning := 0
	outputChan := make(chan runResult, 10)

	for {
		pending := getRunnableNodes(nodes)

		for _, node := range pending {
			numRunning += 1
			node.State = solver.STATE_RUNNING
			go runNode(node, outputChan)
		}

		if numRunning == 0 {
			return nil
		}

		result := <-outputChan

		logrus.Infof("Finished step %s:%s", result.node.Artefact, result.node.Step)
		if !result.success {
			return fmt.Errorf("error running step %s:%s", result.node.Artefact, result.node.Step)
		}

		numRunning -= 1
	}
}

func getRunnableNodes(nodes []*solver.Node) []*solver.Node {
	runnables := []*solver.Node{}

	for _, node := range nodes {
		runnable := node.State == solver.STATE_PENDING

		for _, dep := range node.DependsOn {
			if dep.State != solver.STATE_COMPLETE && dep.State != solver.STATE_SKIPPED {
				runnable = false
			}
		}

		if runnable {
			runnables = append(runnables, node)
		}
	}

	return runnables
}

func depsSkipped(node *solver.Node) bool {
	skipped := node.DependsOn != nil && len(node.DependsOn) > 0

	for _, dep := range node.DependsOn {
		if dep.State != solver.STATE_SKIPPED {
			skipped = false
		}
	}

	return skipped
}

func runNode(node *solver.Node, outputChan chan<- runResult) {
	if depsSkipped(node) {
		node.SharedEnv["DEPS_SKIPPED"] = "true"
	}

	logrus.Infof("Running steps %s:%s", node.Artefact, node.Step)
	result := node.Step.Run(node.Path, node.Artefact, node.SharedEnv)
	success := true

	switch result {
	case steps.COMMAND_SUCCESS:
		node.State = solver.STATE_COMPLETE
	case steps.COMMAND_SKIPPED:
		node.State = solver.STATE_SKIPPED
	case steps.COMMAND_ERROR:
		node.State = solver.STATE_ERROR
		success = false
	}

	outputChan <- runResult{
		node:    node,
		success: success,
	}
}
