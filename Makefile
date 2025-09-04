GITTAG ?= $(shell git describe --tags --always)
GITCOMMIT ?= $(shell git log -1 --pretty=format:"%H")
GOLDFLAGS ?= -s -w -extldflags '-zrelro -znow' -X go.anx.io/go-anxcloud/pkg/client.version=$(GITTAG)
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
	# generate GetIdentifier methods for objects
	tools/tools object-generator --mode runtime --in ./pkg/... --out xxgenerated_object.go
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
	go run github.com/onsi/ginkgo/v2/ginkgo -p  \
	    -timeout 0                              \
	    -race                                   \
	    -coverprofile coverage.out              \
	    --keep-going                            \
	    ./pkg/...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: func-test
func-test:
	go run github.com/onsi/ginkgo/v2/ginkgo -p	\
	    -timeout 180m                          	\
	    -race                                  	\
	    -tags integration                      	\
	    -coverpkg ./...                        	\
	    -coverprofile coverage.out             	\
	    --keep-going                           	\
	    --label-filter="!(old client && slow)" 	\
	    ./pkg/...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: go-lint
go-lint: tools
	@echo "==> Checking source code against linters..."
	@tools/golangci-lint run ./...
	@tools/golangci-lint run --build-tags integration ./...

.PHONY: docs-lint
docs-lint: tools
	@echo "==> Checking docs against linters..."
	@tools/misspell -error -source=text docs/ || (echo; \
		echo "Unexpected misspelling found in docs files."; \
		echo "To automatically fix the misspelling, run 'make docs-lint-fix' and commit the changes."; \
		exit 1)
	@docker run --rm -v $(PWD):/markdown 06kellyjac/markdownlint-cli docs/ || (echo; \
		echo "Unexpected issues found in docs Markdown files."; \
		echo "To apply any automatic fixes, run 'make docs-lint-fix' and commit the changes."; \
		exit 1)

.PHONY: docs-lint-fix
docs-lint-fix: tools
	@echo "==> Applying automatic docs linter fixes..."
	@tools/misspell -w -source=text docs/
	@docker run --rm -v $(PWD):/markdown 06kellyjac/markdownlint-cli --fix docs/

.PHONY: secrets-check
secrets-check: tools
	@echo "==> Scanning for secrets and security issues..."
	@echo "Running gosec security analysis..."
	@tools/golangci-lint run --disable-all --enable gosec ./... || true
	@echo ""
	@echo "Running gitleaks secrets detection..."
	@gitleaks detect --source . --verbose --report-format json --report-path gitleaks-report.json || true
	@if [ -f gitleaks-report.json ] && [ -s gitleaks-report.json ]; then \
		echo "⚠️  Secrets detected! Check gitleaks-report.json for details"; \
		cat gitleaks-report.json | jq -r '.[] | "File: \(.File) Line: \(.StartLine) Secret: \(.Description)"' 2>/dev/null || cat gitleaks-report.json; \
	else \
		echo "✅ No secrets detected by gitleaks"; \
		rm -f gitleaks-report.json; \
	fi

.PHONY: lint
lint: go-lint docs-lint

.PHONY: compat-check
compat-check: tools
	@echo "==> Checking API compatibility..."
	@if git describe --tags --exact-match HEAD >/dev/null 2>&1; then \
		echo "Skipping compatibility check on tagged release"; \
	else \
		LATEST_TAG=$$(git describe --tags --abbrev=0 2>/dev/null || echo ""); \
		if [ -n "$$LATEST_TAG" ]; then \
			echo "Checking compatibility against $$LATEST_TAG"; \
			tools/gorelease -base=$$LATEST_TAG; \
		else \
			echo "No previous release found, skipping compatibility check"; \
		fi \
	fi

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
	cd tools && go build -o . github.com/client9/misspell/cmd/misspell
	cd tools && go build -o . github.com/golangci/golangci-lint/cmd/golangci-lint
	cd tools && go build -o . golang.org/x/exp/cmd/gorelease
	cd tools && go build
	@which gitleaks >/dev/null 2>&1 || (echo "Installing gitleaks..." && go install github.com/zricethezav/gitleaks/v8@latest)
