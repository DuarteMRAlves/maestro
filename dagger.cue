package maestro

import (
    "dagger.io/dagger"
    "dagger.io/dagger/core"
    "universe.dagger.io/bash"
    "universe.dagger.io/docker"
    "universe.dagger.io/go"
)

#Protobuf: {
    input: docker.#Image
    dir: string
    _input: input
    run: bash.#Run & {
        input: _input
        workdir: dir
        script: contents: "protoc -I. --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative ./*.proto"
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
        deps: docker.#Build & {
            steps: [
                docker.#Pull & {
                    source: "duartemralves/maestro.build-workflow:latest"
                },
                docker.#Copy & {
                    contents: client.filesystem."./".read.contents
                    dest: "/workspace"
                }
            ]
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