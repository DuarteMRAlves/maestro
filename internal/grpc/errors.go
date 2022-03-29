package grpc

import (
	"errors"
	"fmt"
)

var (
	notGrpcMessage = errors.New("message is not grpc")
	notOneService  = errors.New("expected 1 available service")
	notOneMethod   = errors.New("expected 1 available method")
)

type fieldNotMessage struct {
	MsgType string
	Field   string
}

func (err *fieldNotMessage) Error() string {
	format := "field '%s' is not a sub message of %s"
	return fmt.Sprintf(format, err.Field, err.MsgType)
}

type fieldNotFound struct {
	msgType string
	field   string
}

func (err *fieldNotFound) NotFound() {}

func (err *fieldNotFound) Error() string {
	return fmt.Sprintf("field '%s' not found in msg %s", err.field, err.msgType)
}

type serviceNotFound struct {
	srv string
}

func (err *serviceNotFound) NotFound() {}

func (err *serviceNotFound) Error() string {
	return fmt.Sprintf("service not found: %s", err.srv)
}

type methodNotFound struct {
	meth string
}

func (err *methodNotFound) NotFound() {}

func (err *methodNotFound) Error() string {
	return fmt.Sprintf("method not found: %s", err.meth)
}
