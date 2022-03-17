package e2e

import (
	"bytes"
	"context"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/logs"
	"github.com/DuarteMRAlves/maestro/internal/server"
	"github.com/DuarteMRAlves/maestro/internal/storage"
	"github.com/DuarteMRAlves/maestro/old/internal/cli/maestroctl/cmd/create"
	"github.com/DuarteMRAlves/maestro/old/internal/cli/maestroctl/cmd/start"
	"github.com/DuarteMRAlves/maestro/test/protobuf"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"gotest.tools/v3/assert"
	"io/ioutil"
	"net"
	"sync"
	"sync/atomic"
	"testing"
)

const (
	sourcePort    = 50052
	transformPort = 50053
	sinkPort      = 50054
)

func TestE2E_Linear(t *testing.T) {
	sourceServer := startLinearSource(t)
	defer sourceServer.Stop()

	transformServer := startLinearTransform(t)
	defer transformServer.Stop()

	max := 3
	collect := make([]*protobuf.LinearMessage, 0, max)
	done := make(chan struct{})

	sinkServer := startLinearSink(t, max, &collect, done)
	defer sinkServer.Stop()

	maestroLis, err := net.Listen("tcp", "localhost:0")
	assert.NilError(t, err, "maestro server failed to listen")
	maestroAddr := maestroLis.Addr().String()

	s := buildMaestro(t)
	defer s.StopGrpc()

	go func() {
		err := s.ServeGrpc(maestroLis)
		assert.NilError(t, err, "serve error")
	}()

	runCreateCmd(t, maestroAddr)
	runStartCmd(t, maestroAddr)

	<-done

	assert.Equal(t, 3, len(collect), "invalid number of collected messages")
	for i, msg := range collect {
		val := int64(i)
		assert.Equal(t, msg.Value, (val+1)*2)
	}
}

func runCreateCmd(t *testing.T, maestroAddr string) {
	buffer := bytes.NewBufferString("")
	resourcesFile := "../data/e2e/linear_pipeline.yaml"

	createCmd := create.NewCmdCreate()
	createCmd.SetArgs([]string{"-f", resourcesFile, "--maestro", maestroAddr})
	createCmd.SetOut(buffer)
	err := createCmd.Execute()
	assert.NilError(t, err, "execute create error")
	out, err := ioutil.ReadAll(buffer)
	assert.NilError(t, err, "create read output error")
	assert.Equal(t, "", string(out), "create output differs")
}

func runStartCmd(t *testing.T, maestroAddr string) {
	buffer := bytes.NewBufferString("")

	startCmd := start.NewCmdStart()
	startCmd.SetArgs([]string{"LinearOrchestration", "--maestro", maestroAddr})
	startCmd.SetOut(buffer)
	err := startCmd.Execute()
	assert.NilError(t, err, "execute start error")
	out, err := ioutil.ReadAll(buffer)
	assert.NilError(t, err, "execute read output error")
	assert.Equal(t, "", string(out), "execute output differs")
}

type LinearSource struct {
	protobuf.UnimplementedLinearSourceServer
	counter int64
}

func (s *LinearSource) Process(
	_ context.Context,
	_ *emptypb.Empty,
) (*protobuf.LinearMessage, error) {
	msg := &protobuf.LinearMessage{Value: s.counter}
	atomic.AddInt64(&s.counter, 1)
	return msg, nil
}

type LinearTransform struct {
	protobuf.UnimplementedLinearTransformServer
}

func (_ *LinearTransform) Process(
	_ context.Context,
	msg *protobuf.LinearMessage,
) (*protobuf.LinearMessage, error) {
	msg.Value *= 2
	return msg, nil
}

type LinearSink struct {
	protobuf.UnimplementedLinearSinkServer
	max     int
	collect *[]*protobuf.LinearMessage
	done    chan<- struct{}
	mu      sync.Mutex
}

func (s *LinearSink) Process(
	_ context.Context,
	msg *protobuf.LinearMessage,
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

func startLinearSource(t *testing.T) *grpc.Server {
	addr := fmt.Sprintf("localhost:%d", sourcePort)
	lis, err := net.Listen("tcp", addr)
	assert.NilError(t, err, "failed to listen source")
	s := grpc.NewServer()
	protobuf.RegisterLinearSourceServer(s, &LinearSource{counter: 1})
	reflection.Register(s)

	go func() {
		err := s.Serve(lis)
		assert.NilError(t, err, "linear source server error")
	}()
	return s
}

func startLinearTransform(t *testing.T) *grpc.Server {
	addr := fmt.Sprintf("localhost:%d", transformPort)
	lis, err := net.Listen("tcp", addr)
	assert.NilError(t, err, "failed to listen transform")

	s := grpc.NewServer()
	protobuf.RegisterLinearTransformServer(s, &LinearTransform{})
	reflection.Register(s)

	go func() {
		err := s.Serve(lis)
		assert.NilError(t, err, "linear transform server error")
	}()
	return s
}

func startLinearSink(
	t *testing.T,
	max int,
	collect *[]*protobuf.LinearMessage,
	done chan<- struct{},
) *grpc.Server {
	addr := fmt.Sprintf("localhost:%d", sinkPort)
	lis, err := net.Listen("tcp", addr)
	assert.NilError(t, err, "failed to listen sink")

	s := grpc.NewServer()
	protobuf.RegisterLinearSinkServer(
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

func buildMaestro(t *testing.T) *server.Server {
	logger, err := logs.DefaultProductionLogger(zap.NewAtomicLevelAt(zap.InfoLevel))
	assert.NilError(t, err, "init logger")

	db, err := storage.NewDb()
	assert.NilError(t, err, "init db")

	s, err := server.NewBuilder().
		WithGrpc().
		WithLogger(logger).
		WithDb(db).
		Build()
	assert.NilError(t, err, "build error")
	return s
}
