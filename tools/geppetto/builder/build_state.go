package builder

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type toyOrder struct {
	toy   string
	state State
}

// BuildState represents the current state of all the builds
type BuildState struct {
	orders []toyOrder
}

// NewBuildState creates a new empty build state
func NewBuildState() BuildState {
	return BuildState{orders: []toyOrder{}}
}

func (s *BuildState) add(toy string, state State) {
	s.orders = append(s.orders, toyOrder{toy: toy, state: state})
}

func (s *BuildState) find(toy string) *toyOrder {
	for index, order := range s.orders {
		if order.toy == toy {
			return &s.orders[index]
		}
	}

	return nil
}

func (s *BuildState) setEndState(toy string, state State) error {
	existing := s.find(toy)

	if existing == nil {
		s.add(toy, state)
		return nil
	}

	if existing.state != BuildingState {
		return fmt.Errorf("Cannot put toy %s into %s state when already in %s", toy, state, existing.state)
	}

	existing.state = state
	return nil
}

// GetState returns the build state for the given toy
func (s *BuildState) GetState(toy string) State {
	order := s.find(toy)
	if order != nil {
		return order.state
	}

	return NoneState
}

// Build signal that a toy is being built
func (s *BuildState) Build(toy string) error {
	state := s.GetState(toy)
	if state != NoneState {
		return fmt.Errorf("Cannot put toy %s into building state when already in %s", toy, state)
	}

	s.add(toy, BuildingState)
	return nil
}

// Complete signal that a toy build has finished
func (s *BuildState) Complete(toy string) error {
	return s.setEndState(toy, BuiltState)
}

// Block signal that a toy build has been blocked by a dependency build failure
func (s *BuildState) Block(toy string) error {
	return s.setEndState(toy, BlockedState)
}

// Error signal that a toy build has failed
func (s *BuildState) Error(toy string) error {
	return s.setEndState(toy, ErroredState)
}

// CanBuild returns true if the toy itself is not in a building, built, errored or blocked state
func (s *BuildState) CanBuild(toy string, deps []string) (ok bool, block bool) {
	state := s.GetState(toy)

	if state != NoneState {
		logrus.Debugf("Not building %s because state is %s", toy, state)
		return false, false
	}

	for _, dep := range deps {
		if state := s.GetState(dep); state != BuiltState {
			logrus.Debugf("Not building %s because dep %s state is %s", toy, dep, state)
			return false, state == ErroredState || state == BlockedState
		}
	}

	return true, false
}
