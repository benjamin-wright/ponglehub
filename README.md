# PongleHub

## Todo:

- Make auth-operator client a dumb client, and add a manager to dedupe noisy kube-events and prevent bombarding auth api with requests

## Dependencies:

- brew install coreutils
- brew install k3d
- brew install operator-sdk
- setup your .envrc (copy .envrc-example and fill in the blanks)

## Build:

- `earthly +generate` to create deployment manifests
- `earthly +repos` to start local private npm repo
- `earthly +infra` to launch the local application cluster

## Architecture

### Application

![](docs/pongle-architecture.png)

### Auth flow

![](docs/pongle-auth.png)
