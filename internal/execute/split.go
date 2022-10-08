package execute

import (
	"context"

	"github.com/DuarteMRAlves/maestro/internal/message"
)

type split struct {
	// fields are the names of the fields of the received message that should
	// be sent through the respective channel. If field is empty, the
	// entire message is sent.
	fields []message.Field
	// input is the channel from which to receive the messages.
	input <-chan state
	// outputs are the several channels where to send messages.
	outputs []chan<- state
}

func newSplit(
	fields []message.Field,
	input <-chan state,
	outputs []chan<- state,
) Stage {
	return &split{
		fields:  fields,
		input:   input,
		outputs: outputs,
	}
}

func (s *split) Run(ctx context.Context) error {
	for {
		var currState state
		select {
		case currState = <-s.input:
		case <-ctx.Done():
			for _, c := range s.outputs {
				close(c)
			}
			return nil
		}
		msg := currState.msg
		for i, out := range s.outputs {
			send := msg
			field := s.fields[i]
			if !field.IsUnspecified() {
				fieldMsg, err := msg.Get(field)
				if err != nil {
					return err
				}
				send = fieldMsg
			}
			sendState := newState(send)
			select {
			case out <- sendState:
			case <-ctx.Done():
				for _, c := range s.outputs {
					close(c)
				}
				return nil
			}
		}
	}
}
