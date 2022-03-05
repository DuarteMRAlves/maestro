package orchestration

import (
	"context"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/storage"
	"github.com/DuarteMRAlves/maestro/test/protobuf/unit"
	"github.com/dgraph-io/badger/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"gotest.tools/v3/assert"
	"net"
	"sync"
	"sync/atomic"
	"testing"
)

func TestExecution_Linear(t *testing.T) {
	var (
		e   *Execution
		err error
	)

	logger, err := zap.NewDevelopment(zap.IncreaseLevel(zap.WarnLevel))
	assert.NilError(t, err, "init logger")

	sourceLis, err := net.Listen("tcp", "localhost:0")
	assert.NilError(t, err, "failed to listen source")
	sourceAddr := sourceLis.Addr().String()
	sourceServer := startLinearSource(t, sourceLis)
	defer sourceServer.Stop()

	transformLis, err := net.Listen("tcp", "localhost:0")
	assert.NilError(t, err, "failed to listen transform")
	transformAddr := transformLis.Addr().String()
	transformServer := startLinearTransform(t, transformLis)
	defer transformServer.Stop()

	max := 3
	collect := make([]*unit.LinearMessage, 0, max)
	done := make(chan struct{})

	sinkLis, err := net.Listen("tcp", "localhost:0")
	assert.NilError(t, err, "failed to listen sink")
	sinkAddr := sinkLis.Addr().String()
	sinkServer := startLinearSink(t, sinkLis, max, &collect, done)
	defer sinkServer.Stop()

	db := storage.NewTestDb(t)
	defer db.Close()

	err = initLinearDb(db, sourceAddr, transformAddr, sinkAddr)
	assert.NilError(t, err, "init db")

	err = db.View(
		func(txn *badger.Txn) error {
			builder := newBuilder(txn)
			e, err = builder.withOrchestration("default").withLogger(logger).build()
			return err
		},
	)
	assert.NilError(t, err, "build error")
	e.Start()
	<-done
	e.Stop()
	assert.Equal(t, 3, len(collect), "invalid length")
	for i, msg := range collect {
		val := int64(i)
		assert.Equal(t, msg.Value, (val+1)*2)
	}
}

type LinearSource struct {
	unit.UnimplementedLinearSourceServer
	counter int64
}

func (s *LinearSource) Process(
	_ context.Context,
	_ *emptypb.Empty,
) (*unit.LinearMessage, error) {
	msg := &unit.LinearMessage{Value: s.counter}
	atomic.AddInt64(&s.counter, 1)
	return msg, nil
}

type LinearTransform struct {
	unit.UnimplementedLinearTransformServer
}

func (_ *LinearTransform) Process(
	_ context.Context,
	msg *unit.LinearMessage,
) (*unit.LinearMessage, error) {
	msg.Value *= 2
	return msg, nil
}

type LinearSink struct {
	unit.UnimplementedLinearSinkServer
	max     int
	collect *[]*unit.LinearMessage
	done    chan<- struct{}
	mu      sync.Mutex
}

func (s *LinearSink) Process(
	_ context.Context,
	msg *unit.LinearMessage,
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
	unit.RegisterLinearSourceServer(s, &LinearSource{counter: 1})
	reflection.Register(s)

	go func() {
		err := s.Serve(lis)
		assert.NilError(t, err, "linear source server error")
	}()
	return s
}

func startLinearTransform(t *testing.T, lis net.Listener) *grpc.Server {
	s := grpc.NewServer()
	unit.RegisterLinearTransformServer(s, &LinearTransform{})
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
	collect *[]*unit.LinearMessage,
	done chan<- struct{},
) *grpc.Server {
	s := grpc.NewServer()
	unit.RegisterLinearSinkServer(
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

func initLinearDb(
	db *badger.DB,
	sourceAddr, transformAddr, sinkAddr string,
) error {
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
		Service:       "unit.LinearSource",
		Rpc:           "Process",
		Address:       sourceAddr,
		Orchestration: "default",
		Asset:         "",
	}

	transform := &api.Stage{
		Name:          "transform",
		Phase:         api.StagePending,
		Service:       "unit.LinearTransform",
		Rpc:           "Process",
		Address:       transformAddr,
		Orchestration: "default",
		Asset:         "",
	}

	sink := &api.Stage{
		Name:          "sink",
		Phase:         api.StagePending,
		Service:       "unit.LinearSink",
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

func TestExecution_SplitAndMerge(t *testing.T) {
	var (
		e   *Execution
		err error
	)

	logger, err := zap.NewDevelopment(zap.IncreaseLevel(zap.WarnLevel))
	assert.NilError(t, err, "init logger")

	sourceLis, err := net.Listen("tcp", "localhost:0")
	assert.NilError(t, err, "failed to listen source")
	sourceAddr := sourceLis.Addr().String()
	sourceServer := startSplitAndMergeSource(t, sourceLis)
	defer sourceServer.Stop()

	transformLis, err := net.Listen("tcp", "localhost:0")
	assert.NilError(t, err, "failed to listen transform")
	transformAddr := transformLis.Addr().String()
	transformServer := startSplitAndMergeTransform(t, transformLis)
	defer transformServer.Stop()

	max := 3
	collect := make([]*unit.SplitAndMergePair, 0, max)
	done := make(chan struct{})

	sinkLis, err := net.Listen("tcp", "localhost:0")
	assert.NilError(t, err, "failed to listen sink")
	sinkAddr := sinkLis.Addr().String()
	sinkServer := startSplitAndMergeSink(t, sinkLis, max, &collect, done)
	defer sinkServer.Stop()

	db := storage.NewTestDb(t)
	defer db.Close()

	err = initSplitAndMergeDb(db, sourceAddr, transformAddr, sinkAddr)
	assert.NilError(t, err, "init db")

	err = db.View(
		func(txn *badger.Txn) error {
			builder := newBuilder(txn)
			e, err = builder.withOrchestration("default").withLogger(logger).build()
			return err
		},
	)

	assert.NilError(t, err, "build error")
	e.Start()
	<-done
	e.Stop()
	assert.Equal(t, 3, len(collect), "invalid length")
	for i, msg := range collect {
		val := int64(i)
		assert.Equal(t, msg.Source.Value, val+1)
		assert.Equal(t, msg.Transformed.Value, (val+1)*2)
	}
}

type SplitAndMergeSource struct {
	unit.UnimplementedSplitAndMergeSourceServer
	counter int64
}

func (s *SplitAndMergeSource) Process(
	_ context.Context,
	_ *emptypb.Empty,
) (*unit.SplitAndMergeMessage, error) {
	msg := &unit.SplitAndMergeMessage{Value: s.counter}
	atomic.AddInt64(&s.counter, 1)
	return msg, nil
}

type SplitAndMergeTransform struct {
	unit.UnimplementedSplitAndMergeTransformServer
}

func (_ *SplitAndMergeTransform) Process(
	_ context.Context,
	msg *unit.SplitAndMergeMessage,
) (*unit.SplitAndMergeMessage, error) {
	msg.Value *= 2
	return msg, nil
}

type SplitAndMergeSink struct {
	unit.UnimplementedSplitAndMergeSinkServer
	max     int
	collect *[]*unit.SplitAndMergePair
	done    chan<- struct{}
	mu      sync.Mutex
}

func (s *SplitAndMergeSink) Process(
	_ context.Context,
	msg *unit.SplitAndMergePair,
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

func startSplitAndMergeSource(t *testing.T, lis net.Listener) *grpc.Server {
	s := grpc.NewServer()
	unit.RegisterSplitAndMergeSourceServer(s, &SplitAndMergeSource{counter: 1})
	reflection.Register(s)

	go func() {
		err := s.Serve(lis)
		assert.NilError(t, err, "split and merge source server error")
	}()
	return s
}

func startSplitAndMergeTransform(t *testing.T, lis net.Listener) *grpc.Server {
	s := grpc.NewServer()
	unit.RegisterSplitAndMergeTransformServer(s, &SplitAndMergeTransform{})
	reflection.Register(s)

	go func() {
		err := s.Serve(lis)
		assert.NilError(t, err, "split and merge transform server error")
	}()
	return s
}

func startSplitAndMergeSink(
	t *testing.T,
	lis net.Listener,
	max int,
	collect *[]*unit.SplitAndMergePair,
	done chan<- struct{},
) *grpc.Server {
	s := grpc.NewServer()
	unit.RegisterSplitAndMergeSinkServer(
		s,
		&SplitAndMergeSink{max: max, collect: collect, done: done},
	)
	reflection.Register(s)

	go func() {
		err := s.Serve(lis)
		assert.NilError(t, err, "split and merge sink server error")
	}()
	return s
}

func initSplitAndMergeDb(
	db *badger.DB,
	sourceAddr, transformAddr, sinkAddr string,
) error {
	var err error

	o := &api.Orchestration{
		Name:   "default",
		Phase:  api.OrchestrationPending,
		Stages: []api.StageName{"source", "transform", "sink"},
		Links: []api.LinkName{
			"link-source-transform",
			"link-transform-sink",
			"link-source-sink",
		},
	}

	source := &api.Stage{
		Name:          "source",
		Phase:         api.StagePending,
		Service:       "unit.SplitAndMergeSource",
		Rpc:           "Process",
		Address:       sourceAddr,
		Orchestration: "default",
		Asset:         "",
	}

	transform := &api.Stage{
		Name:          "transform",
		Phase:         api.StagePending,
		Service:       "unit.SplitAndMergeTransform",
		Rpc:           "Process",
		Address:       transformAddr,
		Orchestration: "default",
		Asset:         "",
	}

	sink := &api.Stage{
		Name:          "sink",
		Phase:         api.StagePending,
		Service:       "unit.SplitAndMergeSink",
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

	sourceToSink := &api.Link{
		Name:          "link-source-sink",
		SourceStage:   "source",
		SourceField:   "",
		TargetStage:   "sink",
		TargetField:   "source",
		Orchestration: "default",
	}

	transformToSink := &api.Link{
		Name:          "link-transform-sink",
		SourceStage:   "transform",
		SourceField:   "",
		TargetStage:   "sink",
		TargetField:   "transformed",
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
			err = helper.SaveLink(sourceToSink)
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
