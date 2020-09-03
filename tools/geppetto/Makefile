.PHONY: test run debug rollback

test:
	go test ./...

run:
	go run . --target ../.. build

debug:
	go run . --debug --target ../.. build

rollback:
	go run . --debug --target ../.. rollback