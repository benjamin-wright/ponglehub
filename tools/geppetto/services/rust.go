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

func (r Rust) isRunning(ctx context.Context, repo types.Repo) bool {
	_, err := r.cmd.Run(
		ctx,
		repo.Path,
		fmt.Sprintf("docker ps | grep %s-builder", repo.Name),
	)

	return err == nil
}

func (r Rust) startBuilder(ctx context.Context, repo types.Repo) error {
	output, err := r.cmd.Run(
		ctx,
		repo.Path,
		fmt.Sprintf(
			"docker run --rm -d --name %s-builder --workdir /volume -v $(pwd):/volume:delegated -v %s-cargo-git:/root/.cargo/git -v %s-cargo-registry:/root/.cargo/registry -v %s-cargo-target:/volume/target rust:1.48.0 tail -f /dev/null",
			repo.Name,
			repo.Name,
			repo.Name,
			repo.Name,
		),
	)

	if err != nil {
		return fmt.Errorf("Error starting builder:\nError\n%+v\nOutput:\n%s", err, output)
	}

	output, err = r.cmd.Run(
		ctx,
		repo.Path,
		fmt.Sprintf(
			"docker exec %s-builder /bin/sh -c \"rustup update nightly; rustup default nightly\"",
			repo.Name,
		),
	)

	if err != nil {
		return fmt.Errorf("Error testing package:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}

// Test run rust unit tests
func (r Rust) Test(ctx context.Context, repo types.Repo) error {
	if !r.isRunning(ctx, repo) {
		err := r.startBuilder(ctx, repo)
		if err != nil {
			return err
		}
	}

	output, err := r.cmd.Run(
		ctx,
		repo.Path,
		fmt.Sprintf(
			"docker exec %s-builder cargo test",
			repo.Name,
		),
	)

	if err != nil {
		return fmt.Errorf("Error testing package:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}

// Build compile the rust binary
func (r Rust) Build(ctx context.Context, repo types.Repo) error {
	if !r.isRunning(ctx, repo) {
		err := r.startBuilder(ctx, repo)
		if err != nil {
			return err
		}
	}

	output, err := r.cmd.Run(
		ctx,
		repo.Path,
		fmt.Sprintf(
			"docker exec %s-builder /bin/sh -c \"cargo build && mkdir -p build && cp target/debug/%s build/%s\"",
			repo.Name,
			repo.Name,
			repo.Name,
		),
	)

	if err != nil {
		return fmt.Errorf("Error building package:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}
