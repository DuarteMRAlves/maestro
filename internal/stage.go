package internal

import (
	"fmt"
	"regexp"
)

var nameRegExp, _ = regexp.Compile(`^[a-zA-Z0-9]+([-:_/][a-zA-Z0-9]+)*$|^$`)

type StageName struct{ val string }

func (s StageName) Unwrap() string { return s.val }

func (s StageName) IsEmpty() bool { return s.val == "" }

func (s StageName) String() string {
	return s.val
}

func NewStageName(name string) (StageName, error) {
	if !isValidStageName(name) {
		return StageName{}, &invalidStageName{name: name}
	}
	return StageName{val: name}, nil
}

func isValidStageName(name string) bool {
	return nameRegExp.MatchString(name)
}

type invalidStageName struct{ name string }

func (err *invalidStageName) Error() string {
	return fmt.Sprintf("invalid stage name: '%s'", err.name)
}

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

type Stage struct {
	name      StageName
	methodCtx MethodContext
}

func (s Stage) Name() StageName {
	return s.name
}

func (s Stage) MethodContext() MethodContext {
	return s.methodCtx
}

func NewStage(
	name StageName,
	methodCtx MethodContext,
) Stage {
	return Stage{
		name:      name,
		methodCtx: methodCtx,
	}
}
