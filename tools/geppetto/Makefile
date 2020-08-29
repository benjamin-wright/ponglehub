.PHONY: test run debug rollback

test:
	go test ./...

run:
	go run . --config ../../.geppetto.json --target ../.. build

debug:
	go run . --debug --config ../../.geppetto.json --target ../.. build

rollback:
	go run . --debug --config ../../.geppetto.json --target ../.. rollback