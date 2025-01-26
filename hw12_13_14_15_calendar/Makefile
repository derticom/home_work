install-lint-deps:
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.57.2

lint: install-lint-deps
	golangci-lint run ./...
	
test:
	go test -race ./internal/...