# PongleHub

## Todo:

- Update auth-user status after setting password
- Properly check token in cookie for normal events
- Set up event-gateway to hold requests by ID and wait for response events

## Dependencies:

- brew install k3d
- setup your .envrc (copy .envrc-example and fill in the blanks)

## Build:

## Architecture

### Application

![](docs/pongle-architecture.png)

### Game move events

![](docs/pongle-game-move.png)
