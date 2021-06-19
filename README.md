# PongleHub

## Todo:

- Make auth-operator client a dumb client, and add a manager to dedupe noisy kube-events and prevent bombarding auth api with requests

## Dependencies:

- brew install coreutils
- brew install k3d
- brew install operator-sdk
- setup your .envrc (copy .envrc-example and fill in the blanks)

## Architecture

### Application

![](docs/pongle-architecture.png)

### Auth flow

![](docs/pongle-auth.png)
