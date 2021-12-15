package stage

import (
	"fmt"
)

// Stage represents a node of the pipeline where a specific rpc method is
// executed.
type Stage struct {
	Name    string
	Asset   string
	Service string
	Method  string
	Address string
}

// Clone creates a copy of the given stage, with the same attributes.
func (s *Stage) Clone() *Stage {
	return &Stage{
		Name:    s.Name,
		Asset:   s.Asset,
		Service: s.Service,
		Method:  s.Method,
		Address: s.Address,
	}
}

// String provides a string representation for the stage.
func (s *Stage) String() string {
	return fmt.Sprintf(
		"Stage{Name:%v,Asset:%v,Service:%v,Method:%v,Address:%v}",
		s.Name,
		s.Asset,
		s.Service,
		s.Method,
		s.Address)
}
