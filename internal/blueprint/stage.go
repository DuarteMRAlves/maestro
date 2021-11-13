package blueprint

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/identifier"
)

type Stage struct {
	Id      identifier.Id
	Name    string
	AssetId identifier.Id
	Service string
	Method  string
}

func (s *Stage) Clone() *Stage {
	return &Stage{
		Id:      s.Id.Clone(),
		Name:    s.Name,
		AssetId: s.AssetId.Clone(),
		Service: s.Service,
		Method:  s.Method,
	}
}

func (s *Stage) String() string {
	return fmt.Sprintf(
		"Stage{Id:%v,Name:'%v',AssetId:%v,Service:'%v',Method:'%v'}",
		s.Id,
		s.Name,
		s.AssetId,
		s.Service,
		s.Method)
}
