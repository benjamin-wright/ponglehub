# PongleHub

## Todo:

- all: Add README.md to everything
- auth-controller to refresh periodically and track rollout with CRD status
- geppetto to find int-test repos inside service repos

## Dependencies:

- brew install coreutils
- brew install k3d
- brew install step
- setup your .envrc (copy .envrc-example and fill in the blanks)
- rust (nightly)
  - `brew tap SergioBenitez/osxct && brew install FiloSottile/musl-cross/musl-cross` for cross-compiling on OSX
  - `rustup target add x86_64-unknown-linux-musl --toolchain=nightly` to add musl compile target
- on OSX: add `export OPENSSL_INCLUDE_DIR=$(brew --prefix openssl)/include` and `export OPENSSL_LIB_DIR=$(brew --prefix openssl)/lib` to your environment

## To get up and running:

`make cluster` to run everything. This will automatically launch `Geppetto` to watch over your local build too!

If `Geppetto` falls over for any reason, run `geppetto watch` in the root dir to restart it

> NB: for a quick cluster-free setup, run `make repos` to just spin up the local repos in docker!

To build and deploy the application images, run `tilt up` in the root dir, then press the `space` bar to open the web UI.

## Clean up

`make clean` will tear down the cluster, or your local repos if you went with that option

`geppetto rollback` will roll all the auto-bumped version numbers back to `1.0.0`

## Architecture

### Application

![](docs/pongle-architecture.png)

### Auth flow

![](docs/pongle-auth.png)
