# PongleHub

## Todo:

- bail out of realm setup if everything looks good
- add new folders to Geppetto watch
- Exit Geppetto rollback automatically if no errors

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