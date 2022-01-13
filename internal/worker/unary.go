package worker

import (
	"context"
	"github.com/DuarteMRAlves/maestro/internal/flow"
	"github.com/DuarteMRAlves/maestro/internal/invoke"
	"github.com/DuarteMRAlves/maestro/internal/reflection"
	"google.golang.org/grpc"
	"time"
)

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
