package exec

import (
	"context"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/rpc"
	"google.golang.org/grpc"
	"io"
	"time"
)

type Stage interface {
	Run(*RunCfg)
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

type StageCfg struct {
	Address string
	Rpc     rpc.RPC
	Input   <-chan *State
	Output  chan<- *State
}

func (c *StageCfg) Clone() *StageCfg {
	return &StageCfg{
		Address: c.Address,
		Rpc:     c.Rpc,
		Input:   c.Input,
		Output:  c.Output,
	}
}

func NewStage(cfg *StageCfg) (Stage, error) {
	cfg = cfg.Clone()
	switch {
	case cfg.Rpc.IsUnary():
		conn, err := grpc.Dial(cfg.Address, grpc.WithInsecure())
		if err != nil {
			return nil, errdefs.InvalidArgumentWithMsg(
				"unable to connect to address: %s",
				cfg.Address,
			)
		}
		w := &UnaryStage{
			Address: cfg.Address,
			conn:    conn,
			rpc:     cfg.Rpc,
			invoker: rpc.NewUnary(cfg.Rpc.InvokePath(), conn),
			input:   cfg.Input,
			output:  cfg.Output,
		}
		return w, nil
	default:
		return nil, errdefs.InvalidArgumentWithMsg("unsupported rpc type")
	}
}

// SourceStage is the source of the orchestration. It defines the initial ids of
// the states and sends empty messages of the received type.
type SourceStage struct {
	id  int32
	msg rpc.Message
	ch  chan *State
}

func NewSourceStage(
	initial int32,
	ch chan *State,
	msg rpc.Message,
) *SourceStage {
	i := &SourceStage{
		id:  initial,
		msg: msg,
		ch:  ch,
	}
	return i
}

func (s *SourceStage) Run(cfg *RunCfg) {
	for {
		select {
		case s.ch <- s.next():
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

// SinkStage defines the last output of the orchestration, where all messages
// are dropped.
type SinkStage struct {
	ch chan *State
}

func NewSinkOutput(ch chan *State) *SinkStage {
	s := &SinkStage{
		ch: ch,
	}
	return s
}

func (s *SinkStage) Run(cfg *RunCfg) {
	for {
		select {
		// Discard results
		case <-s.ch:
		case <-cfg.term:
			close(cfg.done)
			return
		}
	}
}
