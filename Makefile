.PHONY: help

APP=sftp-exporter
APP_EXECUTABLE="./out/$(APP)"
GOLANGCI_LINT_VERSION=v1.53.3
MOCKGEN_VERSION=v0.2.0

ifeq ($(GOLANGCI_LINT),)
	GOLANGCI_LINT=$(shell command -v golangci-lint 2> /dev/null)
endif

help: ## Prints help (only for targets with comment)
	@grep -E '^[a-zA-Z._-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

install-deps: ## install dependencies
	go mod tidy -v

fmt: ## format code
	go fmt

install-golangci-lint: ## install golangci-lint
ifeq ($(GOLANGCI_LINT),)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_LINT_VERSION)
endif

lint: install-golangci-lint ## run lint
	$(GOLANGCI_LINT) run -v

build: ## compile the app
	CGO_ENABLED=0 go build -ldflags "-X main.version=dev" -o $(APP_EXECUTABLE) main.go

test: ## run unit tests
	mkdir -p coverage
	go test -coverprofile coverage/coverage.out ./...
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html

run: ## run the app
	go run main.go

check: install-deps fmt lint mocks test ## runs fmt, lint, test

install-mockgen: ## install mockgen
	go install go.uber.org/mock/mockgen@$(MOCKGEN_VERSION)

mocks: install-mockgen ## generate mocks
	go generate
