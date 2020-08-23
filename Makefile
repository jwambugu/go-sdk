run_examples:
	cd examples
	go run .

test:
	go test -v -race ./...