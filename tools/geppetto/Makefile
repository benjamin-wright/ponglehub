.PHONY: test watch rollback

test:
	go test ./...

watch:
	go run . --target ../.. watch

rollback:
	go run . --target ../.. rollback

install:
	go build -o ../../bin/geppetto