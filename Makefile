GITTAG ?= $(shell git describe --tags --always)
GITCOMMIT ?= $(shell git log -1 --pretty=format:"%H")
GOLDFLAGS ?= -s -w -extldflags '-zrelro -znow' -X github.com/anexia-it/go-anxcloud.version=$(GITTAG) -X github.com/anexia-it/go-anxcloud.commit=$(GITCOMMIT)
GOFLAGS ?= -trimpath
CGO_ENABLED ?= 0

.PHONY: all
all: build

.PHONY: build
build: fmtcheck go-lint
	go build -ldflags "$(GOLDFLAGS)" ./...

.PHONY: generate
generate:
	go generate ./...

.PHONY: benchmark
benchmark:
	go test -bench=. -benchmem ./...

.PHONY: test
test:
	CGO_ENABLED=1 go test -cover -timeout 0 -race ./...

.PHONY: go-lint
go-lint:
	@echo "==> Checking source code against linters..."
	@golangci-lint run ./...

.PHONY: docs-lint
docs-lint:
	@echo "==> Checking docs against linters..."
	@misspell -error -source=text docs/ || (echo; \
		echo "Unexpected misspelling found in docs files."; \
		echo "To automatically fix the misspelling, run 'make docs-lint-fix' and commit the changes."; \
		exit 1)
	@docker run -v $(PWD):/markdown 06kellyjac/markdownlint-cli docs/ || (echo; \
		echo "Unexpected issues found in docs Markdown files."; \
		echo "To apply any automatic fixes, run 'make docs-lint-fix' and commit the changes."; \
		exit 1)
	@terrafmt diff ./docs --check --pattern '*.md' --quiet || (echo; \
		echo "Unexpected differences in docs HCL formatting."; \
		echo "To see the full differences, run: terrafmt diff ./docs --pattern '*.md'"; \
		echo "To automatically fix the formatting, run 'make docs-lint-fix' and commit the changes."; \
		exit 1)

.PHONY: docs-lint-fix
docs-lint-fix:
	@echo "==> Applying automatic docs linter fixes..."
	@misspell -w -source=text docs/
	@docker run -v $(PWD):/markdown 06kellyjac/markdownlint-cli --fix docs/

.PHONY: lint
lint: go-lint docs-lint

.PHONY: vendor
vendor:
	go mod vendor

.PHONY: fmt
fmt:
	gofmt -s -w .

.PHONY: fmtcheck
fmtcheck:
	@./scripts/gofmtcheck.sh

.PHONY: tools
tools:
	cd tools && go install github.com/client9/misspell/cmd/misspell
	cd tools && go install github.com/golangci/golangci-lint/cmd/golangci-lint
	cd tools && go install github.com/katbyte/terrafmt
