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

type FieldNotMessage struct {
	MsgType string
	Field   string
}

func (err *FieldNotMessage) Error() string {
	format := "field '%s' is not a sub message of %s"
	return fmt.Sprintf(format, err.Field, err.MsgType)
}