package builder

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func stateToString(state State) string {
	switch state {
	case NoneState:
		return "none"
	case BuildingState:
		return "build"
	case ErroredState:
		return "error"
	case BlockedState:
		return "block"
	case BuiltState:
		return "built"
	}

	panic(fmt.Sprintf("Unrecognised state: %s", state))
}

func stringToState(code string) State {
	switch code {
	case "none":
		return NoneState
	case "build":
		return BuildingState
	case "error":
		return ErroredState
	case "block":
		return BlockedState
	case "built":
		return BuiltState
	}

	panic(fmt.Sprintf("Unrecognised state rune: %s", code))
}

func stateFromCode(code string) BuildState {
	orders := []repoState{}

	if code == "" {
		return BuildState{orders: orders}
	}

	for _, c := range strings.Split(code, ",") {
		parts := strings.Split(c, ":")
		orders = append(orders, repoState{
			repo:  parts[0],
			state: stringToState(parts[1]),
		})
	}

	return BuildState{orders: orders}
}

func assertCode(t *testing.T, code string, state BuildState) {
	b := strings.Builder{}

	for i, order := range state.orders {
		if i > 0 {
			b.WriteRune(',')
		}
		b.WriteString(order.repo)
		b.WriteRune(':')
		b.WriteString(stateToString(order.state))
	}

	assert.Equal(t, code, b.String())
}

func TestBuildStateGetState(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		s := stateFromCode("")
		assert.Equal(t, NoneState, s.GetState("a"))
	})

	for _, state := range []State{BuildingState, BlockedState, ErroredState, BuiltState} {
		t.Run(fmt.Sprintf("In state: %s", state), func(t *testing.T) {
			s := BuildState{orders: []repoState{{repo: "repo1", state: state}}}
			assert.Equal(t, state, s.GetState("repo1"))
		})
	}

	for _, state := range []State{BuildingState, BlockedState, ErroredState, BuiltState} {
		t.Run(fmt.Sprintf("Other in state: %s", state), func(t *testing.T) {
			s := BuildState{orders: []repoState{{repo: "2", state: state}}}
			assert.Equal(t, NoneState, s.GetState("1"))
		})
	}
}

func TestBuildStateBuild(t *testing.T) {
	t.Run("first run", func(t *testing.T) {
		state := stateFromCode("")

		assert.Nil(t, state.Build("a"))
		assertCode(t, "a:build", state)
	})

	for _, test := range []struct {
		initial  string
		expected string
	}{
		{initial: "b:build", expected: "b:build,a:build"},
		{initial: "b:error", expected: "b:error,a:build"},
		{initial: "b:block", expected: "b:block,a:build"},
		{initial: "b:built", expected: "b:built,a:build"},
	} {
		t.Run(fmt.Sprintf("Initial state: %s", test.initial), func(t *testing.T) {
			s := stateFromCode(test.initial)

			assert.Nil(t, s.Build("a"))
			assertCode(t, test.expected, s)
		})
	}

	for _, code := range []string{"a:build", "a:block", "a:error", "a:built", "a:build,b:block"} {
		t.Run(fmt.Sprintf("Already in %s state", code), func(t *testing.T) {
			s := stateFromCode(code)

			assert.Error(t, s.Build("a"))
			assertCode(t, code, s)
		})
	}
}

func TestBuildStateComplete(t *testing.T) {
	for _, code := range []string{"", "a:build"} {
		t.Run(fmt.Sprintf("succeeds for state: %s", code), func(t *testing.T) {
			state := stateFromCode(code)

			assert.Nil(t, state.Complete("a"))
			assertCode(t, "a:built", state)
		})
	}

	for _, code := range []string{"a:block", "a:error", "a:built"} {
		t.Run(fmt.Sprintf("Fails for state %s", code), func(t *testing.T) {
			state := stateFromCode(code)

			assert.Error(t, state.Complete("a"))
			assertCode(t, code, state)
		})
	}

	for _, test := range []struct {
		initial  string
		expected string
	}{
		{initial: "b:build", expected: "b:build,a:built"},
		{initial: "b:build,a:build", expected: "b:build,a:built"},
		{initial: "a:build,b:build", expected: "a:built,b:build"},
		{initial: "b:error", expected: "b:error,a:built"},
		{initial: "b:error,a:build", expected: "b:error,a:built"},
		{initial: "a:build,b:error", expected: "a:built,b:error"},
		{initial: "b:block", expected: "b:block,a:built"},
		{initial: "b:block,a:build", expected: "b:block,a:built"},
		{initial: "a:build,b:block", expected: "a:built,b:block"},
		{initial: "b:built", expected: "b:built,a:built"},
		{initial: "b:built,a:build", expected: "b:built,a:built"},
		{initial: "a:build,b:built", expected: "a:built,b:built"},
	} {
		t.Run(fmt.Sprintf("Succeeds when another repo is in %s state", test.initial), func(t *testing.T) {
			state := stateFromCode(test.initial)

			assert.Nil(t, state.Complete("a"))
			assertCode(t, test.expected, state)
		})
	}
}

func TestBuildStateBlock(t *testing.T) {
	for _, code := range []string{"", "a:build"} {
		t.Run(fmt.Sprintf("succeeds for state: %s", code), func(t *testing.T) {
			state := stateFromCode(code)

			assert.Nil(t, state.Block("a"))
			assertCode(t, "a:block", state)
		})
	}

	for _, code := range []string{"a:block", "a:error", "a:built"} {
		t.Run(fmt.Sprintf("Fails for state %s", code), func(t *testing.T) {
			state := stateFromCode(code)

			assert.Error(t, state.Block("a"))
			assertCode(t, code, state)
		})
	}

	for _, test := range []struct {
		initial  string
		expected string
	}{
		{initial: "b:build", expected: "b:build,a:block"},
		{initial: "b:build,a:build", expected: "b:build,a:block"},
		{initial: "a:build,b:build", expected: "a:block,b:build"},
		{initial: "b:error", expected: "b:error,a:block"},
		{initial: "b:error,a:build", expected: "b:error,a:block"},
		{initial: "a:build,b:error", expected: "a:block,b:error"},
		{initial: "b:block", expected: "b:block,a:block"},
		{initial: "b:block,a:build", expected: "b:block,a:block"},
		{initial: "a:build,b:block", expected: "a:block,b:block"},
		{initial: "b:built", expected: "b:built,a:block"},
		{initial: "b:built,a:build", expected: "b:built,a:block"},
		{initial: "a:build,b:built", expected: "a:block,b:built"},
	} {
		t.Run(fmt.Sprintf("Succeeds when another repo is in %s state", test.initial), func(t *testing.T) {
			state := stateFromCode(test.initial)

			assert.Nil(t, state.Block("a"))
			assertCode(t, test.expected, state)
		})
	}
}

func TestBuildStateError(t *testing.T) {
	for _, code := range []string{"", "a:build"} {
		t.Run(fmt.Sprintf("succeeds for state: %s", code), func(t *testing.T) {
			state := stateFromCode(code)

			assert.Nil(t, state.Error("a"))
			assertCode(t, "a:error", state)
		})
	}

	for _, code := range []string{"a:block", "a:error", "a:built"} {
		t.Run(fmt.Sprintf("Fails for state %s", code), func(t *testing.T) {
			state := stateFromCode(code)

			assert.Error(t, state.Error("a"))
			assertCode(t, code, state)
		})
	}

	for _, test := range []struct {
		initial  string
		expected string
	}{
		{initial: "b:build", expected: "b:build,a:error"},
		{initial: "b:build,a:build", expected: "b:build,a:error"},
		{initial: "a:build,b:build", expected: "a:error,b:build"},
		{initial: "b:error", expected: "b:error,a:error"},
		{initial: "b:error,a:build", expected: "b:error,a:error"},
		{initial: "a:build,b:error", expected: "a:error,b:error"},
		{initial: "b:block", expected: "b:block,a:error"},
		{initial: "b:block,a:build", expected: "b:block,a:error"},
		{initial: "a:build,b:block", expected: "a:error,b:block"},
		{initial: "b:built", expected: "b:built,a:error"},
		{initial: "b:built,a:build", expected: "b:built,a:error"},
		{initial: "a:build,b:built", expected: "a:error,b:built"},
	} {
		t.Run(fmt.Sprintf("Succeeds when another repo is in %s state", test.initial), func(t *testing.T) {
			state := stateFromCode(test.initial)

			assert.Nil(t, state.Error("a"))
			assertCode(t, test.expected, state)
		})
	}
}

func TestBuildStateCanBuild(t *testing.T) {
	for _, test := range []struct {
		code         string
		repo         string
		dependencies []string
	}{
		{code: "", repo: "a", dependencies: []string{}},
		{code: "b:built", repo: "a", dependencies: []string{"b"}},
		{code: "b:built,c:built", repo: "a", dependencies: []string{"b", "c"}},
	} {
		t.Run(fmt.Sprintf("True with repo %s and deps %v for state %s", test.repo, test.dependencies, test.code), func(t *testing.T) {
			s := stateFromCode(test.code)
			ok, blocked := s.CanBuild(test.repo, test.dependencies)
			assert.True(t, ok)
			assert.False(t, blocked)
		})
	}

	for _, test := range []struct {
		code         string
		repo         string
		dependencies []string
	}{
		{code: "a:build", repo: "a", dependencies: []string{}},
		{code: "a:error", repo: "a", dependencies: []string{}},
		{code: "a:block", repo: "a", dependencies: []string{}},
		{code: "a:built", repo: "a", dependencies: []string{}},
		{code: "a:build,b:built", repo: "a", dependencies: []string{"b"}},
		{code: "a:error,b:built", repo: "a", dependencies: []string{"b"}},
		{code: "a:block,b:built", repo: "a", dependencies: []string{"b"}},
		{code: "a:built,b:built", repo: "a", dependencies: []string{"b"}},
		{code: "", repo: "a", dependencies: []string{"b"}},
		{code: "b:build", repo: "a", dependencies: []string{"b"}},
		{code: "b:built", repo: "a", dependencies: []string{"b", "c"}},
		{code: "b:built,c:build", repo: "a", dependencies: []string{"b", "c"}},
	} {
		t.Run(fmt.Sprintf("False with repo %s and deps %v for state %s", test.repo, test.dependencies, test.code), func(t *testing.T) {
			s := stateFromCode(test.code)
			ok, blocked := s.CanBuild(test.repo, test.dependencies)
			assert.False(t, ok)
			assert.False(t, blocked)
		})
	}

	for _, test := range []struct {
		code         string
		repo         string
		dependencies []string
	}{
		{code: "b:error", repo: "a", dependencies: []string{"b"}},
		{code: "b:block", repo: "a", dependencies: []string{"b"}},
		{code: "b:built,c:error", repo: "a", dependencies: []string{"b", "c"}},
		{code: "b:built,c:block", repo: "a", dependencies: []string{"b", "c"}},
	} {
		t.Run(fmt.Sprintf("Blocked with repo %s and deps %v for state %s", test.repo, test.dependencies, test.code), func(t *testing.T) {
			s := stateFromCode(test.code)
			ok, blocked := s.CanBuild(test.repo, test.dependencies)
			assert.False(t, ok)
			assert.True(t, blocked)
		})
	}
}
