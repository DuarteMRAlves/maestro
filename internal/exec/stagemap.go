package exec

import "github.com/DuarteMRAlves/maestro/internal/api"

// StageMap stores the stages to be used in an Execution.
type StageMap struct {
	// rpcs stores the stages that execute grpc methods. These stages are
	// specified by the user and are indexed by their name.
	rpcs map[api.StageName]Stage
	// inputs stores auxiliary stages that are required to manage the input of
	// a rpc stage. These stages are not specified by the user and are indexed
	// by the name of their respective rpc stage.
	inputs map[api.StageName]Stage
	// outputs stores auxiliary stages that are required to manage the output of
	// a rpc stage. These stages are not specified by the user and are indexed
	// by the name of their respective rpc stage.
	outputs map[api.StageName]Stage
}

type StageVisitor func(s Stage)

func NewStageMap() *StageMap {
	return &StageMap{
		rpcs:    map[api.StageName]Stage{},
		inputs:  map[api.StageName]Stage{},
		outputs: map[api.StageName]Stage{},
	}
}

func (m *StageMap) Len() int {
	return len(m.rpcs) + len(m.inputs) + len(m.outputs)
}

func (m *StageMap) AddRpcStage(name api.StageName, stage Stage) {
	m.rpcs[name] = stage
}

func (m *StageMap) AddInputStage(name api.StageName, stage Stage) {
	m.inputs[name] = stage
}

func (m *StageMap) AddOutputStage(name api.StageName, stage Stage) {
	m.outputs[name] = stage
}

func (m *StageMap) Iter(vis StageVisitor) {
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
