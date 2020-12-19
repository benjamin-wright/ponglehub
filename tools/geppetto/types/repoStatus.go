package types

import "context"

type state string

const (
	pending   state = "pending"
	reinstall state = "reinstall"
	blocked   state = "blocked"
	building  state = "building"
	built     state = "built"
	skipped   state = "skipped"
	errored   state = "errored"
)

// RepoState the build state of a repo
type RepoState struct {
	repo    Repo
	state   state
	err     error
	phase   string
	context context.Context
	cancel  context.CancelFunc
}

// NewRepoState create a new blank repo state
func NewRepoState(repo Repo) RepoState {
	return RepoState{
		repo:  repo,
		state: built,
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
	return r.state == pending || r.state == reinstall
}

// Building returns true if the repo is building
func (r *RepoState) Building() bool {
	return r.state == building
}

// Context return the current build context, or nil if there isn't one
func (r *RepoState) Context() context.Context {
	return r.context
}

// Cancel call the cancel function
func (r *RepoState) Cancel() {
	if r.cancel != nil {
		r.cancel()
	}
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

func (r *RepoState) DependsOn(repo string) bool {
	for _, r := range r.Repo().DependsOn {
		if r == repo {
			return true
		}
	}

	return false
}

// Invalidate the previous build
func (r *RepoState) Invalidate() {
	r.state = pending
	r.phase = ""
	r.err = nil
}

// Reinstall invalidates the previous build and request a dependency reinstall
func (r *RepoState) Reinstall() {
	r.state = reinstall
	r.phase = ""
	r.err = nil
}

// Start set the state to building and return true if a reinstall is needed
func (r *RepoState) Start(ctx context.Context, cancel context.CancelFunc) bool {
	shouldReinstall := r.state == reinstall
	r.state = building
	r.context = ctx
	r.cancel = cancel

	return shouldReinstall
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
	r.context = nil
	r.cancel = nil
}

// Skip set the state to skipped
func (r *RepoState) Skip() {
	r.state = skipped
	r.context = nil
	r.cancel = nil
}

// Error set the state to errored
func (r *RepoState) Error(err error) {
	r.state = errored
	r.err = err
	r.context = nil
	r.cancel = nil
}
