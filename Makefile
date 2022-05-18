.PHONY: build
build:
	@echo 'Start build'
	@go build -o filescanner cmd/scanner/main.go
	@echo 'The app was successfully built at ./filescanner '

.PHONY: run
run: build
	@echo 'Start build'
	@go build -o filescanner cmd/scanner/main.go
	@echo 'Start app ./filescanner'
	@./filescanner