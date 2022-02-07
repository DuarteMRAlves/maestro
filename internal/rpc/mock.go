package rpc

import (
	"context"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"google.golang.org/grpc"
	"sync"
)

type MockManager struct {
	Rpcs sync.Map
}

func (m *MockManager) GetRpc(
	_ context.Context,
	_ grpc.ClientConnInterface,
	stage api.StageName,
) (RPC, bool) {
	rpc, ok := m.Rpcs.Load(stage)
	if !ok {
		return nil, false
	}
	return rpc.(RPC), true
}

// MockService is a mock that implements the rpc.MockService interface to allow
// for easy testing.
type MockService struct {
	Name_ string
	FQN   string
	RPCs_ []RPC
}

func (s *MockService) Name() string {
	return s.Name_
}

func (s *MockService) FullyQualifiedName() string {
	return s.FQN
}

func (s *MockService) RPCs() []RPC {
	return s.RPCs_
}

// MockRPC is a mock struct that implements the rpc.MockRPC interface to
// allow for easy testing.
type MockRPC struct {
	Name_    string
	FQN      string
	Invoke   string
	Service_ Service
	In       Message
	Out      Message
	Unary    bool
}

func (r *MockRPC) Name() string {
	return r.Name_
}

func (r *MockRPC) FullyQualifiedName() string {
	return r.FQN
}

func (r *MockRPC) InvokePath() string {
	return r.Invoke
}

func (r *MockRPC) Service() Service {
	return r.Service_
}

func (r *MockRPC) Input() Message {
	return r.In
}

func (r *MockRPC) Output() Message {
	return r.Out
}

func (r *MockRPC) IsUnary() bool {
	return r.Unary
}
