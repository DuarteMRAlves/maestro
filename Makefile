GO ?= go

HOST_OS = $(shell $(GO) env GOOS)
HOST_ARCH = $(shell $(GO) env GOARCH)

# Target OS for compilation. Defaults to host OS.
OS ?= $(HOST_OS)

# Target architecture for compilation. Defaults to host architecture.
ARCH ?= $(HOST_ARCH)

# Timeout for the go tests
UNIT_TEST_TIMEOUT = 20s

# Directory where to run the tests. Defaults to the internal pkg.
UNIT_TEST_DIR = ./internal/...

UNIT_TEST_FLAGS = --timeout $(UNIT_TEST_TIMEOUT) --shuffle on

PROTOC_FLAGS = -I. --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative

default: help

.PHONY: help
help:
	@echo "Usage:"
	@echo "  make command [options]"
	@echo
	@echo "Commands:"
	@echo "  go/build        builds all project binaries."
	@echo "  go/test         runs automated tests."
	@echo "  go/test/unit    runs automated unit tests."
	@echo "  go/test/e2e     runs automated e2e tests."
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

go/build: pb/api
	GOOS=$(OS) GOARCH=$(ARCH) go build -o target/maestro ./cmd/maestro/maestro.go
	GOOS=$(OS) GOARCH=$(ARCH) go build -o target/maestroctl ./cmd/maestroctl/maestroctl.go

go/test: pb/api pb/test go/test/unit go/test/e2e

.PHONY: go/test/unit
go/test/unit: pb/api pb/test
	go test $(UNIT_TEST_FLAGS) $(UNIT_TEST_DIR)

.PHONE: go/go/test/e2e
go/test/e2e: pb/api pb/test
	go test ./tests/e2e

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

