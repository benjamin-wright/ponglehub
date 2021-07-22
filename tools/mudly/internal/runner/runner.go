package runner

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/tools/mudly/internal/solver"
)

type runResult struct {
	node    *solver.Node
	success bool
}

func Run(nodes []*solver.Node) (err error) {
	numRunning := 0
	outputChan := make(chan runResult)

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
			if dep.State != solver.STATE_COMPLETE {
				runnable = false
			}
		}

		if runnable {
			runnables = append(runnables, node)
		}
	}

	return runnables
}

func runNode(node *solver.Node, outputChan chan<- runResult) {
	logrus.Infof("Running steps %s:%s", node.Artefact, node.Step)
	success := node.Step.Run(node.Artefact, node.SharedEnv)
	if success {
		node.State = solver.STATE_COMPLETE
	} else {
		node.State = solver.STATE_ERROR
	}

	outputChan <- runResult{
		node:    node,
		success: success,
	}
}
