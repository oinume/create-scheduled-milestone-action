APP = create-milestone
BASE_DIR = github.com/oinume/amamonitor
VENDOR_DIR = vendor
GO_GET ?= go get
GO_TEST ?= go test -v -race
GO_TEST_PACKAGES = $(shell go list ./... | grep -v vendor)
GOPATH = $(shell go env GOPATH)
LINT_PACKAGES = $(shell go list ./...)
IMAGE_TAG ?= latest
VERSION_HASH_VALUE = $(shell git rev-parse HEAD | cut -c-7)


all: build

.PHONY: setup
setup: install-linter

.PHONY: install-linter
install-linter:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(GOPATH)/bin v1.18.0

.PHONY: git-config
git-config:
	echo "" > ~/.gitconfig
	git config --global url."https://github.com".insteadOf git://github.com
	git config --global http.https://gopkg.in.followRedirects true

build:
	CGO_ENABLED=0 GO111MODULE=on go build -ldflags="-w -s" -o bin/$(APP) .

clean:
	${RM} $(foreach command,$(COMMANDS),bin/$(command))

.PHONY: test
test:
	$(GO_TEST) $(GO_TEST_PACKAGES)

.PHONY: goimports
goimports:
	goimports -w .

.PHONY: lint
lint: install-linter
	golangci-lint run

# TODO: tag
.PHONY: docker/build
docker/build:
	docker build --pull --no-cache -f Dockerfile .

#.PHONY: gcloud/builds
#gcloud/builds: $(foreach command,$(COMMANDS),gcloud/builds/$(command))
#
#.PHONY: gcloud/builds/%
#gcloud/builds/%:
#	gcloud builds submit . \
#	--project $(GCP_PROJECT_ID) \
#	--config=gcloud-builds.yml \
#	--substitutions=_IMAGE_TAG=$(IMAGE_TAG),_COMMAND=$*
