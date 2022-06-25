## Configuration File Specification

This document details the specification for the `maestro` configuration files. The confirguration files are in yaml format.

## Version 1

This configuration file version details a single resource per yaml document. Multiple resources can be specified in the same file by using multiple yaml documents, delimited by `---`.

Each resource has a `kind` and a `spec` field, both mandatory.

`kind` is specifies a type for the resource to be created. Can be one of:

* pipeline
* stage
* link

`spec` specifies the configuration for the resource. The following sections describe the fields in the section for each resource kind.

### Pipeline Configuration

Here is an example of a Pipeline configuration:

```yaml
kind: pipeline
spec:
    name: hello-world-pipeline
    execution_mode: Offline
```

A Pipeline spec accepts the following fields:

`name` uniquely identifies the resource. (Required)

`execution_mode` specifies the execution mode for the pipeline. Can be either `Offline` or `Online`. (Optional, Default: Offline)

### Stage Configuration

Here is an example of a Stage configuration:

```yaml
kind: stage
spec:
    name: hello-world-stage
    address: localhost:12345
    service: GreetingService
    method: Greet
    pipeline: hello-world-pipeline
```

A Stage spec has the following fields:

`name` uniquely identifies the resource. (Required)

`address` specifies the address the `maestro` should use to connect to the grpc server. (Required)

`service` specifies the name of the grpc service to call. May be ommited if the grpc server only has one service, in which case, that service will be chosen. (Optional)

`method` specifies the name of the grpc method to call. May be ommited if the selected grpc service only has one method, in which case, that method is chosen. (Optional)

`pipeline` is the name of the pipeline that this stage is included in. (Required) 

### Link Configuration

Here is an example of a Link configuration:

```yaml
kind: link
spec:
    name: hello-world-link
    source_stage: hello-world-source
    target_stage: hello-world-target
    pipeline: hello-world-pipeline
```

A Link spec accepts the following fields:

`name` uniquely identifies the resource. (Required)

`source_stage` is the name of the stage that is the source of the link. Messages returned by the rpc executed in this stage are transferred through this link to the next stage. (Required)

`source_field` specifies the field of the message returned by `source_stage` that should be sent through this link. If not specified, the entire message from SourceStage is used. (Optional)

`target_stage` is the name of the stage that is the target of the link. Messages that are transferred through this link are used as input for the rpc method in this stage. (Required)

`target_field` specifies the field of the input message for `target_stage` that should be set with the messages transferred with this link. If not specified, the entire message is sent as input to `target_stage`. (Optional)

`num_empty_messages` specifies the number of empty messages to fill this link with when the pipeline is starting. It allows for cycles, by providing a mechanism to send a first empty message for one of the stages. (Optional, Requires Offline execution mode).