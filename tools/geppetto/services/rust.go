package services

import (
	"context"
	"fmt"

	"ponglehub.co.uk/geppetto/types"
)

type rustIo interface {
	ReadCargoFile(path string) (CargoFile, error)
	FileExists(path string) bool
}

// Rust collects methods related to golang repos
type Rust struct {
	io  rustIo
	cmd commander
}

// NewRustService create a new Rust repo instance, or error if the path doesn't contain a rust project
func NewRustService() Rust {
	return Rust{io: &IO{}, cmd: &Commander{}}
}

// GetRepo returns a repo object representing the node project at the designated file path
func (r Rust) GetRepo(path string) (types.Repo, error) {
	empty := types.Repo{}

	tomlFile := path + "/Cargo.toml"

	data, err := r.io.ReadCargoFile(tomlFile)
	if err != nil {
		return empty, err
	}

	name := data.PackageInfo.Name

	return types.Repo{
		Name:      name,
		Path:      path,
		RepoType:  types.Rust,
		DependsOn: []string{},
	}, nil
}

// Install get cargo deps
func (r Rust) Install(ctx context.Context, repo types.Repo) error {
	output, err := r.cmd.Run(ctx, repo.Path, "docker run --rm -v $(pwd):/home/rust/src:cached rustcc cargo update")
	if err != nil {
		return fmt.Errorf("Error getting cargo deps:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}

// Check run rust compile checks tests
func (r Rust) Check(ctx context.Context, repo types.Repo) error {
	output, err := r.cmd.Run(ctx, repo.Path, "docker run --rm -v $(pwd):/home/rust/src:cached rustcc cargo check")
	if err != nil {
		return fmt.Errorf("Error checking package:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}

// Test run rust unit tests
func (r Rust) Test(ctx context.Context, repo types.Repo) error {
	output, err := r.cmd.Run(ctx, repo.Path, "docker run --rm -v $(pwd):/home/rust/src:cached rustcc cargo test")
	if err != nil {
		return fmt.Errorf("Error testing package:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}

// Build compile the rust binary
func (r Rust) Build(ctx context.Context, repo types.Repo) error {
	output, err := r.cmd.Run(ctx, repo.Path, "docker run --rm -v $(pwd):/home/rust/src:cached rustcc cargo build --release")
	if err != nil {
		return fmt.Errorf("Error building package:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}
