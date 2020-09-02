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
