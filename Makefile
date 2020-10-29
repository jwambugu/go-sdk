run:
	cd examples
	go run .

test:
	go test -v -race ./...

lint:
	golangci-lint run --enable=depguard,gci,gochecknoglobals,errorlint,exportloopref