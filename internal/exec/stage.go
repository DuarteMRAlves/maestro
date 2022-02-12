package exec

import (
	"context"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/rpc"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/dynamic"
	"google.golang.org/grpc"
	"io"
	"time"
)

// Stage executes a stage in an orchestration. A stage can be a rpc stage, that
// executes an rpc and was specified by the user, or an auxiliary stage to
// control the flow of the pipeline, such as source or a sink.
type Stage interface {
	// Run executes the given stage with a config. This function terminates when
	// the term channel is signalled. The function signals the done channel as
	// the last instruction before returning.
	Run(*RunCfg)
	// Close shuts down the output channels of these stages, cleaning up the
	// pipeline.
	Close()
}

// RunCfg specifies the configuration that the Stage should use when running.
type RunCfg struct {
	// term is a channel that will be signaled if the Stage should stop.
	term <-chan struct{}
	// done is a channel that the Stage should close to signal is has finished.
	done chan<- struct{}
	// errs is a channel that the worker should use to send errors in order to
	// be processed.
	// The io.EOF error should not be sent through this channel at is just a
	// termination signal
	errs chan<- error
}

func NewRpcStage(
	address string,
	rpcDesc rpc.RPC,
	input <-chan *State,
	output chan<- *State,
) (Stage, error) {
	switch {
	case rpcDesc.IsUnary():
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			return nil, errdefs.InvalidArgumentWithMsg(
				"unable to connect to address: %s",
				address,
			)
		}
		w := &UnaryStage{
			Address: address,
			conn:    conn,
			rpc:     rpcDesc,
			invoker: rpc.NewUnary(rpcDesc.InvokePath(), conn),
			input:   input,
			output:  output,
		}
		return w, nil
	default:
		return nil, errdefs.InvalidArgumentWithMsg("unsupported rpc type")
	}
}

// UnaryStage manages the execution of a stage in a pipeline.
type UnaryStage struct {
	Address string
	conn    grpc.ClientConnInterface
	rpc     rpc.RPC
	invoker rpc.UnaryClient

	input  <-chan *State
	output chan<- *State
}

func (s *UnaryStage) Run(cfg *RunCfg) {
	var (
		in, out  *State
		req, rep interface{}
		err      error
	)

	for {
		select {
		case in = <-s.input:
		case <-cfg.term:
			close(cfg.done)
			return
		}
		if in.Err() == io.EOF {
			close(cfg.done)
			return
		}
		if in.Err() != nil {
			cfg.errs <- in.Err()
			continue
		}
		req = in.Msg()
		rep = s.rpc.Output().NewEmpty()

		err = s.invoke(req, rep)
		if err != nil {
			cfg.errs <- err
			continue
		}

		out = NewState(in.Id(), rep)

		select {
		case s.output <- out:
		case <-cfg.term:
			close(cfg.done)
			return
		}
	}
}

func (s *UnaryStage) invoke(req interface{}, rep interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return s.invoker.Invoke(ctx, req, rep)
}

func (s *UnaryStage) Close() {
	close(s.output)
}

// SourceStage is the source of the orchestration. It defines the initial ids of
// the states and sends empty messages of the received type.
type SourceStage struct {
	id     int32
	msg    rpc.Message
	output chan<- *State
}

func NewSourceStage(
	initial int32,
	output chan<- *State,
	msg rpc.Message,
) *SourceStage {
	i := &SourceStage{
		id:     initial,
		msg:    msg,
		output: output,
	}
	return i
}

func (s *SourceStage) Run(cfg *RunCfg) {
	for {
		select {
		case s.output <- s.next():
		case <-cfg.term:
			close(cfg.done)
			return
		}
	}
}

func (s *SourceStage) next() *State {
	st := NewState(Id(s.id), s.msg.NewEmpty())
	s.id++
	return st
}

func (s *SourceStage) Close() {
	close(s.output)
}

// SinkStage defines the last output of the orchestration, where all messages
// are dropped.
type SinkStage struct {
	input <-chan *State
}

func NewSinkOutput(input <-chan *State) *SinkStage {
	s := &SinkStage{
		input: input,
	}
	return s
}

func (s *SinkStage) Run(cfg *RunCfg) {
	for {
		select {
		// Discard results
		case <-s.input:
		case <-cfg.term:
			close(cfg.done)
			return
		}
	}
}

func (s *SinkStage) Close() {
	// Do nothing, the input channel will be closed by the upstream stage.
}

// MergeStage collects multiple messages from multiple channels and build a
// single message that sends to the downstream stage.
type MergeStage struct {
	// fields are the names of the fields of the generated message that should
	// be filled with the collected messages.
	fields []string
	// inputs are the several input channels from which to collect the messages.
	inputs []<-chan *State
	// output is the channel used to send messages to the downstream stage.
	output chan<- *State
	// msg describes the message to create and send to the downstream stage.
	msg rpc.Message
	// currId is the current id being constructed.
	currId Id
}

func NewMergeStage(
	fields []string,
	inputs []<-chan *State,
	output chan<- *State,
	msg rpc.Message,
) *MergeStage {
	return &MergeStage{
		fields: fields,
		inputs: inputs,
		output: output,
		msg:    msg,
		currId: 0,
	}
}

func (s *MergeStage) Run(cfg *RunCfg) {
	var (
		// partial is the current message being constructed.
		partial *dynamic.Message
		state   *State
		done    bool
	)

	latest := make([]*State, 0, len(s.inputs))
	for i := 0; i < len(s.inputs); i++ {
		latest = append(latest, nil)
	}
	for {
		partial = s.msg.NewEmpty()
		setFields := 0
		for i, input := range s.inputs {
			state = latest[i]
			if state == nil || state.Id() < s.currId {
				state, done = s.takeUntilCurrId(input, cfg.term)
				if done {
					close(cfg.done)
					return
				}
				latest[i] = state
			}
			// The message with the current id was discarded. The number of set
			// fields is smaller than the number of inputs and so the message
			// will not be discarded.
			if state.Id() > s.currId {
				s.currId = state.Id()
				break
			}
			partial.SetFieldByName(s.fields[i], state.Msg())
			setFields++
		}
		// All fields from inputs were set. The message can be sent
		if setFields == len(s.inputs) {
			sendState := NewState(s.currId, partial)
			select {
			case s.output <- sendState:
			case <-cfg.term:
				close(cfg.done)
				return
			}
			s.currId++
			for i := 0; i < len(s.inputs); i++ {
				latest[i] = nil
			}
		}
	}
}

func (s *MergeStage) takeUntilCurrId(
	input <-chan *State,
	term <-chan struct{},
) (*State, bool) {
	for {
		select {
		case state := <-input:
			if state.Id() >= s.currId {
				return state, false
			}
		case <-term:
			return nil, true
		}
	}
}

func (s *MergeStage) Close() {
	close(s.output)
}

// SplitStage divides a stage output into multiple channels. It can send the
// entire message, or a field.
type SplitStage struct {
	// fields are the names of the fields of the received message that should
	// be sent through the respective channel. If field is the empty string, the
	// entire message is sent.
	fields []string
	// input is the channel from which to receive the messages.
	input <-chan *State
	// outputs are the several channels where to send messages.
	outputs []chan<- *State
}

func NewSplitStage(
	fields []string,
	input <-chan *State,
	outputs []chan<- *State,
) *SplitStage {
	return &SplitStage{
		fields:  fields,
		input:   input,
		outputs: outputs,
	}
}

func (s *SplitStage) Run(cfg *RunCfg) {
	var (
		state *State
		send  interface{}
	)
	for {
		select {
		case state = <-s.input:
		case <-cfg.term:
			close(cfg.done)
			return
		}
		msg, ok := state.Msg().(proto.Message)
		if !ok {
			cfg.errs <- errdefs.InternalWithMsg("Invalid message type")
			continue
		}
		dyn, err := dynamic.AsDynamicMessage(msg)
		if err != nil {
			cfg.errs <- errdefs.InternalWithMsg(
				"convert proto msg to dynamic: %s",
				err,
			)
			continue
		}
		for i, out := range s.outputs {
			send = dyn
			field := s.fields[i]
			if field != "" {
				send, err = dyn.TryGetFieldByName(field)
				if err != nil {
					cfg.errs <- errdefs.InternalWithMsg(
						"get field '%s': %s",
						field,
						err,
					)
					continue
				}
			}
			newState := NewState(state.Id(), send)
			select {
			case out <- newState:
			case <-cfg.term:
				close(cfg.done)
				return
			}
		}
	}
}

func (s *SplitStage) Close() {
	for _, c := range s.outputs {
		close(c)
	}
}
