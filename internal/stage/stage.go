package stage

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/api/types"
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
}

// Clone creates a copy of the given stage, with the same attributes.
func (s *Stage) Clone() *Stage {
	return &Stage{
		Name:    s.Name,
		Asset:   s.Asset,
		Address: s.Address,
		Rpc:     s.Rpc,
	}
}

func (s *Stage) ToApi() *types.Stage {
	return &types.Stage{
		Name:    s.Name,
		Asset:   s.Asset,
		Service: s.Rpc.Service().Name(),
		Rpc:     s.Rpc.Name(),
		Address: s.Address,
	}
}

// String provides a string representation for the stage.
func (s *Stage) String() string {
	return fmt.Sprintf(
		"Stage{Name:%v,Asset:%v,Rpc:%v,Address:%v}",
		s.Name,
		s.Asset,
		s.Rpc.FullyQualifiedName(),
		s.Address)
}
