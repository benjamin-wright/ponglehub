package builder

import (
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/types"
)

type buildState struct {
	repos []types.RepoState
}

func newBuildState(repos []types.Repo) buildState {
	state := buildState{repos: []types.RepoState{}}

	for _, repo := range repos {
		state.repos = append(state.repos, types.NewRepoState(repo))
	}

	return state
}

func (s *buildState) find(repo string) *types.RepoState {
	for index, state := range s.repos {
		if state.Repo().Name == repo {
			return &s.repos[index]
		}
	}

	return nil
}

func (s *buildState) invalidate(repo string, path string, reinstall bool) {
	logrus.Debugf("Invalidating %s", repo)
	r := s.find(repo)
	repoObj := r.Repo()

	buildBypass := false

	if r.Building() {
		if repoObj.Application {
			buildBypass = !reinstall
		} else {
			targets := repoObj.BuildTargets()
			for _, target := range targets {
				if target == path {
					buildBypass = true
					break
				}
			}
		}
	}

	if buildBypass {
		logrus.Debugf("Bypassing in-build invalidation of %s (%s) for file %s", repo, repoObj.RepoType, path)
		return
	}

	if reinstall {
		r.Cancel()
		r.Reinstall()
	} else if !repoObj.Application {
		r.Cancel()
		r.Invalidate()
	}
}

func (s *buildState) numBuilding() int {
	counted := 0

	for _, repo := range s.repos {
		if repo.Building() {
			counted++
		}
	}

	return counted
}

func (s *buildState) canBuild(repo string) (ok bool, block bool) {
	state := s.find(repo)

	if !state.Pending() {
		logrus.Debugf("Not building %s because build is finished", repo)
		return false, false
	}

	for _, dep := range state.Repo().DependsOn {
		depState := s.find(dep)

		if depState.Success() {
			continue
		}

		if depState.Failed() {
			logrus.Debugf("Not building %s because dep %s state is blocking", repo, dep)
			return false, true
		}

		logrus.Debugf("Not building %s because dep %s state is not built yet", repo, dep)
		return false, false
	}

	return true, false
}
