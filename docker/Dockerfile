ARG GO=1.17.6
ARG PROTOC="3.19.4"
ARG WORKSPACE=/opt/maestro

FROM golang:${GO}-bullseye as builder
ARG PROTOC
ARG WORKSPACE

RUN apt-get update && apt-get install -y curl unzip && \
    curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC}/protoc-${PROTOC}-linux-x86_64.zip &&  \
    unzip protoc-${PROTOC}-linux-x86_64.zip -d /opt/protoc && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0

ENV PATH="${PATH}:/opt/protoc/bin"

WORKDIR ${WORKSPACE}

COPY ./ ./

RUN make go/build

FROM golang:${GO}-bullseye
ARG WORKSPACE

ARG USER=maestro
ARG GROUP=maestro-group

# Create non-privileged user to run
RUN addgroup --system ${GROUP} && \
    adduser --system --no-create-home --ingroup ${GROUP} ${USER} && \
    mkdir ${WORKSPACE} && \
    chown -R ${USER}:${GROUP} ${WORKSPACE}

WORKDIR ${WORKSPACE}

COPY --from=builder --chown=${USER}:${GROUP} ${WORKSPACE}/target/maestro maestro

RUN chmod +x maestro

EXPOSE 50051

CMD ./maestro