package stage

import (
	"fmt"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/reflection"
)

// Stage represents a node of the pipeline where a specific rpc method is
// executed.
type Stage struct {
	name    apitypes.StageName
	phase   apitypes.StagePhase
	asset   apitypes.AssetName
	address string

	// Descriptor for the rpc that this stage calls.
	rpc reflection.RPC
}

func New(
	name apitypes.StageName,
	address string,
	asset apitypes.AssetName,
	rpc reflection.RPC,
) *Stage {
	return &Stage{
		name:    name,
		address: address,
		asset:   asset,
		rpc:     rpc,
		phase:   apitypes.StagePending,
	}
}

func (s *Stage) Name() apitypes.StageName {
	return s.name
}

func (s *Stage) Address() string {
	return s.address
}

func (s *Stage) Rpc() reflection.RPC {
	return s.rpc
}

func (s *Stage) IsPending() bool {
	return s.phase == apitypes.StagePending
}

// Clone creates a copy of the given stage, with the same attributes.
func (s *Stage) Clone() *Stage {
	return &Stage{
		name:    s.name,
		asset:   s.asset,
		address: s.address,
		rpc:     s.rpc,

		phase: s.phase,
	}
}

func (s *Stage) ToApi() *apitypes.Stage {
	return &apitypes.Stage{
		Name:    s.name,
		Phase:   s.phase,
		Asset:   s.asset,
		Service: s.rpc.Service().Name(),
		Rpc:     s.rpc.Name(),
		Address: s.address,
	}
}

// String provides a string representation for the stage.
func (s *Stage) String() string {
	return fmt.Sprintf(
		"Stage{Name:%v,Phase%v,Asset:%v,Rpc:%v,Address:%v}",
		s.name,
		s.phase,
		s.asset,
		s.rpc.FullyQualifiedName(),
		s.address)
}
