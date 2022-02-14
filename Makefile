GO ?= go

HOST_OS = $(shell $(GO) env GOOS)
HOST_ARCH = $(shell $(GO) env GOARCH)

# Target OS for compilation. Defaults to host OS.
OS ?= $(HOST_OS)

# Target architecture for compilation. Defaults to host architecture.
ARCH ?= $(HOST_ARCH)

# Architecture for the generated docker image
DOCKER_ARCH ?= $(ARCH)

# Timeout for the go tests
TEST_TIMEOUT = 10s

# Directory where to run the tests. Defaults to the internal pkg.
TEST_DIR = ./internal/...

TEST_FLAGS = --timeout $(TEST_TIMEOUT)

default: help

.PHONY: help
help:
	@echo "Usage:"
	@echo "  make command [options]"
	@echo
	@echo "Commands:"
	@echo "  go/build        builds maestro server binary."
	@echo "  go/test         runs automated tests."
	@echo "  docker/build    builds maestro docker image."


docker/build: go/build
	docker build -t duartemralves/maestro:latest \
		-f docker/Dockerfile \
		--platform $(DOCKER_ARCH)  .

go/build: grpc
	GOOS=$(OS) GOARCH=$(ARCH) go build -o target/maestro ./cmd/maestro/maestro.go

go/test: grpc
	@echo
	go test $(TEST_FLAGS) $(TEST_DIR)

grpc:
	@echo
	@./scripts/genpb.sh

clean:
	rm -rf target api/pb/**.pb.go tests/pb/**.pb.go

