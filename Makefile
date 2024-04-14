MAIN_FILE=cmd/main.go

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development
.PHONY: lint
lint: assert_docker_installed ## Run code linters.
	docker run -v .:/src -w "/src" golangci/golangci-lint:latest-alpine golangci-lint run ./... --concurrency 2 -v -c golangci.yaml

.PHONY: mod
mod:
	docker run -v .:/src -w "/src" golang:alpine go get -u -t ./...
	docker run -v .:/src -w "/src" golang:alpine go mod tidy
	docker run -v .:/src -w "/src" golang:alpine go mod vendor


##@ Assertions
.PHONY: assert_go_installed
assert_go_installed: ## Assert go is installed.
	@if ! command -v go &> /dev/null; then \
		echo "go is not installed; you need to install it in order to run this command"; \
		exit 1; \
	fi

.PHONY: assert_docker_installed
assert_docker_installed: ## Assert docker is installed.
	@if ! command -v docker &> /dev/null; then \
		echo "docker is not installed; you need to install it in order to run this command"; \
		exit 1; \
	fi
