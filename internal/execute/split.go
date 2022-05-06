package execute

import (
	"context"

	"github.com/DuarteMRAlves/maestro/internal"
)

type offlineSplit struct {
	// fields are the names of the fields of the received message that should
	// be sent through the respective channel. If field is empty, the
	// entire message is sent.
	fields []internal.MessageField
	// input is the channel from which to receive the messages.
	input <-chan offlineState
	// outputs are the several channels where to send messages.
	outputs []chan<- offlineState
}

func newOfflineSplit(
	fields []internal.MessageField,
	input <-chan offlineState,
	outputs []chan<- offlineState,
) Stage {
	return &offlineSplit{
		fields:  fields,
		input:   input,
		outputs: outputs,
	}
}

func (s *offlineSplit) Run(ctx context.Context) error {
	for {
		var currState offlineState
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
			if !field.IsEmpty() {
				fieldMsg, err := msg.GetField(field)
				if err != nil {
					return err
				}
				send = fieldMsg
			}
			sendState := newOfflineState(send)
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

type onlineSplit struct {
	// fields are the names of the fields of the received message that should
	// be sent through the respective channel. If field is empty, the
	// entire message is sent.
	fields []internal.MessageField
	// input is the channel from which to receive the messages.
	input <-chan onlineState
	// outputs are the several channels where to send messages.
	outputs []chan<- onlineState
}

func newOnlineSplit(
	fields []internal.MessageField,
	input <-chan onlineState,
	outputs []chan<- onlineState,
) Stage {
	return &onlineSplit{
		fields:  fields,
		input:   input,
		outputs: outputs,
	}
}

func (s *onlineSplit) Run(ctx context.Context) error {
	var currState onlineState
	for {
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
			if !field.IsEmpty() {
				fieldMsg, err := msg.GetField(field)
				if err != nil {
					return err
				}
				send = fieldMsg
			}
			sendState := newOnlineState(currState.id, send)
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
