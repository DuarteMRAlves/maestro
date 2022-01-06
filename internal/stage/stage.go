package stage

import (
	"fmt"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/link"
	"github.com/DuarteMRAlves/maestro/internal/reflection"
)

// Stage represents a node of the pipeline where a specific rpc method is
// executed.
type Stage struct {
	Name    string
	Phase   apitypes.StagePhase
	Asset   string
	Address string

	// Descriptor for the Rpc that this stage calls.
	Rpc reflection.RPC

	// Input and Output describe the connections to other stages
	Input  Input
	Output Output
}

func NewDefault() *Stage {
	return &Stage{
		Phase: apitypes.StagePending,
		Input: Input{
			connections: map[string]link.Link{},
		},
		Output: Output{
			connections: map[string]link.Link{},
		},
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

// Input represents the several input fields
type Input struct {
	connections map[string]link.Link
}

// Output represents the several input fields
type Output struct {
	connections map[string]link.Link
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
