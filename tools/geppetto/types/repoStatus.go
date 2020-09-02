package types

// RepoStatus the build state of a repo
type RepoStatus struct {
	Repo     Repo
	Blocked  bool
	Building bool
	Built    bool
	Skipped  bool
	Error    error
	Phase    string
}

// Pending returns true if the repo is still waiting to build
func (r *RepoStatus) Pending() bool {
	return !r.Blocked && !r.Building && !r.Built && !r.Skipped && r.Error == nil
}

// Blocker returns true if the repo has failed to build
func (r *RepoStatus) Blocker() bool {
	return r.Blocked || r.Error != nil
}

// Success returns true if the repo is successfully built
func (r *RepoStatus) Success() bool {
	return r.Built || r.Skipped
}

// SetBuilding set the state to building
func (r *RepoStatus) SetBuilding() {
	r.Building = true
}

// SetBlocked set the state to blocked
func (r *RepoStatus) SetBlocked() {
	r.Blocked = true
}

// SetComplete set the state to built
func (r *RepoStatus) SetComplete() {
	r.Building = false
	r.Built = true
}

// SetSkipped set the state to skipped
func (r *RepoStatus) SetSkipped() {
	r.Building = false
	r.Skipped = true
}

// SetError set the state to errored
func (r *RepoStatus) SetError(err error) {
	r.Building = false
	r.Error = err
}
