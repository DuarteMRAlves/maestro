package execute

import (
	"context"

	"github.com/DuarteMRAlves/maestro/internal/message"
)

type offlineMerge struct {
	// fields are the names of the fields of the generated message that should
	// be filled with the collected messages.
	fields []message.Field
	// inputs are the several input channels from which to collect the messages.
	inputs []<-chan offlineState
	// output is the channel used to send messages to the downstream stage.
	output chan<- offlineState
	// builder generates empty messages for the output type.
	builder message.Builder
}

func newOfflineMerge(
	fields []message.Field,
	inputs []<-chan offlineState,
	output chan<- offlineState,
	gen message.Builder,
) Stage {
	return &offlineMerge{
		fields:  fields,
		inputs:  inputs,
		output:  output,
		builder: gen,
	}
}

func (s *offlineMerge) Run(ctx context.Context) error {
	for {
		var (
			currState offlineState
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
		sendState := newOfflineState(partial)
		select {
		case s.output <- sendState:
		case <-ctx.Done():
			close(s.output)
			return nil
		}
	}
}
