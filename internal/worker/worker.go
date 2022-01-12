package worker

import (
	"context"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/flow"
	"github.com/DuarteMRAlves/maestro/internal/invoke"
	"github.com/DuarteMRAlves/maestro/internal/reflection"
	"google.golang.org/grpc"
	"time"
)

type Worker interface {
	Run()
}

// UnaryWorker manages the execution of a stage in a pipeline
type UnaryWorker struct {
	Address string
	conn    grpc.ClientConnInterface
	rpc     reflection.RPC
	invoker invoke.UnaryClient

	input  flow.Input
	output flow.Output

	done   chan<- bool
	maxMsg int
}

func NewWorker(
	address string,
	rpc reflection.RPC,
	input flow.Input,
	output flow.Output,
	done chan<- bool,
	maxMsg int,
) (Worker, error) {
	switch {
	case rpc.IsUnary():
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			return nil, errdefs.InvalidArgumentWithMsg(
				"unable to connect to address: %s",
				address)
		}
		w := &UnaryWorker{
			Address: address,
			conn:    conn,
			rpc:     rpc,
			invoker: invoke.NewUnary(rpc.FullyQualifiedName(), conn),
			input:   input,
			output:  output,
			done:    done,
			maxMsg:  maxMsg,
		}
		return w, nil
	default:
		return nil, errdefs.InvalidArgumentWithMsg("unsupported rpc type")
	}
}

func (w *UnaryWorker) Run() {
	var (
		in, out  *flow.State
		req, rep interface{}
	)

	for msgCount := 0; msgCount < w.maxMsg; msgCount++ {
		in = w.input.Next()

		req = in.Msg()
		rep = w.rpc.Output().NewEmpty()

		err := w.invoke(req, rep)
		if err != nil {
			panic(err)
		}

		out = flow.New(in.Id(), rep)
		w.output.Yield(out)
	}
	w.done <- true
}

func (w *UnaryWorker) invoke(req interface{}, rep interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return w.invoker.Invoke(ctx, req, rep)
}
