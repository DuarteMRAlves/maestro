package stage

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/types"
	"regexp"
)

var nameRegExp, _ = regexp.Compile(`^[a-zA-Z0-9]+([-:_/][a-zA-Z0-9]+)*$`)

type stageName string

func NewStageName(name string) (types.StageName, error) {
	if len(name) == 0 {
		return nil, errdefs.InvalidArgumentWithMsg("empty name")
	}
	if !isValidStageName(name) {
		return nil, errdefs.InvalidArgumentWithMsg("invalid name '%v'", name)
	}
	return stageName(name), nil
}

func (s stageName) Unwrap() string {
	return string(s)
}

func isValidStageName(name string) bool {
	return nameRegExp.MatchString(name)
}

type service string

func NewService(s string) (types.Service, error) {
	if len(s) == 0 {
		return nil, errdefs.InvalidArgumentWithMsg("empty service")
	}
	return service(s), nil
}

func (s service) Unwrap() string {
	return string(s)
}

type presentService struct{ types.Service }

func (s presentService) Unwrap() types.Service { return s.Service }

func (s presentService) Present() bool { return true }

type emptyService struct{}

func (s emptyService) Unwrap() types.Service {
	panic("Service not available in an empty service optional")
}

func (s emptyService) Present() bool { return false }

func NewPresentService(s types.Service) types.OptionalService {
	return presentService{s}
}

func NewEmptyService() types.OptionalService {
	return emptyService{}
}

type method string

func (m method) Unwrap() string {
	return string(m)
}

func NewMethod(m string) (types.Method, error) {
	if len(m) == 0 {
		return nil, errdefs.InvalidArgumentWithMsg("empty method")
	}
	return method(m), nil
}

type presentMethod struct{ types.Method }

func (m presentMethod) Unwrap() types.Method { return m.Method }

func (m presentMethod) Present() bool { return true }

type emptyMethod struct{}

func (m emptyMethod) Unwrap() types.Method {
	panic("Method not available in an empty method optional")
}

func (m emptyMethod) Present() bool { return false }

func NewPresentMethod(m types.Method) types.OptionalMethod {
	return presentMethod{m}
}

func NewEmptyMethod() types.OptionalMethod {
	return emptyMethod{}
}

type address string

func (a address) Unwrap() string { return string(a) }

func NewAddress(a string) (types.Address, error) {
	if len(a) == 0 {
		return nil, errdefs.InvalidArgumentWithMsg("empty address")
	}
	return address(a), nil
}

type methodContext struct {
	address types.Address
	service types.OptionalService
	method  types.OptionalMethod
}

func (m methodContext) Address() types.Address { return m.address }

func (m methodContext) Service() types.OptionalService { return m.service }

func (m methodContext) Method() types.OptionalMethod { return m.method }

func NewMethodContext(
	address types.Address,
	service types.OptionalService,
	method types.OptionalMethod,
) types.MethodContext {
	return methodContext{
		address: address,
		service: service,
		method:  method,
	}
}

type stage struct {
	name      types.StageName
	methodCtx types.MethodContext
}

func (s stage) Name() types.StageName {
	return s.name
}

func (s stage) MethodContext() types.MethodContext {
	return s.methodCtx
}

func NewStage(name types.StageName, methodCtx types.MethodContext) types.Stage {
	return stage{
		name:      name,
		methodCtx: methodCtx,
	}
}
