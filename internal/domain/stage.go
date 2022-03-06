package domain

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"regexp"
)

var nameRegExp, _ = regexp.Compile(`^[a-zA-Z0-9]+([-:_/][a-zA-Z0-9]+)*$`)

type stageName string

func (s stageName) StageName() {}

func (s stageName) Unwrap() string {
	return string(s)
}

func NewStageName(name string) (StageName, error) {
	if len(name) == 0 {
		return nil, errdefs.InvalidArgumentWithMsg("empty name")
	}
	if !isValidStageName(name) {
		return nil, errdefs.InvalidArgumentWithMsg("invalid name '%v'", name)
	}
	return stageName(name), nil
}

func isValidStageName(name string) bool {
	return nameRegExp.MatchString(name)
}

type service string

func (s service) Service() {}

func (s service) Unwrap() string {
	return string(s)
}

func NewService(s string) (Service, error) {
	if len(s) == 0 {
		return nil, errdefs.InvalidArgumentWithMsg("empty service")
	}
	return service(s), nil
}

type presentService struct{ Service }

func (s presentService) Unwrap() Service { return s.Service }

func (s presentService) Present() bool { return true }

type emptyService struct{}

func (s emptyService) Unwrap() Service {
	panic("Service not available in an empty service optional")
}

func (s emptyService) Present() bool { return false }

func NewPresentService(s Service) OptionalService {
	return presentService{s}
}

func NewEmptyService() OptionalService {
	return emptyService{}
}

type method string

func (m method) Method() {}

func (m method) Unwrap() string {
	return string(m)
}

func NewMethod(m string) (Method, error) {
	if len(m) == 0 {
		return nil, errdefs.InvalidArgumentWithMsg("empty method")
	}
	return method(m), nil
}

type presentMethod struct{ Method }

func (m presentMethod) Unwrap() Method { return m.Method }

func (m presentMethod) Present() bool { return true }

type emptyMethod struct{}

func (m emptyMethod) Unwrap() Method {
	panic("Method not available in an empty method optional")
}

func (m emptyMethod) Present() bool { return false }

func NewPresentMethod(m Method) OptionalMethod {
	return presentMethod{m}
}

func NewEmptyMethod() OptionalMethod {
	return emptyMethod{}
}

type address string

func (a address) Address() {}

func (a address) Unwrap() string { return string(a) }

func NewAddress(a string) (Address, error) {
	if len(a) == 0 {
		return nil, errdefs.InvalidArgumentWithMsg("empty address")
	}
	return address(a), nil
}

type methodContext struct {
	address Address
	service OptionalService
	method  OptionalMethod
}

func (m methodContext) Address() Address { return m.address }

func (m methodContext) Service() OptionalService { return m.service }

func (m methodContext) Method() OptionalMethod { return m.method }

func NewMethodContext(
	address Address,
	service OptionalService,
	method OptionalMethod,
) MethodContext {
	return methodContext{
		address: address,
		service: service,
		method:  method,
	}
}

type stage struct {
	name      StageName
	methodCtx MethodContext
}

func (s stage) Name() StageName {
	return s.name
}

func (s stage) MethodContext() MethodContext {
	return s.methodCtx
}

func NewStage(
	name StageName,
	methodCtx MethodContext,
) Stage {
	return stage{
		name:      name,
		methodCtx: methodCtx,
	}
}
