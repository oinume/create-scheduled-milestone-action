APP = create-scheduled-milestone
BASE_DIR = github.com/oinume/create-scheduled-milestone
VENDOR_DIR = vendor
GO_GET ?= go get
GO_TEST ?= go test -v -race
GO_TEST_PACKAGES = $(shell go list ./... | grep -v vendor)
GOPATH = $(shell go env GOPATH)
LINT_PACKAGES = $(shell go list ./...)
IMAGE_TAG ?= latest


all: build

.PHONY: git-config
git-config:
	echo "" > ~/.gitconfig
	git config --global url."https://github.com".insteadOf git://github.com
	git config --global http.https://gopkg.in.followRedirects true

build:
	CGO_ENABLED=0 GO111MODULE=on go build -ldflags="-w -s" -o $(APP) .

clean:
	${RM} $(foreach command,$(COMMANDS),bin/$(command))

.PHONY: test
test:
	$(GO_TEST) $(GO_TEST_PACKAGES)

lint: ## Run golangci-lint
	docker run --rm -v ${GOPATH}/pkg/mod:/go/pkg/mod -v $(shell pwd):/app -v $(shell go env GOCACHE):/cache/go -e GOCACHE=/cache/go -e GOLANGCI_LINT_CACHE=/cache/go -w /app golangci/golangci-lint:v1.59.1 golangci-lint run --modules-download-mode=readonly /app/...
.PHONY: lint

lint/fix: ## Run golangci-lint with --fix
	docker run --rm -v ${GOPATH}/pkg/mod:/go/pkg/mod -v $(shell pwd):/app -v $(shell go env GOCACHE):/cache/go -e GOCACHE=/cache/go -e GOLANGCI_LINT_CACHE=/cache/go -w /app golangci/golangci-lint:v1.59.1 golangci-lint run --fix --modules-download-mode=readonly /app/...
.PHONY: lint/fix

# TODO: tag
.PHONY: docker/build
docker/build:
	docker build --pull --no-cache -f Dockerfile .
