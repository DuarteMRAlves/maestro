## Orchestration Modes Proposal

This document describes a proposal for two execution modes for orchestrations: online and offline.

In orchestrations, messages are created in a source stage and suffer a series of transformations in intermediate until they reach the sink stage where they are consumed.
These modes specify the behaviour of the connections between the stages.

**Offline mode** focuses on data processing. 
With this mode, all messages are processed.
If necessary, faster stages block to allow other stages to process their data.


**Online mode** targets real-time applications.
It aims to reduce the time between the creation of a message and the arrival of the respective transformed message to the sink stage.
If necessary, this mode discards older messages in order to keep the processed messages up-to-date.

## Table of contents

This document details each mode, and outlines a possible specification on how to configure the execution modes, dividing it into three sections:

* [Offline Mode](#offline-mode)
* [Online Mode](#online-mode)
* [Configuration](#configuration)

## Offline Mode

**Offline mode** is focused on data processing without any time constraints.
Data is transferred from one stage to another. 
With this mode, eventually all data will be processed and no messages will be lost.

### Implementation Details

In **offline mode**, the stages are connected with buffered queues. 
The messages are processed in order and no message is discarded.
If a stage is slow, the message queues for the inputs of that stage will fill up.
When the queues reach a maximum size, the upstream stage blocks until the messages are processed.


## Online Mode

**Online mode** is ideal for real-time applications, where some loss of data is acceptable for better latency.
This mode aims to reduce the time between the creation of a message in the source stage and the arrival of the respective final message to the sink stage.
With this mode, some older messages may be lost if a stage is too slow.

### Implementation Details

With **online mode**, the stages are also connected with buffered queues, in order to ensure they are processed in order.
Furthermore, an extra routine periodically checks the size of each queue.
If the size of a queue is above a given threshold, it is drained, allowing for more recent messages to be added for processing.

## Configuration

The execution can be configured using an enumeration:

```go
type OrchestrationExecutionMode int

const (
	OfflineExecution OrchestrationExecutionMode = iota
	OnlineExecution OrchestrationExecutionMode
)
```

This allows for support for future modes. 
By default, the offline mode should be used such that the user does not lose that without noticing.

In yaml, the mode can be configured withing the `orchestration` spec as follows:
```yaml
kind: orchestration
spec:
  name: orchestration
  execution_mode: offline / online
```

If the `execution_mode` tag is not specified, the **offline** mode should be used.