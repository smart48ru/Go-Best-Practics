.PHONY: build
build:
	@echo 'Start build'
	@go build -o bin/filescanner cmd/scanner/main.go
	@echo 'The app was successfully built at bin/filescanner'

.PHONY: run
run: build
	@echo 'Start app ./bin/filescanner'
	@bin/filescanner

.PHONY: help
help: build
	@echo 'Start app ./bin/filescanner -h'
	@bin/filescanner -h

.PHONY: test_cover
test_cover:
	@echo 'Start Unit test cover'
	@go test ./... -coverprofile fmt

.PHONY: test
test:
	@echo 'Start Unit test'
	@go test -race -covermode=atomic ./...

.PHONY: integration
integration:
	@echo 'Start Integration test'
	@go test -race -covermode=atomic --tags=integration ./...

.PHONY: lint
lint:
	@echo 'Start lint'
	golangci-lint run
