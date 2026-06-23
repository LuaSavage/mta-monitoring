.PHONY: test vet generate build-example run-example

test:
	go test ./... -race

vet:
	go vet ./...

generate:
	go generate ./...

build-example:
	go build -o example-bin ./cmd/example

run-example:
	go run ./cmd/example
