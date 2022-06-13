package compiled

import (
	"fmt"
)

// Stage defines a step of a Pipeline
type Stage struct {
	name  StageName
	sType StageType

	// ictx stores all the necessary information
	// to invoke the method for this stage. In the case of aux
	// stages, this information is also filled to allow for
	// generation of messages.
	ictx *InvocationContext

	// define the connections for this stage.
	inputs  []*Link
	outputs []*Link
}

func (s *Stage) Name() StageName {
	return s.name
}

func (s *Stage) Type() StageType { return s.sType }

func (s *Stage) InvocationContext() *InvocationContext {
	return s.ictx
}

func (s *Stage) Inputs() []*Link {
	return s.inputs
}

func (s *Stage) Outputs() []*Link {
	return s.outputs
}

// InvocationContext stores all necessary information to invoke a
// remote method.
type InvocationContext struct {
	// static attributes for the method invocation
	address Address
	service Service
	method  Method

	// runtime attributes that can be computed from
	// the static attributes
	unaryMethod UnaryMethod
}

func newInvocationContext(
	methodLoader MethodLoader, address Address, service Service, method Method,
) (*InvocationContext, error) {
	methodCtx := NewMethodContext(address, service, method)
	unaryMethod, err := methodLoader.Load(&methodCtx)
	if err != nil {
		return nil, fmt.Errorf("load method context %s: %w", methodCtx, err)
	}
	ictx := &InvocationContext{
		address:     address,
		service:     service,
		method:      method,
		unaryMethod: unaryMethod,
	}
	return ictx, nil
}

func (ictx *InvocationContext) Address() Address {
	if ictx == nil {
		return Address{}
	}
	return ictx.address
}

func (ictx *InvocationContext) Service() Service {
	if ictx == nil {
		return Service{}
	}
	return ictx.service
}

func (ictx *InvocationContext) Method() Method {
	if ictx == nil {
		return Method{}
	}
	return ictx.method
}

func (ictx *InvocationContext) ClientBuilder() UnaryClientBuilder {
	if ictx == nil {
		return nil
	}
	return ictx.unaryMethod.ClientBuilder()
}

func (ictx *InvocationContext) Input() MessageDesc {
	if ictx == nil {
		return nil
	}
	return ictx.unaryMethod.Input()
}

func (ictx *InvocationContext) Output() MessageDesc {
	if ictx == nil {
		return nil
	}
	return ictx.unaryMethod.Output()
}

func (ictx *InvocationContext) String() string {
	methodInfo := joinAddressServiceMethod(ictx.address, ictx.service, ictx.method)
	return fmt.Sprintf("InvokationContext{%s}", methodInfo)
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

type Service struct{ val string }

func (s Service) Unwrap() string {
	return s.val
}

func (s Service) IsEmpty() bool { return s.val == "" }

func NewService(s string) Service {
	return Service{val: s}
}

type Method struct{ val string }

func (m Method) Unwrap() string {
	return m.val
}

func (m Method) IsEmpty() bool { return m.val == "" }

func NewMethod(m string) Method { return Method{val: m} }

type Address struct{ val string }

func (a Address) Unwrap() string { return a.val }

func (a Address) IsEmpty() bool { return a.val == "" }

func NewAddress(a string) Address { return Address{val: a} }

type MethodContext struct {
	address Address
	service Service
	method  Method
}

func (m MethodContext) Address() Address { return m.address }

func (m MethodContext) Service() Service { return m.service }

func (m MethodContext) Method() Method { return m.method }

func (m MethodContext) String() string {
	return joinAddressServiceMethod(m.address, m.service, m.method)
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

func joinAddressServiceMethod(address Address, service Service, method Method) string {
	addr := "*"
	srv := "*"
	meth := "*"
	if !address.IsEmpty() {
		addr = address.val
	}
	if !service.IsEmpty() {
		srv = service.val
	}
	if !method.IsEmpty() {
		meth = method.val
	}
	return fmt.Sprintf("%s/%s/%s", addr, srv, meth)
}
