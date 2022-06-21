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

type onlineMerge struct {
	// fields are the names of the fields of the generated message that should
	// be filled with the collected messages.
	fields []message.Field
	// inputs are the several input channels from which to collect the messages.
	inputs []<-chan onlineState
	// output is the channel used to send messages to the downstream stage.
	output chan<- onlineState
	// builder generates empty messages for the output type.
	builder message.Builder
	// currId is the current id being constructed.
	currId id
}

func newOnlineMerge(
	fields []message.Field,
	inputs []<-chan onlineState,
	output chan<- onlineState,
	gen message.Builder,
) Stage {
	return &onlineMerge{
		fields:  fields,
		inputs:  inputs,
		output:  output,
		builder: gen,
		currId:  0,
	}
}

func (s *onlineMerge) Run(ctx context.Context) error {
	var (
		// partial is the current message being constructed.
		partial   message.Instance
		currState onlineState
		done      bool
	)
	// latest stores the most recent state received from any channel.
	latest := make([]onlineState, 0, len(s.inputs))
	for i := 0; i < len(s.inputs); i++ {
		latest = append(latest, emptyOnlineState)
	}
	for {
		partial = s.builder.Build()
		// number of fields in the partial message that are set.
		setFields := 0
		for i, input := range s.inputs {
			currState = latest[i]
			if currState.id < s.currId {
				currState, done = s.takeUntilCurrId(ctx, input)
				if done {
					close(s.output)
					return nil
				}
				latest[i] = currState
			}
			// If the currState id is higher than the current id, we will never
			// be able to construct the message for the current id. As such, we
			// need to break this cycle, discard the work and move to the next
			// iteration with the currState id.
			if currState.id > s.currId {
				s.currId = currState.id
				break
			}
			err := partial.Set(s.fields[i], currState.msg)
			if err != nil {
				return err
			}
			setFields++
		}
		// All fields from inputs were set. The message can be sent
		if setFields == len(s.inputs) {
			sendState := newOnlineState(s.currId, partial)
			select {
			case s.output <- sendState:
			case <-ctx.Done():
				close(s.output)
				return nil
			}
			s.currId++
			for i := 0; i < len(s.inputs); i++ {
				latest[i] = emptyOnlineState
			}
		}
	}
}

func (s *onlineMerge) takeUntilCurrId(
	ctx context.Context,
	input <-chan onlineState,
) (onlineState, bool) {
	for {
		select {
		case st, more := <-input:
			if !more {
				return emptyOnlineState, true
			}
			if st.id >= s.currId {
				return st, false
			}
		case <-ctx.Done():
			return emptyOnlineState, true
		}
	}
}
