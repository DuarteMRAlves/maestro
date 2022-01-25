package orchestration

import (
	"fmt"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
)

// Stage represents a node of the pipeline where a specific rpc method is
// executed.
type Stage struct {
	name    apitypes.StageName
	phase   apitypes.StagePhase
	rpcSpec *RpcSpec
	asset   apitypes.AssetName

	// orchestration is the Orchestration where this stage is inserted.
	orchestration *Orchestration
}

func NewStage(
	name apitypes.StageName,
	rpcSpec *RpcSpec,
	asset apitypes.AssetName,
	orchestration *Orchestration,
) *Stage {
	return &Stage{
		name:          name,
		rpcSpec:       rpcSpec,
		asset:         asset,
		orchestration: orchestration,
		phase:         apitypes.StagePending,
	}
}

func (s *Stage) Name() apitypes.StageName {
	return s.name
}

func (s *Stage) Address() string {
	return s.rpcSpec.address
}

func (s *Stage) IsPending() bool {
	return s.phase == apitypes.StagePending
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

func (s *Stage) ToApi() *apitypes.Stage {
	return &apitypes.Stage{
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
		s.rpcSpec.address)
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
