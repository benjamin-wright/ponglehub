.PHONY: test run debug rollback

test:
	go test ./...

run:
	go run . --config ../../.geppetto.json build

debug:
	go run . --debug --config ../../.geppetto.json build

rollback:
	go run . --debug --config ../../.geppetto.json rollback