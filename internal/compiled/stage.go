package compiled

import (
	"context"
	"fmt"
)

// Stage defines a step of a Pipeline
type Stage struct {
	name  StageName
	sType StageType

	// static attributes for the method invocation
	mid MethodID

	// runtime attributes that can be computed from
	// the static attributes
	method MethodDesc

	// define the connections for this stage.
	inputs  []*Link
	outputs []*Link
}

func (s *Stage) Name() StageName {
	return s.name
}

func (s *Stage) Type() StageType { return s.sType }

func (s *Stage) Inputs() []*Link {
	return s.inputs
}

func (s *Stage) Outputs() []*Link {
	return s.outputs
}

type StageName struct{ val string }

func (s StageName) Unwrap() string { return s.val }

func (s StageName) IsEmpty() bool { return s.val == "" }

func (s StageName) String() string {
	return s.val
}

func NewStageName(name string) (StageName, error) {
	if !validateResourceName(name) {
		return StageName{}, &invalidStageName{name: name}
	}
	return StageName{val: name}, nil
}

type invalidStageName struct{ name string }

func (err *invalidStageName) Error() string {
	return fmt.Sprintf("invalid stage name: '%s'", err.name)
}

type StageType string

const (
	StageTypeUnary  StageType = "UnaryStage"
	StageTypeSource StageType = "SourceStage"
	StageTypeSink   StageType = "SinkStage"
	StageTypeMerge  StageType = "MergeStage"
	StageTypeSplit  StageType = "SplitStage"
)

// MethodID uniquely identifies a given method.
type MethodID interface {
	String() string
}

// MethodDesc contains the information to create a method.
type MethodDesc interface {
	Dial() Conn
	Input() MessageDesc
	Output() MessageDesc
}

type Conn interface {
	Call(ctx context.Context, req Message) (Message, error)
	Close() error
}

// Address specifies the location of the server executing the
// stage method.
type Address string

func (a Address) IsEmpty() bool { return a == "" }

func (a Address) String() string {
	if a.IsEmpty() {
		return "*"
	}
	return string(a)
}

// Service specifies the name of the grpc service to execute.
type Service string

// IsUnspecified reports whether this service is either "" or "*".
func (s Service) IsUnspecified() bool { return s == "" || s == "*" }

func (s Service) String() string {
	if s.IsUnspecified() {
		return "*"
	}
	return string(s)
}

// Method specified the name of the grpc method to execute.
type Method string

// IsUnspecified reports whether this method is either "" or "*".
func (m Method) IsUnspecified() bool { return m == "" || m == "*" }

func (m Method) String() string {
	if m.IsUnspecified() {
		return "*"
	}
	return string(m)
}

type MethodContext struct {
	address Address
	service Service
	method  Method
}

func (m MethodContext) Address() Address { return m.address }

func (m MethodContext) Service() Service { return m.service }

func (m MethodContext) Method() Method { return m.method }

func (m MethodContext) String() string {
	return fmt.Sprintf("%s/%s/%s", m.address, m.service, m.method)
}

func NewMethodContext(
	address Address,
	service Service,
	method Method,
) MethodContext {
	return MethodContext{
		address: address,
		service: service,
		method:  method,
	}
}
