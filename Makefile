.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: lint
## lint: Runs golanglint-ci
lint:
	@golangci-lint run

.PHONY: run
## run: Runs the action
run:
	@go run ./...
