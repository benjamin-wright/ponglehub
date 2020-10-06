package types

// RepoUpdate an update message
type RepoUpdate struct {
	Name    string
	Path    string
	Install bool
}
