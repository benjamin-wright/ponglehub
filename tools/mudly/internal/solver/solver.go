package solver

import (
	"fmt"
	"path"

	"ponglehub.co.uk/tools/mudly/internal/config"
	"ponglehub.co.uk/tools/mudly/internal/runner"
	"ponglehub.co.uk/tools/mudly/internal/steps"
	"ponglehub.co.uk/tools/mudly/internal/target"
	"ponglehub.co.uk/tools/mudly/internal/utils"
)

type link struct {
	Target target.Target
	Source target.Target
}

func (l link) isSame(m link) bool {
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

func getPipeline(cfg *config.Config, artefact *config.Artefact) (*config.Pipeline, error) {
	if artefact.Steps != nil && len(artefact.Steps) > 0 {
		return &config.Pipeline{
			Name:  "",
			Steps: artefact.Steps,
		}, nil
	} else if artefact.Pipeline != "" {
		for _, pipeline := range cfg.Pipelines {
			if pipeline.Name == artefact.Pipeline {
				return &pipeline, nil
			}
		}
	}

	return nil, fmt.Errorf("failed to get pipeline from artefact %s (%s)", artefact.Name, cfg.Path)
}

func collectDependencies(targets []target.Target, configs []config.Config) ([]link, error) {
	running := true
	links := []link{}

	for running {
		newLinks := []link{}
		newTargets := []target.Target{}

		for _, target := range targets {
			_, artefact, err := getArtefact(target, configs)
			if err != nil {
				return nil, err
			}

			for _, dep := range artefact.DependsOn {
				rebased := dep.Rebase(target)

				link := link{
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

	output := []link{}

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

func getDedupedTargets(targets []target.Target, links []link) []target.Target {
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

func createRunnable(step config.Step) (runner.Runnable, error) {
	if step.Command != "" {
		return steps.CommandStep{
			Name:      step.Name,
			Condition: step.Condition,
			Watch:     step.Watch,
			Command:   step.Command,
			Env:       step.Env,
		}, nil
	}

	if step.Dockerfile != "" {
		return steps.DockerStep{
			Name:       step.Name,
			Dockerfile: step.Dockerfile,
			Context:    step.Context,
			Tag:        step.Tag,
		}, nil
	}

	return nil, fmt.Errorf("failed to convert config step into runnable step: %+v", step)
}

func createNodes(targets []target.Target, configs []config.Config) (*NodeList, error) {
	nodes := NodeList{list: []nodeListElement{}}

	for _, target := range targets {
		cfg, artefact, err := getArtefact(target, configs)
		if err != nil {
			return &nodes, err
		}

		pipeline, err := getPipeline(cfg, artefact)
		if err != nil {
			return &nodes, err
		}
		for _, step := range pipeline.Steps {
			runnable, err := createRunnable(step)
			if err != nil {
				return &nodes, err
			}

			newNode := runner.Node{
				SharedEnv: utils.MergeMaps(cfg.Env, artefact.Env, pipeline.Env),
				Path:      cfg.Path,
				Artefact:  artefact.Name,
				Step:      runnable,
				State:     runner.STATE_PENDING,
				DependsOn: []*runner.Node{},
			}

			nodes.AddNode(cfg.Path, artefact.Name, &newNode)
		}
	}

	return &nodes, nil
}

func linkNodes(links []link, nodes *NodeList) error {
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

func Solve(targets []target.Target, configs []config.Config) ([]*runner.Node, error) {
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
