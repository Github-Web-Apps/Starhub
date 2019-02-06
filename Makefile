SOURCE_FILES?=./...
TEST_PATTERN?=.
TEST_OPTIONS?=

export PATH := ./bin:$(PATH)
export GO111MODULE := on

# Install all the build and lint dependencies
setup:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh
	curl -sfL https://install.goreleaser.com/github.com/gohugoio/hugo.sh | sh
	go mod download
.PHONY: setup

# Run all the tests
test:
	go test $(TEST_OPTIONS) -failfast -race -coverpkg=./... -covermode=atomic -coverprofile=coverage.txt $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=2m
.PHONY: test

# Run all the tests and opens the coverage report
cover: test
	go tool cover -html=coverage.txt
.PHONY: cover

# gofmt and goimports all go files
fmt:
	gofmt -w -s $(SOURCE_FILES)
	goimports -w $(SOURCE_FILES)
.PHONY: fmt

# Run all the linters
lint:
	./bin/golangci-lint run --tests=false --enable-all ./...
.PHONY: lint

# Run all the tests and code checks
ci: lint test
.PHONY: ci

# Build a beta version of releaser
build:
	go build
.PHONY: build

.DEFAULT_GOAL := build
