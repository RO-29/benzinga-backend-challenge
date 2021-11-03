.DEFAULT_GOAL=build

VERSION?=$(shell (git describe --tags --exact-match 2> /dev/null || git rev-parse HEAD) | sed "s/^v//")
.PHONY: version
version:
	@echo $(VERSION)

PROJECT?=$(shell basename $(shell pwd))

GO_BUILD_DIR=build
.PHONY: build
build:
	mkdir -p $(GO_BUILD_DIR)
	go build -v -ldflags="-s -w -X main.version=$(VERSION)" -o $(GO_BUILD_DIR) ./cmd/...

.PHONY: test
test:
	go test -v -cover -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out -o=coverage.txt
	cat coverage.txt
	go tool cover -html=coverage.out -o=coverage.html


.PHONY: lint
lint:: install-golangci-lint
	$(GOLANGCI_LINT_RUN) --fix

GOLANGCI_LINT_VERSION=v1.41.1
GOLANGCI_LINT_DIR=$(shell go env GOPATH)/pkg/golangci-lint/$(GOLANGCI_LINT_VERSION)
GOLANGCI_LINT_BIN=$(GOLANGCI_LINT_DIR)/golangci-lint
$(GOLANGCI_LINT_BIN):
	curl -vfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOLANGCI_LINT_DIR) $(GOLANGCI_LINT_VERSION)

.PHONY: install-golangci-lint
install-golangci-lint: $(GOLANGCI_LINT_BIN)

GOLANGCI_LINT_RUN=$(GOLANGCI_LINT_BIN) -v run

.PHONY: ensure-command-docker
ensure-command-docker:
	$(call ENSURE_COMMAND,docker,See https://docs.docker.com/install/)

export DOCKER_BUILDKIT=1
docker-build: ensure-command-docker
	docker image build\
	 --build-arg VERSION=$(VERSION)\
	 --tag=$(PROJECT):$(VERSION)\
	 --tag=$(PROJECT):latest\
	 .
