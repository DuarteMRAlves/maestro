package execute

import (
	"context"

	"github.com/DuarteMRAlves/maestro/internal/message"
)

type merge struct {
	// fields are the names of the fields of the generated message that should
	// be filled with the collected messages.
	fields []message.Field
	// inputs are the several input channels from which to collect the messages.
	inputs []<-chan state
	// output is the channel used to send messages to the downstream stage.
	output chan<- state
	// builder generates empty messages for the output type.
	builder message.Builder
}

func newMerge(
	fields []message.Field,
	inputs []<-chan state,
	output chan<- state,
	gen message.Builder,
) Stage {
	return &merge{
		fields:  fields,
		inputs:  inputs,
		output:  output,
		builder: gen,
	}
}

func (s *merge) Run(ctx context.Context) error {
	for {
		var (
			currState state
			more      bool
		)
		// partial is the current message being constructed.
		partial := s.builder.Build()
		for i, input := range s.inputs {
			select {
			case currState, more = <-input:
			case <-ctx.Done():
				close(s.output)
				return nil
			}
			if !more {
				close(s.output)
				return nil
			}
			err := partial.Set(s.fields[i], currState.msg)
			if err != nil {
				return err
			}
		}
		sendState := newState(partial)
		select {
		case s.output <- sendState:
		case <-ctx.Done():
			close(s.output)
			return nil
		}
	}
}
