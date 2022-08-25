package compiled

import (
	"fmt"

	"github.com/DuarteMRAlves/maestro/internal/message"
	"github.com/DuarteMRAlves/maestro/internal/method"
)

// Stage defines a step of a Pipeline
type Stage struct {
	name  StageName
	sType StageType

	// static attributes for the method invocation
	address string

	// runtime attributes that can be computed from
	// the static attributes
	desc method.Desc

	// define the connections for this stage.
	inputs  []*Link
	outputs []*Link
}

func (s *Stage) Name() StageName {
	if s == nil {
		return StageName{}
	}
	return s.name
}

func (s *Stage) Type() StageType {
	if s == nil {
		return StageTypeUnknown
	}
	return s.sType
}

func (s *Stage) Dialer() method.Dialer {
	if s == nil {
		return nil
	}
	return s.desc
}

func (s *Stage) InputDesc() message.Type {
	if s == nil {
		return nil
	}
	return s.desc.Input()
}

func (s *Stage) OutputDesc() message.Type {
	if s == nil {
		return nil
	}
	return s.desc.Output()
}

// Iterates over the input links of Stage while fn returns true.
func (s *Stage) RangeInputs(fn func(*Link) bool) {
	if s == nil {
		return
	}
	for _, i := range s.inputs {
		if !fn(i) {
			return
		}
	}
}

func (s *Stage) CopyInputs() []*Link {
	if s == nil {
		return nil
	}
	cp := make([]*Link, 0, len(s.inputs))
	s.RangeInputs(func(l *Link) bool {
		cp = append(cp, l)
		return true
	})
	return cp
}

// Iterates over the output links of Stage while fn returns true.
func (s *Stage) RangeOutputs(fn func(*Link) bool) {
	if s == nil {
		return
	}
	for _, i := range s.outputs {
		if !fn(i) {
			return
		}
	}
}

func (s *Stage) CopyOutputs() []*Link {
	if s == nil {
		return nil
	}
	cp := make([]*Link, 0, len(s.outputs))
	s.RangeOutputs(func(l *Link) bool {
		cp = append(cp, l)
		return true
	})
	return cp
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
	StageTypeUnknown StageType = "UnknownStage"
	StageTypeUnary   StageType = "UnaryStage"
	StageTypeSource  StageType = "SourceStage"
	StageTypeSink    StageType = "SinkStage"
	StageTypeMerge   StageType = "MergeStage"
	StageTypeSplit   StageType = "SplitStage"
)
