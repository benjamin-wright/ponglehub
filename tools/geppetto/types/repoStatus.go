package types

type state string

const (
	pending  state = "pending"
	blocked  state = "blocked"
	building state = "building"
	built    state = "built"
	skipped  state = "skipped"
	errored  state = "errored"
)

// RepoState the build state of a repo
type RepoState struct {
	repo  Repo
	state state
	err   error
	phase string
}

// NewRepoState create a new blank repo state
func NewRepoState(repo Repo) RepoState {
	return RepoState{
		repo:  repo,
		state: pending,
		err:   nil,
		phase: "",
	}
}

// Repo return the underlying repo struct
func (r *RepoState) Repo() Repo {
	return r.repo
}

// Pending returns true if the repo is still waiting to build
func (r *RepoState) Pending() bool {
	return r.state == pending
}

// Building returns true if the repo is building
func (r *RepoState) Building() bool {
	return r.state == building
}

// Phase returns the build phase
func (r *RepoState) Phase() string {
	return r.phase
}

// Built returns true if the repo has been built
func (r *RepoState) Built() bool {
	return r.state == built
}

// Skipped returns true if the repo has been built before and doesn't need updating
func (r *RepoState) Skipped() bool {
	return r.state == skipped
}

// Errored returns an error if the repo failed to build
func (r *RepoState) Errored() error {
	if r.state == errored {
		return r.err
	}

	return nil
}

// Blocked returns true if one of the repo's dependencies failed to build
func (r *RepoState) Blocked() bool {
	return r.state == blocked
}

// Failed returns true if the repo or one of its dependencies failed to build
func (r *RepoState) Failed() bool {
	switch r.state {
	case blocked:
		fallthrough
	case errored:
		return true
	default:
		return false
	}
}

// Success returns true if the repo has successfully built or didn't need updating
func (r *RepoState) Success() bool {
	switch r.state {
	case built:
		fallthrough
	case skipped:
		return true
	default:
		return false
	}
}

// Start set the state to building
func (r *RepoState) Start() {
	r.state = building
}

// Progress advance the build phase
func (r *RepoState) Progress(phase string) {
	r.phase = phase
}

// Block set the state to blocked
func (r *RepoState) Block() {
	r.state = blocked
}

// Complete set the state to built
func (r *RepoState) Complete() {
	r.state = built
}

// Skip set the state to skipped
func (r *RepoState) Skip() {
	r.state = skipped
}

// Error set the state to errored
func (r *RepoState) Error(err error) {
	r.state = errored
	r.err = err
}
