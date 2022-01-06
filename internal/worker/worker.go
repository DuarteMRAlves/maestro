package worker

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/invoke"
	"github.com/DuarteMRAlves/maestro/internal/reflection"
	"google.golang.org/grpc"
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
}

func NewWorker(address string, rpc reflection.RPC) (Worker, error) {
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
		}
		return w, nil
	default:
		return nil, errdefs.InvalidArgumentWithMsg("unsupported rpc type")
	}
}

func (w *UnaryWorker) Run() {
}
