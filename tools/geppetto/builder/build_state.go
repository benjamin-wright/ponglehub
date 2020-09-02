package builder

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// RepoState the state of a repo
type RepoState struct {
	repo  string
	state State
}

// BuildState represents the current state of all the builds
type BuildState struct {
	Repos []RepoState
}

// NewBuildState creates a new empty build state
func NewBuildState() BuildState {
	return BuildState{Repos: []RepoState{}}
}

func (s *BuildState) add(repo string, state State) {
	s.Repos = append(s.Repos, RepoState{repo: repo, state: state})
}

func (s *BuildState) find(repo string) *RepoState {
	for index, order := range s.Repos {
		if order.repo == repo {
			return &s.Repos[index]
		}
	}

	return nil
}

func (s *BuildState) setEndState(repo string, state State) error {
	existing := s.find(repo)

	if existing == nil {
		s.add(repo, state)
		return nil
	}

	if existing.state != BuildingState {
		return fmt.Errorf("Cannot put repo %s into %s state when already in %s", repo, state, existing.state)
	}

	existing.state = state
	return nil
}

// Count return the number of Repos in the requested state
func (s *BuildState) Count(state State) int {
	counted := 0

	for _, repo := range s.Repos {
		if repo.state == state {
			counted++
		}
	}

	return counted
}

// GetState returns the build state for the given repo
func (s *BuildState) GetState(repo string) State {
	order := s.find(repo)
	if order != nil {
		return order.state
	}

	return NoneState
}

// Build signal that a repo is being built
func (s *BuildState) Build(repo string) error {
	state := s.GetState(repo)
	if state != NoneState {
		return fmt.Errorf("Cannot put repo %s into building state when already in %s", repo, state)
	}

	s.add(repo, BuildingState)
	return nil
}

// Complete signal that a repo build has finished
func (s *BuildState) Complete(repo string) error {
	return s.setEndState(repo, BuiltState)
}

// Block signal that a repo build has been blocked by a dependency build failure
func (s *BuildState) Block(repo string) error {
	return s.setEndState(repo, BlockedState)
}

// Error signal that a repo build has failed
func (s *BuildState) Error(repo string) error {
	return s.setEndState(repo, ErroredState)
}

// CanBuild returns true if the repo itself is not in a building, built, errored or blocked state
func (s *BuildState) CanBuild(repo string, deps []string) (ok bool, block bool) {
	state := s.GetState(repo)

	if state != NoneState {
		logrus.Debugf("Not building %s because state is %s", repo, state)
		return false, false
	}

	for _, dep := range deps {
		if state := s.GetState(dep); state != BuiltState {
			logrus.Debugf("Not building %s because dep %s state is %s", repo, dep, state)
			return false, state == ErroredState || state == BlockedState
		}
	}

	return true, false
}
