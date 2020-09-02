package builder

// state represents a build state
type state string

const (
	// NoneState repo build is pending
	noneState state = "None"
	// BuildingState repo is currently being built
	buildingState state = "Building"
	// ErroredState repo build failed unexpectedly
	erroredState state = "Errored"
	// BlockedState repo cannot be built because one of its dependencies failed to build
	blockedState state = "Blocked"
	// BuiltState repo has been successfully built
	builtState state = "Built"
	// skippedState repo has been built before
	skippedState state = "Skipped"
)
