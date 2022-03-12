package domain

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"regexp"
)

var nameRegExp, _ = regexp.Compile(`^[a-zA-Z0-9]+([-:_/][a-zA-Z0-9]+)*$|^$`)

type StageName struct{ val string }

func (s StageName) Unwrap() string { return s.val }

func (s StageName) IsEmpty() bool { return s.val == "" }

func NewStageName(name string) (StageName, error) {
	if !isValidStageName(name) {
		err := errdefs.InvalidArgumentWithMsg("invalid name '%v'", name)
		return StageName{}, err
	}
	return StageName{val: name}, nil
}

func isValidStageName(name string) bool {
	return nameRegExp.MatchString(name)
}

type Service struct{ val string }

func (s Service) Unwrap() string {
	return s.val
}

func (s Service) IsEmpty() bool { return s.val == "" }

func NewService(s string) Service {
	return Service{val: s}
}

type OptionalService struct {
	val     Service
	present bool
}

func (o OptionalService) Unwrap() Service { return o.val }

func (o OptionalService) Present() bool { return o.present }

func NewPresentService(s Service) OptionalService {
	return OptionalService{val: s, present: true}
}

func NewEmptyService() OptionalService {
	return OptionalService{}
}

type Method struct{ val string }

func (m Method) Unwrap() string {
	return m.val
}

func (m Method) IsEmpty() bool { return m.val == "" }

func NewMethod(m string) Method { return Method{val: m} }

type OptionalMethod struct {
	val     Method
	present bool
}

func (m OptionalMethod) Unwrap() Method { return m.val }

func (m OptionalMethod) Present() bool { return m.present }

func NewPresentMethod(m Method) OptionalMethod {
	return OptionalMethod{val: m, present: true}
}

func NewEmptyMethod() OptionalMethod {
	return OptionalMethod{}
}

type Address struct{ val string }

func (a Address) Unwrap() string { return a.val }

func (a Address) IsEmpty() bool { return a.val == "" }

func NewAddress(a string) Address { return Address{val: a} }

type MethodContext struct {
	address Address
	service OptionalService
	method  OptionalMethod
}

func (m MethodContext) Address() Address { return m.address }

func (m MethodContext) Service() OptionalService { return m.service }

func (m MethodContext) Method() OptionalMethod { return m.method }

func NewMethodContext(
	address Address,
	service OptionalService,
	method OptionalMethod,
) MethodContext {
	return MethodContext{
		address: address,
		service: service,
		method:  method,
	}
}

type Stage struct {
	name          StageName
	methodCtx     MethodContext
	orchestration OrchestrationName
}

func (s Stage) Name() StageName {
	return s.name
}

func (s Stage) MethodContext() MethodContext {
	return s.methodCtx
}

func (s Stage) Orchestration() OrchestrationName {
	return s.orchestration
}

func NewStage(
	name StageName,
	methodCtx MethodContext,
	orchestration OrchestrationName,
) Stage {
	return Stage{
		name:          name,
		methodCtx:     methodCtx,
		orchestration: orchestration,
	}
}
