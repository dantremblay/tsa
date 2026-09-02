MODULE     := github.com/kassisol/tsa
BINARIES   := tsa tsad

# Version info
COMMIT     := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DIRTY      := $(shell git status --porcelain --untracked-files=no 2>/dev/null)
GIT_TAG    := $(shell git tag -l --contains HEAD 2>/dev/null | head -n 1)

ifdef DIRTY
  GITSTATE := dirty
  VERSION  ?= $(COMMIT)-dirty
else ifdef GIT_TAG
  GITSTATE := clean
  VERSION  ?= $(GIT_TAG)
else
  GITSTATE := clean
  VERSION  ?= $(COMMIT)
endif

LDFLAGS := -s -w \
           -X $(MODULE)/version.Version=$(VERSION) \
           -X $(MODULE)/version.GitCommit=$(COMMIT) \
           -X $(MODULE)/version.GitState=$(GITSTATE) \
           -X $(MODULE)/version.BuildDate=$(shell date +%s)

# Default target
.DEFAULT_GOAL := build

## build: Build both binaries (tsa, tsad)
build:
	@echo "Building $(VERSION) ($(COMMIT))"
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o bin/tsa ./cmd/tsa
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o bin/tsad ./cmd/tsad

## cross: Cross-compile both binaries for multiple platforms
cross:
	@echo "Cross-compiling $(VERSION)"
	@for bin in $(BINARIES); do \
		CGO_ENABLED=0 GOOS=linux  GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/$$bin-linux-amd64 ./cmd/$$bin; \
		CGO_ENABLED=0 GOOS=linux  GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o bin/$$bin-linux-arm64 ./cmd/$$bin; \
		CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/$$bin-darwin-amd64 ./cmd/$$bin; \
		CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o bin/$$bin-darwin-arm64 ./cmd/$$bin; \
	done

## image: Build the tsad daemon Docker image
image: build
	docker build --build-arg version=$(VERSION) -t kassisol/tsa:$(VERSION) -f Dockerfile .

## test: Run tests
test:
	go test ./...

## vet: Run go vet
vet:
	go vet ./cmd/... ./api/... ./cli/... ./daemon/... ./pkg/... ./client/... ./version/...

## tidy: Tidy and vendor dependencies
tidy:
	go mod tidy
	go mod vendor

## generate: Generate man pages and shell completions for both binaries
generate:
	go run gen/man/genman-tsa.go
	go run gen/man/genman-tsad.go
	go run gen/shellcompletion/genshellcompletion-tsa.go
	go run gen/shellcompletion/genshellcompletion-tsad.go

check-nfpm:
	@command -v nfpm >/dev/null 2>&1 || { \
		echo "Error: nfpm not found on PATH."; \
		echo "Install it with:"; \
		echo "    go install github.com/goreleaser/nfpm/v2/cmd/nfpm@latest"; \
		echo "and make sure \"\$$(go env GOPATH)/bin\" is on your PATH."; \
		exit 1; \
	}

## package-rpm: Build RPM packages for tsa and tsa-daemon (requires nfpm)
package-rpm: RELEASE ?= 1
package-rpm: check-nfpm
	@mkdir -p dist
	VERSION=$(VERSION) RELEASE=$(RELEASE) GOARCH=$(or $(TARGETARCH),amd64) \
		nfpm pkg -f nfpm-tsa.yaml --packager rpm --target dist/
	VERSION=$(VERSION) RELEASE=$(RELEASE) GOARCH=$(or $(TARGETARCH),amd64) \
		nfpm pkg -f nfpm-tsad.yaml --packager rpm --target dist/

## package-deb: Build DEB packages for tsa and tsa-daemon (requires nfpm)
package-deb: RELEASE ?= 1
package-deb: check-nfpm
	@mkdir -p dist
	VERSION=$(VERSION) RELEASE=$(RELEASE) GOARCH=$(or $(TARGETARCH),amd64) \
		nfpm pkg -f nfpm-tsa.yaml --packager deb --target dist/
	VERSION=$(VERSION) RELEASE=$(RELEASE) GOARCH=$(or $(TARGETARCH),amd64) \
		nfpm pkg -f nfpm-tsad.yaml --packager deb --target dist/

## packages: Build all packages
packages: generate package-rpm package-deb

## clean: Remove build artifacts
clean:
	rm -rf bin/ dist/

## version: Print version info
version:
	@echo "Version:  $(VERSION)"
	@echo "Commit:   $(COMMIT)"
	@echo "GitState: $(GITSTATE)"

## help: Show this help
help:
	@echo "Usage: make [target]"
	@echo ""
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':'

.PHONY: build cross image test vet tidy generate check-nfpm package-rpm package-deb packages clean version help
