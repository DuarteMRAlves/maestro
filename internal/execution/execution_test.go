package execution

import (
	"context"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/rpc"
	"github.com/DuarteMRAlves/maestro/internal/storage"
	"github.com/DuarteMRAlves/maestro/internal/util"
	"github.com/DuarteMRAlves/maestro/tests/pb"
	"github.com/dgraph-io/badger/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"gotest.tools/v3/assert"
	"net"
	"sync"
	"sync/atomic"
	"testing"
)

func TestExecution_LinearPipeline(t *testing.T) {
	var (
		e   *Execution
		err error
	)

	sourceLis := util.NewTestListener(t)
	sourceAddr := sourceLis.Addr().String()
	sourceServer := startLinearSource(t, sourceLis)
	defer sourceServer.Stop()

	transformLis := util.NewTestListener(t)
	transformAddr := transformLis.Addr().String()
	transformServer := startLinearTransform(t, transformLis)
	defer transformServer.Stop()

	max := 3
	collect := make([]*pb.LinearMessage, 0, max)
	done := make(chan struct{})

	sinkLis := util.NewTestListener(t)
	sinkAddr := sinkLis.Addr().String()
	sinkServer := startLinearSink(t, sinkLis, max, &collect, done)
	defer sinkServer.Stop()

	rpcManager := rpc.NewManager()

	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	assert.NilError(t, err, "create db")
	defer db.Close()

	err = initDb(db, sourceAddr, transformAddr, sinkAddr)
	assert.NilError(t, err, "init db")

	err = db.View(
		func(txn *badger.Txn) error {
			builder := newBuilder(txn, rpcManager)
			e, err = builder.withOrchestration("default").build()
			return err
		},
	)
	assert.NilError(t, err, "build error")
	e.Start()
	<-done
	e.Stop()
	var prev int64 = 0
	for _, msg := range collect {
		assert.Assert(t, msg.Value >= prev)
		assert.Assert(t, msg.Value%2 == 0)
		prev = msg.Value
	}
}

type LinearSource struct {
	pb.UnimplementedLinearSourceServer
	counter int64
}

func (s *LinearSource) Process(
	_ context.Context,
	_ *emptypb.Empty,
) (*pb.LinearMessage, error) {
	msg := &pb.LinearMessage{Value: s.counter}
	atomic.AddInt64(&s.counter, 1)
	return msg, nil
}

type LinearTransform struct {
	pb.UnimplementedLinearTransformServer
}

func (_ *LinearTransform) Process(
	_ context.Context,
	msg *pb.LinearMessage,
) (*pb.LinearMessage, error) {
	msg.Value *= 2
	return msg, nil
}

type LinearSink struct {
	pb.UnimplementedLinearSinkServer
	max     int
	collect *[]*pb.LinearMessage
	done    chan<- struct{}
	mu      sync.Mutex
}

func (s *LinearSink) Process(
	_ context.Context,
	msg *pb.LinearMessage,
) (*emptypb.Empty, error) {

	s.mu.Lock()
	defer s.mu.Unlock()
	// Receive while not at full capacity
	if len(*s.collect) < s.max {
		*s.collect = append(*s.collect, msg)
	}
	// Notify when full. Remaining messages are discarded.
	if len(*s.collect) == s.max && s.done != nil {
		close(s.done)
		s.done = nil
	}
	return &emptypb.Empty{}, nil
}

func startLinearSource(t *testing.T, lis net.Listener) *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterLinearSourceServer(s, &LinearSource{counter: 1})
	reflection.Register(s)

	go func() {
		err := s.Serve(lis)
		assert.NilError(t, err, "linear source server error")
	}()
	return s
}

func startLinearTransform(t *testing.T, lis net.Listener) *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterLinearTransformServer(s, &LinearTransform{})
	reflection.Register(s)

	go func() {
		err := s.Serve(lis)
		assert.NilError(t, err, "linear transform server error")
	}()
	return s
}

func startLinearSink(
	t *testing.T,
	lis net.Listener,
	max int,
	collect *[]*pb.LinearMessage,
	done chan<- struct{},
) *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterLinearSinkServer(
		s,
		&LinearSink{max: max, collect: collect, done: done},
	)
	reflection.Register(s)

	go func() {
		err := s.Serve(lis)
		assert.NilError(t, err, "linear sink server error")
	}()
	return s
}

func initDb(db *badger.DB, sourceAddr, transformAddr, sinkAddr string) error {
	var err error

	o := &api.Orchestration{
		Name:   "default",
		Phase:  api.OrchestrationPending,
		Stages: []api.StageName{"source", "transform", "sink"},
		Links:  []api.LinkName{"link-source-transform", "link-transform-sink"},
	}

	source := &api.Stage{
		Name:          "source",
		Phase:         api.StagePending,
		Service:       "pb.LinearSource",
		Rpc:           "Process",
		Address:       sourceAddr,
		Orchestration: "default",
		Asset:         "",
	}

	transform := &api.Stage{
		Name:          "transform",
		Phase:         api.StagePending,
		Service:       "pb.LinearTransform",
		Rpc:           "Process",
		Address:       transformAddr,
		Orchestration: "default",
		Asset:         "",
	}

	sink := &api.Stage{
		Name:          "sink",
		Phase:         api.StagePending,
		Service:       "pb.LinearSink",
		Rpc:           "Process",
		Address:       sinkAddr,
		Orchestration: "default",
		Asset:         "",
	}

	sourceToTransform := &api.Link{
		Name:          "link-source-transform",
		SourceStage:   "source",
		SourceField:   "",
		TargetStage:   "transform",
		TargetField:   "",
		Orchestration: "default",
	}

	transformToSink := &api.Link{
		Name:          "link-transform-sink",
		SourceStage:   "transform",
		SourceField:   "",
		TargetStage:   "sink",
		TargetField:   "",
		Orchestration: "default",
	}

	return db.Update(
		func(txn *badger.Txn) error {
			helper := storage.NewTxnHelper(txn)
			err = helper.SaveOrchestration(o)
			if err != nil {
				return err
			}
			err = helper.SaveStage(source)
			if err != nil {
				return err
			}
			err = helper.SaveStage(transform)
			if err != nil {
				return err
			}
			err = helper.SaveStage(sink)
			if err != nil {
				return err
			}
			err = helper.SaveLink(sourceToTransform)
			if err != nil {
				return err
			}
			err = helper.SaveLink(transformToSink)
			if err != nil {
				return err
			}
			return nil
		},
	)
}
