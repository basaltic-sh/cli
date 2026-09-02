# The command tree is generated from the SDK's api.json manifest, which
# describes the SDK's Go surface. Point SDK at a local checkout.
SDK ?= ../sdks/sdk-go

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
	@test -f "$(SDK)/api.json" || { \
		echo "No api.json at $(SDK). Point SDK at an sdk-go checkout:"; \
		echo "    make generate SDK=/path/to/sdk-go"; \
		exit 1; \
	}
	cd internal/gen && GOWORK=off go run . -manifest "$(abspath $(SDK))/api.json" -out ../generated
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
