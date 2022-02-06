package execution

import (
	"context"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/rpc"
	"google.golang.org/grpc"
	"io"
	"time"
)

type Worker interface {
	Run(*RunCfg)
}

// RunCfg specifies the configuration that the Worker should use when running.
type RunCfg struct {
	// term is a channel that will be signaled if the Worker should stop.
	term <-chan struct{}
	// errs is a channel that the worker should use to send errors in order to
	// be processed.
	// The io.EOF error should not be sent through this channel at is just a
	// termination signal
	errs chan<- error
}

// UnaryWorker manages the execution of a stage in a pipeline
type UnaryWorker struct {
	Address string
	conn    grpc.ClientConnInterface
	rpc     rpc.RPC
	invoker rpc.UnaryClient

	input  Input
	output Output

	done chan<- bool
}

func (w *UnaryWorker) Run(cfg *RunCfg) {
	var (
		in, out  *State
		req, rep interface{}
		err      error
	)

	for {
		select {
		case in = <-w.input.Chan():
			if in.Err() == io.EOF {
				w.done <- true
				return
			}
			if in.Err() != nil {
				cfg.errs <- in.Err()
				continue
			}
			req = in.Msg()
			rep = w.rpc.Output().NewEmpty()

			err = w.invoke(req, rep)
			if err != nil {
				cfg.errs <- err
				continue
			}

			out = NewState(in.Id(), rep)
			w.output.Chan() <- out
		case <-cfg.term:
			w.done <- true
			return
		}
	}
}

func (w *UnaryWorker) invoke(req interface{}, rep interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return w.invoker.Invoke(ctx, req, rep)
}

type WorkerCfg struct {
	Address string
	Rpc     rpc.RPC
	Input   Input
	Output  Output
	Done    chan<- bool
}

func (c *WorkerCfg) Clone() *WorkerCfg {
	return &WorkerCfg{
		Address: c.Address,
		Rpc:     c.Rpc,
		Input:   c.Input,
		Output:  c.Output,
		Done:    c.Done,
	}
}

func NewWorker(cfg *WorkerCfg) (Worker, error) {
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
		w := &UnaryWorker{
			Address: cfg.Address,
			conn:    conn,
			rpc:     cfg.Rpc,
			invoker: rpc.NewUnary(cfg.Rpc.FullyQualifiedName(), conn),
			input:   cfg.Input,
			output:  cfg.Output,
			done:    cfg.Done,
		}
		return w, nil
	default:
		return nil, errdefs.InvalidArgumentWithMsg("unsupported rpc type")
	}
}
