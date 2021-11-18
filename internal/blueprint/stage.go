package blueprint

import (
	"fmt"
)

type Stage struct {
	Name    string
	Asset   string
	Service string
	Method  string
}

func (s *Stage) Clone() *Stage {
	return &Stage{
		Name:    s.Name,
		Asset:   s.Asset,
		Service: s.Service,
		Method:  s.Method,
	}
}

func (s *Stage) String() string {
	return fmt.Sprintf(
		"Stage{Name:'%v',Asset:%v,Service:'%v',Method:'%v'}",
		s.Name,
		s.Asset,
		s.Service,
		s.Method)
}
