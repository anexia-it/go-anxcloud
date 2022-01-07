GITTAG ?= $(shell git describe --tags --always)
GITCOMMIT ?= $(shell git log -1 --pretty=format:"%H")
GOLDFLAGS ?= -s -w -extldflags '-zrelro -znow' -X github.com/anexia-it/go-anxcloud/pkg/client.version=$(GITTAG)
GOFLAGS ?= -trimpath
CGO_ENABLED ?= 0

.PHONY: all
all: build

.PHONY: build
build: fmtcheck go-lint
	go build -ldflags "$(GOLDFLAGS)" ./...

.PHONY: generate
generate: tools
	# generate object tests
	tools/tools object-generator --mode tests --in ./pkg/... --out xxgenerated_object_test.go
	# run golang default generator
	go generate ./...

.PHONY: depscheck
depscheck:
	@echo "==> Checking source code dependencies..."
	@go mod tidy
	@git diff --exit-code -- go.mod go.sum || \
		(echo; echo "Found differences in go.mod/go.sum files. Run 'go mod tidy' or revert go.mod/go.sum changes."; exit 1)
	@# reset go.sum to state before checking if it is clean
	@git checkout -q go.sum

.PHONY: benchmark
benchmark:
	go test -bench=. -benchmem ./...

.PHONY: test
test:
	CGO_ENABLED=1 go test -coverpkg ./pkg/... -coverprofile coverage.out -timeout 0 -race ./pkg/...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: func-test
func-test:
	# we run the tests in our tests directory and for every package listed additionally
	# some of the newer tests (starting with the refactored CloudDNS tests) can run as
	# either unit or integration tests, using mocks of the API for unit tests. They need
	# the -integration-test flag to run as integration tests.
	CGO_ENABLED=1 go test -timeout 180m                \
	    -coverpkg ./pkg/... -coverprofile coverage.out \
	    ./tests/...                                    \
	    ./pkg/clouddns/zone/...                        \
	    -args -integration-test
	go tool cover -html=coverage.out -o coverage.html

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
	cd tools && go build

.PHONY: install-precommit-hook
install-precommit-hook: .git/hooks/pre-commit

.PHONY: .git/hooks/pre-commit
.git/hooks/pre-commit: scripts/pre-commit
	cp $< $@
	chmod +x $@
