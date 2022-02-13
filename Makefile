GOOS := linux
GOARCH := amd64
DOCKER_PLATFORM := $(GOOS)/$(GOARCH)

build-docker: build
	docker build -t duartemralves/maestro:latest \
		-f docker/Dockerfile \
		--platform $(DOCKER_PLATFORM)  .

build: grpc
	GOOS=linux GOARCH=amd64 go build -o target/maestro ./cmd/maestro/maestro.go

grpc:
	./scripts/genpb.sh

clean:
	rm -rf target api/pb/**.pb.go tests/pb/**.pb.go