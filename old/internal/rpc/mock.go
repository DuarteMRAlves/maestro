package rpc

// MockRPC is a mock struct that implements the rpc.MockRPC interface to
// allow for easy testing.
type MockRPC struct {
	Name_    string
	FQN      string
	Invoke   string
	Service_ Service
	In       MessageDesc
	Out      MessageDesc
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

func (r *MockRPC) Input() MessageDesc {
	return r.In
}

func (r *MockRPC) Output() MessageDesc {
	return r.Out
}

func (r *MockRPC) IsUnary() bool {
	return r.Unary
}
