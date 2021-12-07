FROM golang:1.17.4-alpine

RUN apk add --no-cache bash protoc protobuf-dev gcc libc-dev && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0 \