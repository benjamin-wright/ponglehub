package builder

// State represents a build state
type State string

const (
	// NoneState toy build is pending
	NoneState State = "None"
	// BuildingState toy is currently being built
	BuildingState State = "Building"
	// ErroredState toy build failed unexpectedly
	ErroredState State = "Errored"
	// BlockedState toy cannot be built because one of its dependencies failed to build
	BlockedState State = "Blocked"
	// BuiltState toy has been successfully built
	BuiltState State = "Built"
)
