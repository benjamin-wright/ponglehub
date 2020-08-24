package builder

import (
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/config"
)

type planner struct {
	built   []string
	running []string
}

func (p *planner) reset() {
	p.built = []string{}
}

func (p *planner) run(task string) {
	p.running = append(p.running, task)
}

func (p *planner) complete(task string) {
	p.built = append(p.built, task)

	if ok, index := p.isRunning(task); ok {
		end := len(p.running) - 1
		p.running[index] = p.running[end]
		p.running[end] = ""
		p.running = p.running[0:end]
	}
}

func (p *planner) isBuilt(task string) bool {
	for _, built := range p.built {
		if task == built {
			return true
		}
	}

	return false
}

func (p *planner) areBuilt(tasks []string) bool {
	for _, task := range tasks {
		if !p.isBuilt(task) {
			return false
		}
	}

	return true
}

func (p *planner) isRunning(task string) (bool, int) {
	index := -1
	for i, t := range p.running {
		if task == t {
			index = i
			break
		}
	}

	return index != -1, index
}

func (p *planner) canBuild(repo config.Repo) bool {
	if p.isBuilt(repo.Name) {
		logrus.Debugf("Already built: %s", repo.Name)
		return false
	}

	for _, dep := range repo.DependsOn {
		if !p.isBuilt(dep) {
			logrus.Debugf("Dependency not built for %s: %s", repo.Name, dep)
			return false
		}
	}

	return true
}
