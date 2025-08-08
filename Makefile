HELM_PLUGIN_NAME := cleanup
LDFLAGS := "-X main.version=${VERSION}"
MOD_PROXY_URL ?= https://goproxy.io

.PHONY: build
build:
	export CGO_ENABLED=0 && \
	go build -o bin/${HELM_PLUGIN_NAME} -ldflags $(LDFLAGS) ./cmd

.PHONY: format
format:
	go fmt ./...

.PHONY: lint
lint: install-build-deps
	golangci-lint run

.PHONY: test
test:
	go test -v ./...

.PHONY: tag
tag:
	@scripts/tag.sh

install-build-deps:
ifeq (, $(shell which mockery))
	go install github.com/vektra/mockery/v3@v3.2.1
endif
ifeq (, $(shell which go-licenses))
	go install github.com/google/go-licenses@v1.6.0
endif
ifeq (, $(shell which copywrite))
	go install github.com/hashicorp/copywrite@v0.22.0
endif
ifeq (, $(shell which golangci-lint))
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.2
endif