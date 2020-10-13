GITTAG ?= $(shell git describe --tags --always)
GITCOMMIT ?= $(shell git log -1 --pretty=format:"%H")
GOLDFLAGS ?= -s -w -extldflags '-zrelro -znow' -X github.com/anexia-it/go-anxcloud.version=$(GITTAG) -X github.com/anexia-it/go-anxcloud.commit=$(GITCOMMIT)
GOFLAGS ?= -trimpath
CGO_ENABLED ?= 0

.PHONY: all
all: build

.PHONY: build
build:
	go build -ldflags "$(GOLDFLAGS)" ./...

.PHONY: generate
generate:
	go generate ./...

.PHONY: benchmark
benchmark:
	go test -bench=. -benchmem ./...

.PHONY: test
test:
	CGO_ENABLED=1 go test -cover -timeout 5m -race ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: download
download:
	go mod download

.PHONY: fmt
fmt:
	gofmt -s -w .

.PHONY: update
update:
	go get -t -u=patch ./...
