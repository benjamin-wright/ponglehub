# PongleHub

## Todo:

- update keycloak-init config load from environment to use tags on the field names in the config struct
- factor out config load from environment into a separate golang library
- update geppetto to find golang dependencies
- update geppetto to close on Ctrl+C (sigint?)
- bail out of realm setup if everything looks good

## Dependencies:

- brew install coreutils
- brew install k3d
- brew install step
- setup your .envrc (copy .envrc-example and fill in the blanks)

## To get up and running:

`make cluster` to run everything. This will automatically launch `Geppetto` to watch over your local build too!

If `Geppetto` falls over for any reason, run `geppetto watch` in the root dir to restart it

> NB: for a quick cluster-free setup, run `make repos` to just spin up the local repos in docker!

## Clean up

`make clean` will tear down the cluster, or your local repos if you went with that option

`geppetto rollback` will roll all the auto-bumped version numbers back to `1.0.0`