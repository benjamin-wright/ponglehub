package builder

import "ponglehub.co.uk/geppetto/config"

type planner struct {
	built []string
	cfg   *config.Config
}

func (p planner) complete(task string) {
	p.built = append(p.built, task)
}

func (p planner) isBuilt(task string) bool {
	for _, built := range p.built {
		if task == built {
			return true
		}
	}

	return false
}

func (p planner) canBuild(repo config.Repo) bool {
	if p.isBuilt(repo.Name) {
		return false
	}

	for _, dep := range repo.DependsOn {
		if !p.isBuilt(dep) {
			return false
		}
	}

	return false
}
