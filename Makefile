# The command tree is generated from the SDK's api.json manifest, which
# describes the SDK's Go surface.
#
# By default the manifest is read from the module cache, so a fresh clone can
# regenerate without a checkout of the SDK next door. Set SDK to a local
# checkout to generate against unreleased SDK changes.
SDK ?=
SDK_DIR = $(if $(SDK),$(abspath $(SDK)),$(shell GOWORK=off go list -m -f '{{.Dir}}' github.com/basaltic-sh/sdk-go))

BINARY := basaltic
VERSION ?= dev
LDFLAGS := -X github.com/basaltic-sh/cli/internal/auth.Version=$(VERSION)

.PHONY: all
all: build test

.PHONY: build
build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) .

.PHONY: generate
generate:
	@test -n "$(SDK_DIR)" || { \
		echo "Could not locate the SDK module. Either fetch it:"; \
		echo "    go mod download github.com/basaltic-sh/sdk-go"; \
		echo "or point SDK at a checkout:"; \
		echo "    make generate SDK=/path/to/sdk-go"; \
		exit 1; \
	}
	@test -f "$(SDK_DIR)/api.json" || { \
		echo "No api.json at $(SDK_DIR)."; \
		exit 1; \
	}
	cd internal/gen && GOWORK=off go run . -manifest "$(SDK_DIR)/api.json" -out ../generated
	gofmt -w internal/generated

.PHONY: test
test:
	go test ./...
	cd internal/gen && GOWORK=off go test ./...

.PHONY: vet
vet:
	go vet ./...
	cd internal/gen && GOWORK=off go vet ./...

# Fails when the committed command tree no longer matches the SDK manifest.
.PHONY: check-generated
check-generated: generate
	@git diff --exit-code -- internal/generated || { \
		echo; \
		echo "The command tree is out of date. Commit the regenerated files."; \
		exit 1; \
	}

.PHONY: install
install:
	go install -ldflags "$(LDFLAGS)" .

.PHONY: clean
clean:
	rm -f $(BINARY)
