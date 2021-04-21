run_test:
	go test -v -race ./...
run_lint:
	golangci-lint run --disable-all --enable=golint,depguard,gci,gochecknoglobals,errorlint,exportloopref,typecheck,goimports,misspell,govet,ineffassign,gosimple,deadcode,structcheck