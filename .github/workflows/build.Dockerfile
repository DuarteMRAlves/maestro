ARG GO=1.17.6
ARG PROTOC="3.19.3"

FROM debian:bullseye-slim AS builder
ARG PROTOC

RUN apt-get update && apt-get install -y curl unzip && \
    curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC}/protoc-${PROTOC}-linux-x86_64.zip &&  \
    unzip protoc-${PROTOC}-linux-x86_64.zip -d /opt/protoc

FROM golang:${GO}-bullseye
ARG PROTOC

COPY --from=builder /opt/protoc /opt/protoc
ENV PATH="${PATH}:/opt/protoc/bin"

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0