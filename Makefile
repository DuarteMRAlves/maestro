GO ?= go

HOST_OS = $(shell $(GO) env GOOS)
HOST_ARCH = $(shell $(GO) env GOARCH)

# Target OS for compilation. Defaults to host OS.
OS ?= $(HOST_OS)

# Target architecture for compilation. Defaults to host architecture.
ARCH ?= $(HOST_ARCH)

# Timeout for the go tests
TEST_TIMEOUT = 10s

# Directory where to run the tests. Defaults to the internal pkg.
TEST_DIR = ./internal/...

TEST_FLAGS = --timeout $(TEST_TIMEOUT)

PROTOC_FLAGS = -I. --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative

default: help

.PHONY: help
help:
	@echo "Usage:"
	@echo "  make command [options]"
	@echo
	@echo "Commands:"
	@echo "  go/build        builds maestro server binary."
	@echo "  go/test         runs automated tests."
	@echo
	@echo "  docker/build    builds maestro docker image."
	@echo
	@echo "  pb              generates the all .pb.go files for this project."
	@echo "  pb/api          generates the .pb.go files for the grpc api."
	@echo "  pb/test         generates the .pb.go files for the grpc tests."
	@echo "  pb/clean        removes all generated .pb.go files."


.PHONY: docker/build
docker/build:
	docker build -t duartemralves/maestro:latest -f docker/Dockerfile --platform linux/amd64  .

go/build: pb/api pb/test
	GOOS=$(OS) GOARCH=$(ARCH) go build -o target/maestro ./cmd/maestro/maestro.go

go/test: pb/api pb/test
	go test $(TEST_FLAGS) $(TEST_DIR)

.PHONY: pb
pb: pb/api pb/test

.PHONY: pb/api
pb/api:
	cd ./api/pb && protoc $(PROTOC_FLAGS) ./*.proto

.PHONY: pb/test
pb/test:
	cd ./tests/pb/ && protoc $(PROTOC_FLAGS) ./*.proto

pb/clean:
	rm -rf ./api/pb/**.pb.go ./tests/pb/**.pb.go

clean: pb/clean
	rm -rf target

