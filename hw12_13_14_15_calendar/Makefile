BIN := "./bin"

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.63.4

lint: install-lint-deps
	golangci-lint --config=.golangci.yml run ./...

generate:
	protoc --go_out=./pb/ --go-grpc_out=./pb/ ./proto/calendar.proto

build:
	mkdir -p $(BIN)
	go build -v -o $(BIN)/calendar ./cmd/calendar
	go build -v -o $(BIN)/scheduler ./cmd/scheduler
	go build -v -o $(BIN)/sender ./cmd/sender

up:
	docker compose -f deployments/docker-compose.yaml up --build -d

down:
	docker compose -f deployments/docker-compose.yaml down

integration-tests:
	@docker compose -f tests/docker-compose.yaml up --build -d postgres rabbitmq sender scheduler calendar
	@docker compose -f tests/docker-compose.yaml up --abort-on-container-exit --exit-code-from integration-tests integration-tests; \
	EXIT_CODE=$$?; \
	docker compose -f tests/docker-compose.yaml down; \
	exit $$EXIT_CODE