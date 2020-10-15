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

// Test run rust unit tests
func (r Rust) Test(ctx context.Context, repo types.Repo) error {
	output, err := r.cmd.Run(ctx, repo.Path, "cargo test --release")
	if err != nil {
		return fmt.Errorf("Error testing package:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}

// Build compile the rust binary
func (r Rust) Build(ctx context.Context, repo types.Repo) error {
	output, err := r.cmd.Run(ctx, repo.Path, "TARGET_CC=x86_64-linux-musl-gcc RUSTFLAGS=\"-C linker=x86_64-linux-musl-gcc\" cargo build --release --target=x86_64-unknown-linux-musl")
	if err != nil {
		return fmt.Errorf("Error building package:\nError\n%+v\nOutput:\n%s", err, output)
	}

	output, err = r.cmd.Run(ctx, repo.Path, "mkdir -p build && cp target/x86_64-unknown-linux-musl/release/"+repo.Name+" build/")
	if err != nil {
		return fmt.Errorf("Error building package:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}
