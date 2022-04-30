package execute

import (
	"context"

	"github.com/DuarteMRAlves/maestro/internal"
)

type Stage interface {
	Run(context.Context) error
}

// stageMap stores the stages to be used in an Execution.
type stageMap struct {
	// rpcs stores the stages that execute grpc methods. These stages are
	// specified by the user and are indexed by their name.
	rpcs map[internal.StageName]Stage
	// inputs stores auxiliary stages that are required to manage the input of
	// a rpc stage. These stages are not specified by the user and are indexed
	// by the name of their respective rpc stage.
	inputs map[internal.StageName]Stage
	// outputs stores auxiliary stages that are required to manage the output of
	// a rpc stage. These stages are not specified by the user and are indexed
	// by the name of their respective rpc stage.
	outputs map[internal.StageName]Stage
}

type StageVisitor func(s Stage)

func newStageMap() *stageMap {
	return &stageMap{
		rpcs:    map[internal.StageName]Stage{},
		inputs:  map[internal.StageName]Stage{},
		outputs: map[internal.StageName]Stage{},
	}
}

func (m *stageMap) len() int {
	return len(m.rpcs) + len(m.inputs) + len(m.outputs)
}

func (m *stageMap) addRpcStage(name internal.StageName, stage Stage) {
	m.rpcs[name] = stage
}

func (m *stageMap) addInputStage(name internal.StageName, stage Stage) {
	m.inputs[name] = stage
}

func (m *stageMap) addOutputStage(name internal.StageName, stage Stage) {
	m.outputs[name] = stage
}

func (m *stageMap) iter(vis StageVisitor) {
	for _, s := range m.rpcs {
		vis(s)
	}
	for _, s := range m.inputs {
		vis(s)
	}
	for _, s := range m.outputs {
		vis(s)
	}
}
