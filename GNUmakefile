# Makefile

# Variables
GO_BIN_DIR := $(PWD)/bin
PROVIDER_DIR := $(PWD)/internal
EXAMPLES_DIR := $(PWD)/examples/provider-install-verification
GOBIN := $(GO_BIN_DIR)

# Phony targets
.PHONY: all build install update-deps run-tests generate-docs debug-client tfplugindocs check-deps clean

# Default target
all: update-deps build install generate-docs

# Check for required dependencies
check-deps:
	@command -v go >/dev/null 2>&1 || { echo >&2 "Go is not installed. Please install Go >=1.21."; exit 1; }
	@command -v terraform >/dev/null 2>&1 || { echo >&2 "Terraform is not installed. Please install Terraform >=1.5+"; exit 1; }

# Build the provider
build: check-deps
	cd $(PROVIDER_DIR) && go build ./...

# Install the provider
install: check-deps
	cd $(PROVIDER_DIR) && go install ./...

# Update Go module dependencies
update-deps:
	cd $(PROVIDER_DIR) && go mod tidy
	$(MAKE) install

# Install tfplugindocs tool
tfplugindocs:
	GOBIN=$(GOBIN) go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

# Generate documentation
generate-docs: tfplugindocs
	$(GOBIN)/tfplugindocs generate --rendered-provider-name Tilgangsportalen

# Check required environment variables
check-env:
	@if [ -z "$$ACC_TEST_SYSTEM_ROLE_OWNER" ]; then \
		echo "Error: ACC_TEST_SYSTEM_ROLE_OWNER is not set."; exit 1; \
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
run-tests: check-env
	TF_ACC=1 go test ./... $(TESTARGS) -cover -coverprofile=c.out -timeout 120m
	go tool cover -html=c.out -o coverage.html

# Clean build artifacts
clean:
	cd $(PROVIDER_DIR) && go clean
