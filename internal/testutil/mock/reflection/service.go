package reflection

import "github.com/DuarteMRAlves/maestro/internal/reflection"

// Service is a mock that implements the reflection.Service interface to allow
// for easy testing.
type Service struct {
	Name_ string
	FQN   string
	RPCs_ []reflection.RPC
}

func (s *Service) Name() string {
	return s.Name_
}

func (s *Service) FullyQualifiedName() string {
	return s.FQN
}

func (s *Service) RPCs() []reflection.RPC {
	return s.RPCs_
}
