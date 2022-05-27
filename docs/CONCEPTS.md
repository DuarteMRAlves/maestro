## Maestro Concepts

Maestro is a tool for connecting grpc services into pipelines. Therefore, the main concept behind `maestro` is a `Pipeline`.

A `Pipeline` is a graph like structure where we have Stages and Links.

A `Stage` is a node of the `Pipeline` graph that processed a message. It specifies a concrete grpc method
to be executed, and has the following fields:

* `Name` that uniquely identifies the stage.
* `Address` where the stage server is running.
* `Service` to specify the grpc service within the sever. Required if multiple services exist, otherwise can be omitted.
* `Method` to specify the grpc method within the service. Required if multiple methods exist, otherwise can be omitted.

A `Link` specifies a connection between two stages. A Link has:

* `Name` to uniquely identify the link.

* `SourceStage` to identify the stage that is the source of the link. Messages returned by the rpc executed in this stage are transferred through this link to the next stage.

* `SourceField` to specify the field of the message returned by `SourceStage` that should be sent through this link. If not specified, the entire message from SourceStage is used.

* `TargetStage` to identify the stage that is the target of the link. Messages that are transferred through this link are used as input for the rpc method in this stage.

* `TargetField` to specify the field of the input message for `TargetStage` that should be set with the messages transferred with this link. If not specified, the entire message is sent as input to `TargetStage`.
