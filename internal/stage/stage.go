package stage

import (
	"fmt"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/reflection"
	"github.com/DuarteMRAlves/maestro/internal/stage/input"
	"github.com/DuarteMRAlves/maestro/internal/stage/output"
	"github.com/DuarteMRAlves/maestro/internal/worker"
)

// Stage represents a node of the pipeline where a specific rpc method is
// executed.
type Stage struct {
	Name    string
	Phase   apitypes.StagePhase
	Asset   apitypes.AssetName
	Address string

	// Descriptor for the Rpc that this stage calls.
	Rpc reflection.RPC

	Worker worker.Worker

	// Input and Output describe the connections to other stages
	Input  *input.Input
	Output *output.Output
}

func NewDefault() *Stage {
	return &Stage{
		Phase:  apitypes.StagePending,
		Input:  input.NewDefault(),
		Output: output.NewDefault(),
	}
}

// Clone creates a copy of the given stage, with the same attributes.
func (s *Stage) Clone() *Stage {
	return &Stage{
		Name:    s.Name,
		Asset:   s.Asset,
		Address: s.Address,
		Rpc:     s.Rpc,

		Phase: s.Phase,
	}
}

func (s *Stage) ToApi() *apitypes.Stage {
	return &apitypes.Stage{
		Name:    s.Name,
		Phase:   s.Phase,
		Asset:   s.Asset,
		Service: s.Rpc.Service().Name(),
		Rpc:     s.Rpc.Name(),
		Address: s.Address,
	}
}

// String provides a string representation for the stage.
func (s *Stage) String() string {
	return fmt.Sprintf(
		"Stage{Name:%v,Phase%v,Asset:%v,Rpc:%v,Address:%v}",
		s.Name,
		s.Phase,
		s.Asset,
		s.Rpc.FullyQualifiedName(),
		s.Address)
}
