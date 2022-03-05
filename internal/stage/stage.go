package stage

import (
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"regexp"
)

var nameRegExp, _ = regexp.Compile(`^[a-zA-Z0-9]+([-:_/][a-zA-Z0-9]+)*$`)

type stageName string

func NewStageName(name string) (domain.StageName, error) {
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

func NewService(s string) (domain.Service, error) {
	if len(s) == 0 {
		return nil, errdefs.InvalidArgumentWithMsg("empty service")
	}
	return service(s), nil
}

func (s service) Unwrap() string {
	return string(s)
}

type presentService struct{ domain.Service }

func (s presentService) Unwrap() domain.Service { return s.Service }

func (s presentService) Present() bool { return true }

type emptyService struct{}

func (s emptyService) Unwrap() domain.Service {
	panic("Service not available in an empty service optional")
}

func (s emptyService) Present() bool { return false }

func NewPresentService(s domain.Service) domain.OptionalService {
	return presentService{s}
}

func NewEmptyService() domain.OptionalService {
	return emptyService{}
}

type method string

func (m method) Unwrap() string {
	return string(m)
}

func NewMethod(m string) (domain.Method, error) {
	if len(m) == 0 {
		return nil, errdefs.InvalidArgumentWithMsg("empty method")
	}
	return method(m), nil
}

type presentMethod struct{ domain.Method }

func (m presentMethod) Unwrap() domain.Method { return m.Method }

func (m presentMethod) Present() bool { return true }

type emptyMethod struct{}

func (m emptyMethod) Unwrap() domain.Method {
	panic("Method not available in an empty method optional")
}

func (m emptyMethod) Present() bool { return false }

func NewPresentMethod(m domain.Method) domain.OptionalMethod {
	return presentMethod{m}
}

func NewEmptyMethod() domain.OptionalMethod {
	return emptyMethod{}
}

type address string

func (a address) Unwrap() string { return string(a) }

func NewAddress(a string) (domain.Address, error) {
	if len(a) == 0 {
		return nil, errdefs.InvalidArgumentWithMsg("empty address")
	}
	return address(a), nil
}

type methodContext struct {
	address domain.Address
	service domain.OptionalService
	method  domain.OptionalMethod
}

func (m methodContext) Address() domain.Address { return m.address }

func (m methodContext) Service() domain.OptionalService { return m.service }

func (m methodContext) Method() domain.OptionalMethod { return m.method }

func NewMethodContext(
	address domain.Address,
	service domain.OptionalService,
	method domain.OptionalMethod,
) domain.MethodContext {
	return methodContext{
		address: address,
		service: service,
		method:  method,
	}
}

type stage struct {
	name      domain.StageName
	methodCtx domain.MethodContext
}

func (s stage) Name() domain.StageName {
	return s.name
}

func (s stage) MethodContext() domain.MethodContext {
	return s.methodCtx
}

func NewStage(
	name domain.StageName,
	methodCtx domain.MethodContext,
) domain.Stage {
	return stage{
		name:      name,
		methodCtx: methodCtx,
	}
}
