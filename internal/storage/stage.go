package storage

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/api"
)

// Stage represents a node of the pipeline where a specific rpc method is
// executed.
type Stage struct {
	name    api.StageName
	phase   api.StagePhase
	rpcSpec *RpcSpec
	asset   api.AssetName

	// orchestration is the Orchestration where this stage is inserted.
	orchestration *api.Orchestration
}

func NewStage(
	name api.StageName,
	rpcSpec *RpcSpec,
	asset api.AssetName,
	orchestration *api.Orchestration,
) *Stage {
	return &Stage{
		name:          name,
		rpcSpec:       rpcSpec,
		asset:         asset,
		orchestration: orchestration,
		phase:         api.StagePending,
	}
}

func (s *Stage) Name() api.StageName {
	return s.name
}

func (s *Stage) Address() string {
	return s.rpcSpec.address
}

func (s *Stage) IsPending() bool {
	return s.phase == api.StagePending
}

// Clone creates a copy of the given stage, with the same attributes.
func (s *Stage) Clone() *Stage {
	return &Stage{
		name:    s.name,
		asset:   s.asset,
		rpcSpec: s.rpcSpec.Clone(),

		phase: s.phase,

		orchestration: s.orchestration,
	}
}

func (s *Stage) ToApi() *api.Stage {
	return &api.Stage{
		Name:    s.name,
		Phase:   s.phase,
		Asset:   s.asset,
		Service: s.rpcSpec.service,
		Rpc:     s.rpcSpec.rpc,
		Address: s.rpcSpec.address,
	}
}

// String provides a string representation for the stage.
func (s *Stage) String() string {
	return fmt.Sprintf(
		"Stage{Name:%v,Phase%v,Asset:%v,Rpc:%v,Address:%v}",
		s.name,
		s.phase,
		s.asset,
		fmt.Sprintf("%s/%s", s.rpcSpec.service, s.rpcSpec.rpc),
		s.rpcSpec.address,
	)
}

type RpcSpec struct {
	address string
	service string
	rpc     string
}

func NewRpcSpec(address string, service string, rpc string) *RpcSpec {
	return &RpcSpec{
		address: address,
		service: service,
		rpc:     rpc,
	}
}

func (r *RpcSpec) Clone() *RpcSpec {
	return &RpcSpec{
		address: r.address,
		service: r.service,
		rpc:     r.rpc,
	}
}
