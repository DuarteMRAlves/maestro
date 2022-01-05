package stage

import (
	"fmt"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/reflection"
)

// Stage represents a node of the pipeline where a specific rpc method is
// executed.
type Stage struct {
	Name    string
	Asset   string
	Address string
	// Descriptor for the Rpc that this stage calls.
	Rpc reflection.RPC

	phase apitypes.StagePhase
}

func NewDefault() *Stage {
	return &Stage{phase: apitypes.StagePending}
}

// Clone creates a copy of the given stage, with the same attributes.
func (s *Stage) Clone() *Stage {
	return &Stage{
		Name:    s.Name,
		Asset:   s.Asset,
		Address: s.Address,
		Rpc:     s.Rpc,

		phase: s.phase,
	}
}

func (s *Stage) ToApi() *apitypes.Stage {
	return &apitypes.Stage{
		Name:    s.Name,
		Phase:   s.phase,
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
		s.phase,
		s.Asset,
		s.Rpc.FullyQualifiedName(),
		s.Address)
}
