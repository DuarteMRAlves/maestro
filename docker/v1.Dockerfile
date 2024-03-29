ARG GO=1.19
ARG PROTOC="3.19.4"
ARG WORKSPACE=/opt/maestro

FROM golang:${GO}-bullseye as builder
ARG PROTOC
ARG WORKSPACE

RUN apt-get update && apt-get install -y curl unzip && \
    curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC}/protoc-${PROTOC}-linux-x86_64.zip &&  \
    unzip protoc-${PROTOC}-linux-x86_64.zip -d /opt/protoc && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0

ENV PATH="${PATH}:/opt/protoc/bin"

WORKDIR ${WORKSPACE}

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./ ./

RUN make go/build

FROM gcr.io/distroless/base-debian11
ARG WORKSPACE

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot ${WORKSPACE}/target/maestro /maestro

USER nonroot:nonroot

CMD ["/maestro", "run", "-f", "config.yaml"]