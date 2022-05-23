.PHONY: build
build:
	@echo 'Start build'
	@go build -o bin/filescanner cmd/scanner/main.go
	@echo 'The app was successfully built at bin/filescanner'

.PHONY: run
run: build
	@echo 'Start app bin/filescanner'
	@bin/filescanner

.PHONY: help
help: build
	@echo 'Start app bin/filescanner -h'
	@bin/filescanner -h
