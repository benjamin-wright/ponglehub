package services

import (
	"context"
	"fmt"
	"time"

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
	if err := r.waitForImage(ctx, repo); err != nil {
		return err
	}

	output, err := r.cmd.Run(
		ctx,
		repo.Path,
		fmt.Sprintf(
			"docker run --rm --name %s-builder -v $(pwd):/home/rust/app:cached -v %s-cargo-git:/root/.cargo/git -v %s-cargo-registry:/root/.cargo/registry rust-builder cargo test --release",
			repo.Name,
			repo.Name,
			repo.Name,
		),
	)

	if err != nil {
		return fmt.Errorf("Error testing package:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}

func (r Rust) waitForImage(ctx context.Context, repo types.Repo) error {
	for {
		_, err := r.cmd.Run(
			ctx,
			repo.Path,
			fmt.Sprintf("docker ps | grep %s-builder", repo.Name),
		)

		if err != nil {
			return nil
		}

		time.Sleep(time.Millisecond * 200)
	}
}

// Build compile the rust binary
func (r Rust) Build(ctx context.Context, repo types.Repo) error {
	if err := r.waitForImage(ctx, repo); err != nil {
		return err
	}

	output, err := r.cmd.Run(
		ctx,
		repo.Path,
		fmt.Sprintf(
			"docker run --rm --name %s-builder -v $(pwd):/home/rust/app:cached -v %s-cargo-git:/root/.cargo/git -v %s-cargo-registry:/root/.cargo/registry rust-builder cargo build --release",
			repo.Name,
			repo.Name,
			repo.Name,
		),
	)

	if err != nil {
		return fmt.Errorf("Error building package:\nError\n%+v\nOutput:\n%s", err, output)
	}

	output, err = r.cmd.Run(ctx, repo.Path, "mkdir -p build && cp target/release/"+repo.Name+" build/")
	if err != nil {
		return fmt.Errorf("Error building package:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}
