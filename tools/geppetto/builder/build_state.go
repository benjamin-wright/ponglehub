package builder

import (
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/types"
)

// buildState represents the current state of all the builds
type buildState struct {
	repos []types.RepoStatus
}

// NewbuildState creates a new empty build state
func newBuildState(repos []types.Repo) buildState {
	state := buildState{repos: []types.RepoStatus{}}

	for _, repo := range repos {
		state.repos = append(state.repos, types.RepoStatus{
			Repo: repo,
		})
	}

	return state
}

func (s *buildState) find(repo string) *types.RepoStatus {
	for index, state := range s.repos {
		if state.Repo.Name == repo {
			return &s.repos[index]
		}
	}

	return nil
}

// NumBuilding return the number of repos in the requested state
func (s *buildState) numBuilding() int {
	counted := 0

	for _, repo := range s.repos {
		if repo.Building {
			counted++
		}
	}

	return counted
}

// CanBuild returns true if the repo itself is not in a building, built, errored or blocked state
func (s *buildState) canBuild(repo string) (ok bool, block bool) {
	state := s.find(repo)

	if !state.Pending() {
		logrus.Debugf("Not building %s because build is finished", repo, state)
		return false, false
	}

	for _, dep := range state.Repo.DependsOn {
		depState := s.find(dep)

		if depState.Success() {
			continue
		}

		if depState.Blocker() {
			logrus.Debugf("Not building %s because dep %s state is blocking", repo, dep)
			return false, true
		}

		logrus.Debugf("Not building %s because dep %s state is not built yet", repo, dep)
		return false, false
	}

	return true, false
}
