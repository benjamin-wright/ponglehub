package types

// RepoType indicates the type of data in a repo
type RepoType string

const (
	// Node repo is an NPM module
	Node RepoType = "Node"
	// Golang repo is a Go module / application
	Golang RepoType = "Golang"
	// Helm repo is a helm chart
	Helm RepoType = "Helm"
)

// Repo represents a code repo
type Repo struct {
	// Name a unique name for the dependency
	Name string
	// Path the location of the code on disk
	Path string
	// The kind of code in the repo
	RepoType RepoType
	// The paths of other repos one which this one depends
	DependsOn []string
	// Can be built
	Application bool
}
