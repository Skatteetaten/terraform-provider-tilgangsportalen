# Makefile

# Variables
PROVIDER_DIR := $(PWD)/internal
EXAMPLES_DIR := $(PWD)/examples/provider-install-verification
# GO_MIN_VERSION is the minimum Go version required by the module (the go directive in go.mod). It must not be
# newer than the Go bundled with MegaLinter, which runs with GOTOOLCHAIN=local and lags behind the CI image.
# GO_TOOLCHAIN_VERSION is the Go version used to build (the toolchain directive in go.mod). It should match the
# Go version in the CI image (IMAGE_GO in .gitlab-ci.yml).
GO_MIN_VERSION := 1.26.3
GO_TOOLCHAIN_VERSION := 1.27.2

# Phony targets
.PHONY: all build install update-deps upgrade-deps run-tests generate-docs check-deps clean

# Default target
all: update-deps build install generate-docs go-vet

# Check for required dependencies
check-deps:
	@command -v go >/dev/null 2>&1 || { echo >&2 "Go is not installed. Please install Go >=1.21."; exit 1; }
	@command -v terraform >/dev/null 2>&1 || { echo >&2 "Terraform is not installed. Please install Terraform >=1.5+"; exit 1; }

# Build the provider
build: check-deps
	go build ./...

# Install the provider
install: check-deps
	go install ./...

# Update Go module dependencies
update-deps:
	go mod tidy

# Upgrade all Go module dependencies (including tool dependencies in go.mod) to the latest versions
# compatible with GO_MIN_VERSION. The toolchain directive is removed during the upgrade and set to
# GO_TOOLCHAIN_VERSION afterwards. Check the latest Go version in the CI image (IMAGE_GO in .gitlab-ci.yml) manually in:
# https://gitlab.skead.no/azure-tools/container-registry/go/-/blob/main/CHANGELOG.md?ref_type=heads
upgrade-deps:
	go get -u ./... go@$(GO_MIN_VERSION) toolchain@none
	go get -u tool go@$(GO_MIN_VERSION) toolchain@none
	go mod tidy
	go mod edit -go=$(GO_MIN_VERSION) -toolchain=go$(GO_TOOLCHAIN_VERSION)

# Generate documentation
generate-docs:
	go tool tfplugindocs generate --rendered-provider-name Tilgangsportalen

# Vet examines Go source code and reports suspicious constructs
go-vet:
	go vet ./...

# Check required environment variables
check-env:
	@if [ -z "$$ACC_TEST_SYSTEM_ROLE_OWNER" ]; then \
		echo "Error: ACC_TEST_SYSTEM_ROLE_OWNER is not set."; exit 1; \
	fi
	@if [ -z "$$ACC_TEST_SYSTEM_ROLE_SECURITY_OWNER" ]; then \
		echo "Error: ACC_TEST_SYSTEM_ROLE_SECURITY_OWNER is not set."; exit 1; \
	fi
	@if [ -z "$$TILGANGSPORTALEN_USERNAME" ]; then \
		echo "Error: TILGANGSPORTALEN_USERNAME is not set."; exit 1; \
	fi
	@if [ -z "$$TILGANGSPORTALEN_PASSWORD" ]; then \
		echo "Error: TILGANGSPORTALEN_PASSWORD is not set."; exit 1; \
	fi
	@if [ -z "$$TILGANGSPORTALEN_URL" ]; then \
		echo "Error: TILGANGSPORTALEN_URL is not set."; exit 1; \
	fi

# Run tests
run-tests: check-env go-vet
	TF_ACC=1 go test ./... $(TESTARGS) -cover -coverprofile=c.out -timeout 120m
	go tool cover -html=c.out -o coverage.html

# Run a specific test
run-test: check-env
	@if [ -z "$(test)" ]; then \
		echo "Test name is required. Usage: make run-test test=<test-name>"; \
		exit 1; \
	fi
	TF_ACC=1 go test ./... -count=1 -run='^$(test)$$' -v

# Clean build artifacts
clean:
	cd $(PROVIDER_DIR) && go clean
