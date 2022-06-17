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

func (s *Stage) Dialer() Dialer {
	if s == nil {
		return nil
	}
	return s.method
}

func (s *Stage) InputDesc() MessageDesc {
	if s == nil {
		return nil
	}
	return s.method.Input()
}

func (s *Stage) OutputDesc() MessageDesc {
	if s == nil {
		return nil
	}
	return s.method.Output()
}

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
	Dialer
	Input() MessageDesc
	Output() MessageDesc
}

type Dialer interface {
	Dial() (Conn, error)
}

type DialFunc func() (Conn, error)

func (fn DialFunc) Dial() (Conn, error) { return fn() }

type Conn interface {
	Call(ctx context.Context, req Message) (Message, error)
	Close() error
}
