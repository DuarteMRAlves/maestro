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
	Run()
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

func (w *UnaryWorker) Run() {
	var (
		in, out  *State
		req, rep interface{}
		err      error
	)

	for {
		in, err = w.input.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		req = in.Msg()
		rep = w.rpc.Output().NewEmpty()

		err = w.invoke(req, rep)
		if err != nil {
			panic(err)
		}

		out = NewState(in.Id(), rep)
		w.output.Yield(out)
	}
	w.done <- true
}

func (w *UnaryWorker) invoke(req interface{}, rep interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return w.invoker.Invoke(ctx, req, rep)
}
