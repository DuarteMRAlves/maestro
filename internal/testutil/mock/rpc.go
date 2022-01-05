package mock

import (
	"github.com/DuarteMRAlves/maestro/internal/reflection"
)

// RPC is a mock struct that implements the reflection.RPC interface to
// allow for easy testing.
type RPC struct {
	Name_    string
	FQN      string
	Service_ reflection.Service
	In       reflection.Message
	Out      reflection.Message
}

func (r *RPC) Name() string {
	return r.Name_
}

func (r *RPC) FullyQualifiedName() string {
	return r.FQN
}

func (r *RPC) Service() reflection.Service {
	return r.Service_
}

func (r *RPC) Input() reflection.Message {
	return r.In
}

func (r *RPC) Output() reflection.Message {
	return r.Out
}
