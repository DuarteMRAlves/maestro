# Maestro

## Overview

Maestro is a tool for orchestrating grpc services into pipelines.
It connects the services by delivering the returned messages from one service
as an input to the next.

## Main Concepts
There are
three main concepts inside maestro - `Asset`, `Blueprint`, `Orchestration`:

### Asset

An `Asset` is a component that can be added to a pipeline.
It may have an associated docker image that exposes the grpc api.

An Asset has the following properties:

* `Id` that uniquely identifies the Asset.
* `Name` that is a human-readable string to identify the Asset.
* `Image` (optional) which is the name of the image associated with this Asset.

### Blueprint

A `Blueprint` defines the architecture of the pipeline.
A Blueprint is a graph like structure where we have Stages and Links.

A `Stage` is an instantiation of an Asset. 
It specifies a concrete grpc method to be executed, and has the following fields:

* `Id` Uniquely identifies the stage inside a Blueprint.
* `AssetId` Identifies the associated Asset.
* `Address` Defines where the grpc service is running.
* `Service` (optional) Specifies the grpc service within the sever (required if multiple services exist, otherwise can be omitted)
* `Method` (optional) Specifies the grpc method within the service (required if multiple methods exist, otherwise can be omitted)

A `Link` specifies a connection between two stages. A Link has:

* `SourceStageId` Id of the Stage that is the source of the connection.
* `SourceFieldName` (optional) If not specified the whole message is transferred. 
  It defines the field name of the source output that will be transferred through the connection.
* `TargetStageId` Id of the Stage that is the target of the connection.
* `TargetFieldName`(optional) If not defined, the whole received message is delivered to the stage. 
  Otherwise, the default message for the target stage is created, and the field with the given variable name is set with the received message.

### Orchestration

An `Orchestration` is an instantiation of a Blueprint where the pipeline is executed. 
 