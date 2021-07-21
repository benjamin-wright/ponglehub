package solver

import (
	"fmt"
	"path"

	"ponglehub.co.uk/tools/mudly/internal/config"
	"ponglehub.co.uk/tools/mudly/internal/target"
)

type NodeState int

const (
	STATE_PENDING NodeState = iota
	STATE_RUNNING
	STATE_ERROR
	STATE_COMPLETE
)

type Node struct {
	SharedEnv map[string]string
	Path      string
	Artefact  string
	Step      config.Runnable
	State     NodeState
	DependsOn []*Node
}

type Link struct {
	Target target.Target
	Source target.Target
}

func (l Link) isSame(m Link) bool {
	return l.Source.IsSame(m.Source) && l.Target.IsSame(m.Target)
}

func getArtefact(target target.Target, configs []config.Config) (*config.Config, *config.Artefact, error) {
	var cfg config.Config
	missing := true
	for _, c := range configs {
		if path.Clean(target.Dir) == c.Path {
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

func collectDependencies(targets []target.Target, configs []config.Config) ([]Link, error) {
	running := true
	links := []Link{}

	for running {
		newLinks := []Link{}
		newTargets := []target.Target{}

		for _, target := range targets {
			_, artefact, err := getArtefact(target, configs)
			if err != nil {
				return nil, err
			}

			for _, dep := range artefact.Dependencies {
				rebased := dep.Rebase(target)

				link := Link{
					Target: rebased,
					Source: target,
				}

				missing := true

				for _, existing := range links {
					if link.isSame(existing) {
						missing = false
						break
					}
				}

				if missing {
					newLinks = append(newLinks, link)
					newTargets = append(newTargets, rebased)
				}
			}
		}

		if len(newLinks) == 0 {
			running = false
			continue
		}

		links = append(links, newLinks...)
		targets = append(targets, newTargets...)
	}

	output := []Link{}

	for _, link := range links {
		missing := true

		for _, existing := range output {
			if link.isSame(existing) {
				missing = false
				break
			}
		}

		if missing {
			output = append(output, link)
		}
	}

	return output, nil
}

func getDedupedTargets(targets []target.Target, links []Link) []target.Target {
	for _, link := range links {
		targets = append(targets, link.Target)
	}

	output := []target.Target{}

	for _, target := range targets {
		missing := true

		for _, existing := range output {
			if target.IsSame(existing) {
				missing = false
			}
		}

		if missing {
			output = append(output, target)
		}
	}

	return output
}

func mergeMaps(maps ...map[string]string) map[string]string {
	output_map := map[string]string{}

	hasAny := false
	for _, obj := range maps {
		for key, value := range obj {
			output_map[key] = value
			hasAny = true
		}
	}

	if hasAny {
		return output_map
	} else {
		return nil
	}
}

func createNodes(targets []target.Target, configs []config.Config) (*NodeList, error) {
	nodes := NodeList{list: []nodeListElement{}}

	for _, target := range targets {
		cfg, artefact, err := getArtefact(target, configs)
		if err != nil {
			return &nodes, err
		}

		pipeline := &artefact.Pipeline
		for _, step := range pipeline.Steps {
			newNode := Node{
				SharedEnv: mergeMaps(cfg.Env, artefact.Env, pipeline.Env),
				Path:      cfg.Path,
				Artefact:  artefact.Name,
				Step:      step,
				State:     STATE_PENDING,
				DependsOn: []*Node{},
			}

			nodes.AddNode(cfg.Path, artefact.Name, &newNode)
		}
	}

	return &nodes, nil
}

func linkNodes(links []Link, nodes *NodeList) error {
	for _, link := range links {
		sourceNode := nodes.getFirstElement(link.Source.Dir, link.Source.Artefact)
		if sourceNode == nil {
			return fmt.Errorf("failed to generate link for %+v, couldn't find source element", link)
		}

		targetNode := nodes.getLastElement(link.Target.Dir, link.Target.Artefact)
		if targetNode == nil {
			return fmt.Errorf("failed to generate link for %+v, couldn't find target element", link)
		}

		sourceNode.node.DependsOn = append(sourceNode.node.DependsOn, targetNode.node)
	}

	return nil
}

func Solve(targets []target.Target, configs []config.Config) ([]*Node, error) {
	// Recursively compile the chain of dependency links between the input targets and their references
	// and their references references.
	links, err := collectDependencies(targets, configs)
	if err != nil {
		return nil, err
	}

	// Reduce the target and dependency list down to just unique config and artefact combinations
	deduped := getDedupedTargets(targets, links)

	// Create the solver node list for all the unique config and artefact combinations
	nodes, err := createNodes(deduped, configs)
	if err != nil {
		return nil, err
	}

	// Decorate the node list with the dependency links, so that we can figure out the build order
	err = linkNodes(links, nodes)
	if err != nil {
		return nil, err
	}

	return nodes.GetList(), nil
}
