package worker

import (
	"context"
	flowinput "github.com/DuarteMRAlves/maestro/internal/flow/input"
	flowoutput "github.com/DuarteMRAlves/maestro/internal/flow/output"
	"github.com/DuarteMRAlves/maestro/internal/flow/state"
	"github.com/DuarteMRAlves/maestro/internal/invoke"
	"github.com/DuarteMRAlves/maestro/internal/reflection"
	"google.golang.org/grpc"
	"io"
	"time"
)

// UnaryWorker manages the execution of a stage in a pipeline
type UnaryWorker struct {
	Address string
	conn    grpc.ClientConnInterface
	rpc     reflection.RPC
	invoker invoke.UnaryClient

	input  flowinput.Input
	output flowoutput.Output

	done   chan<- bool
	maxMsg int
}

func (w *UnaryWorker) Run() {
	var (
		in, out  *state.State
		req, rep interface{}
		err      error
	)

	for msgCount := 0; msgCount < w.maxMsg; msgCount++ {
		in, err = w.input.Next()
		if err != nil {
			if err == io.EOF {
				return
			}
			panic(err)
		}

		req = in.Msg()
		rep = w.rpc.Output().NewEmpty()

		err = w.invoke(req, rep)
		if err != nil {
			panic(err)
		}

		out = state.New(in.Id(), rep)
		w.output.Yield(out)
	}
	w.done <- true
}

func (w *UnaryWorker) invoke(req interface{}, rep interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return w.invoker.Invoke(ctx, req, rep)
}
