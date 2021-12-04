# Image to run test workflow with all dependencies installed
FROM golang:1.17.4-bullseye

RUN apt-get update && \
    apt-get install -y protobuf-compiler && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1 && \
    go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0