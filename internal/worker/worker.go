package worker

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/flow"
	"github.com/DuarteMRAlves/maestro/internal/invoke"
	"github.com/DuarteMRAlves/maestro/internal/reflection"
	"google.golang.org/grpc"
)

type Worker interface {
	Run()
}

type Cfg struct {
	Address string
	Rpc     reflection.RPC
	Input   flow.Input
	Output  flow.Output
	Done    chan<- bool
	MaxMsg  int
}

func (c *Cfg) Clone() *Cfg {
	return &Cfg{
		Address: c.Address,
		Rpc:     c.Rpc,
		Input:   c.Input,
		Output:  c.Output,
		Done:    c.Done,
		MaxMsg:  c.MaxMsg,
	}
}

func NewWorker(cfg *Cfg) (Worker, error) {
	cfg = cfg.Clone()
	switch {
	case cfg.Rpc.IsUnary():
		conn, err := grpc.Dial(cfg.Address, grpc.WithInsecure())
		if err != nil {
			return nil, errdefs.InvalidArgumentWithMsg(
				"unable to connect to address: %s",
				cfg.Address)
		}
		w := &UnaryWorker{
			Address: cfg.Address,
			conn:    conn,
			rpc:     cfg.Rpc,
			invoker: invoke.NewUnary(cfg.Rpc.FullyQualifiedName(), conn),
			input:   cfg.Input,
			output:  cfg.Output,
			done:    cfg.Done,
			maxMsg:  cfg.MaxMsg,
		}
		return w, nil
	default:
		return nil, errdefs.InvalidArgumentWithMsg("unsupported rpc type")
	}
}
