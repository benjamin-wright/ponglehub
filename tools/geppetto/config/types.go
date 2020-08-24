package config

// RepoType indicates the type of data in a repo
type RepoType string

const (
	// Node repo is an NPM module
	Node RepoType = "Node"
	// Go repo is a GOLANG module
	Go RepoType = "Go"
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
}

// Config represents the app configuration
type Config struct {
	Repos []Repo
}

// FileStruct struct for umarshalling config data
type FileStruct struct {
	Node []RepoStruct `json:"node"`
	Go   []RepoStruct `json:"go"`
}

// RepoStruct struct for unmarshalling config data
type RepoStruct struct {
	Name         string   `json:"name"`
	Path         string   `json:"path"`
	Dependencies []string `json:"dependencies"`
}
