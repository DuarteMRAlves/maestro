package maestro

import (
    "dagger.io/dagger"
    "dagger.io/dagger/core"
    "universe.dagger.io/bash"
    "universe.dagger.io/docker"
)

#Protobuf: {
    input: docker.#Image
    dir: string
    _input: input
    run: bash.#Run & {
        input: _input
        workdir: dir
        script: contents: #"""
        protoc -I. --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative ./*.proto
        """#
    }, 
    contents: core.#Subdir & {
        input: run.output.rootfs
        path:  dir
    }
}

dagger.#Plan & {
    client: filesystem: {
        "./": read: contents: dagger.#FS,
        "./test/protobuf/unit": write: contents: actions.pb.test.contents.output
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
            test: #Protobuf & {
                input: deps.output
                dir: "/workspace/test/protobuf/unit"
            }
            integration: #Protobuf & {
                input: deps.output
                dir: "/workspace/test/protobuf/integration"
            }
        }
        test: bash.#Run & {
            input: pb.test.run.output
            workdir: "/workspace"
            script: contents: #"""
            go test --timeout 20s --shuffle on ./internal/...
            """#
        }
    }
}