package builder

// State represents a build state
type State string

const (
	// NoneState repo build is pending
	NoneState State = "None"
	// BuildingState repo is currently being built
	BuildingState State = "Building"
	// ErroredState repo build failed unexpectedly
	ErroredState State = "Errored"
	// BlockedState repo cannot be built because one of its dependencies failed to build
	BlockedState State = "Blocked"
	// BuiltState repo has been successfully built
	BuiltState State = "Built"
)
