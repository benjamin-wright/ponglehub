package builder

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// repoState the state of a repo
type repoState struct {
	repo  string
	state state
	phase string
}

// buildState represents the current state of all the builds
type buildState struct {
	repos []repoState
}

// NewbuildState creates a new empty build state
func newBuildState() buildState {
	return buildState{repos: []repoState{}}
}

func (s *buildState) add(repo string, state state) {
	s.repos = append(s.repos, repoState{repo: repo, state: state})
}

func (s *buildState) find(repo string) *repoState {
	for index, order := range s.repos {
		if order.repo == repo {
			return &s.repos[index]
		}
	}

	return nil
}

func (s *buildState) setEndState(repo string, state state) error {
	existing := s.find(repo)

	if existing == nil {
		s.add(repo, state)
		return nil
	}

	if existing.state != buildingState {
		return fmt.Errorf("Cannot put repo %s into %s state when already in %s", repo, state, existing.state)
	}

	existing.state = state
	return nil
}

// Count return the number of repos in the requested state
func (s *buildState) Count(state state) int {
	counted := 0

	for _, repo := range s.repos {
		if repo.state == state {
			counted++
		}
	}

	return counted
}

// GetState returns the build state for the given repo
func (s *buildState) GetState(repo string) state {
	order := s.find(repo)
	if order != nil {
		return order.state
	}

	return noneState
}

func (s *buildState) GetPhase(repo string) string {
	order := s.find(repo)
	if order != nil {
		return order.phase
	}

	return ""
}

func (s *buildState) Progress(repo string, phase string) error {
	order := s.find(repo)
	if order.state != buildingState {
		return fmt.Errorf("Cannot set phase %s on repo %s, state is %s", phase, repo, order.state)
	}

	order.phase = phase
	return nil
}

// Build signal that a repo is being built
func (s *buildState) Build(repo string) error {
	state := s.GetState(repo)
	if state != noneState {
		return fmt.Errorf("Cannot put repo %s into building state when already in %s", repo, state)
	}

	s.add(repo, buildingState)
	return nil
}

// Complete signal that a repo build has finished
func (s *buildState) Complete(repo string) error {
	return s.setEndState(repo, builtState)
}

// Skip signal that a repo build has been built before
func (s *buildState) Skip(repo string) error {
	return s.setEndState(repo, skippedState)
}

// Block signal that a repo build has been blocked by a dependency build failure
func (s *buildState) Block(repo string) error {
	return s.setEndState(repo, blockedState)
}

// Error signal that a repo build has failed
func (s *buildState) Error(repo string) error {
	return s.setEndState(repo, erroredState)
}

// CanBuild returns true if the repo itself is not in a building, built, errored or blocked state
func (s *buildState) CanBuild(repo string, deps []string) (ok bool, block bool) {
	state := s.GetState(repo)

	if state != noneState {
		logrus.Debugf("Not building %s because state is %s", repo, state)
		return false, false
	}

	for _, dep := range deps {
		if state := s.GetState(dep); state != builtState && state != skippedState {
			logrus.Debugf("Not building %s because dep %s state is %s", repo, dep, state)
			return false, state == erroredState || state == blockedState
		}
	}

	return true, false
}
