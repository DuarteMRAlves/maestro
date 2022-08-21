# Maestro

`maestro` is a tool for developing pipelines of grpc services. It connects the services by delivering messages returned from one service as an input to
the next.

![GitHub](https://img.shields.io/github/license/duarteMRAlves/maestro?label=License)

## Getting Started

The `maestro` tool is available as a [docker image](https://hub.docker.com/r/duartemralves/maestro). To pull the image, perform the following command:

```shell
docker pull duartemralves/maestro:v1-latest
```

In order to run `maestro`, you need to specify the pipeline configuration with a .yaml file. The specification for the `maestro` configuration file is detailed [here](docs/CONFIG_FILE.md).

You can then run the pipeline by executing:

```shell
docker run --mount type=bind,source=<config file absolute path>,target=/config.yaml duartemralves/maestro:v1-latest
```

## Developing

* Install golang version 1.19
* Install protobuf version 21.5
* Run the following commands:

```shell
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0
```