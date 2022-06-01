package maestro

import (
    "dagger.io/dagger"
    "dagger.io/dagger/core"
    "universe.dagger.io/bash"
    "universe.dagger.io/docker"
    "universe.dagger.io/go"
)

#PlanImage: {
    go: string
    protoc: string
    workdir: string

    docker.#Dockerfile & {
        dockerfile: contents: """
            ARG GO=\(go)
            ARG PROTOC=\(protoc)

            FROM debian:bullseye-slim AS builder
            ARG PROTOC

            RUN apt-get update && apt-get install -y curl unzip && \\
                curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC}/protoc-${PROTOC}-linux-x86_64.zip &&  \\
                unzip protoc-${PROTOC}-linux-x86_64.zip -d /opt/protoc

            FROM --platform=linux/amd64 golang:${GO}-bullseye
            ARG PROTOC

            COPY --from=builder /opt/protoc /opt/protoc
            ENV PATH="${PATH}:/opt/protoc/bin"

            RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1 && \\
                go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0

            COPY . \(workdir)
            """
    }
}

#Protobuf: {
    input: docker.#Image
    dir: string
    protoc: string | *"protoc"
    _input: input
    run: bash.#Run & {
        input: _input
        workdir: dir
        script: contents: "\(protoc) -I. --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative ./*.proto"
    }, 
    contents: core.#Subdir & {
        input: run.output.rootfs
        path:  dir
    }
}

dagger.#Plan & {
    client: filesystem: {
        "./": read: contents: dagger.#FS,
        "./test/protobuf/unit": write: contents: actions.pb.unit.contents.output
        "./test/protobuf/integration": write: contents: actions.pb.integration.contents.output
    }

    actions: {
        params: {
            go: version: "1.18.2"
            protoc: version: "3.19.4"
        }
        deps: #PlanImage & {
            go: params.go.version
            protoc: params.protoc.version
            workdir: "/workspace"
            source: client.filesystem."./".read.contents
        }
        pb: {
            unit: #Protobuf & {
                input: deps.output
                dir: "/workspace/test/protobuf/unit"
            }
            integration: #Protobuf & {
                input: deps.output
                dir: "/workspace/test/protobuf/integration"
            }
        }
        test: {
            _prep_unit: core.#Subdir & {
                input: pb.unit.run.output.rootfs
                path: "/workspace"
            }
            unit: go.#Test & {
                source: _prep_unit.output 
                package: "./internal/..."
            }
            _prep_integration: core.#Subdir & {
                input: pb.integration.run.output.rootfs
                path: "/workspace"
            }
            integration: go.#Test & {
                source: _prep_integration.output 
                package: "./test/integration/..."
            }
        }
    }
}