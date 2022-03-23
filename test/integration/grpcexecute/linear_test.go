package grpcexecute

import (
	"context"
	"errors"
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/internal/execute"
	igrpc "github.com/DuarteMRAlves/maestro/internal/grpc"
	"github.com/DuarteMRAlves/maestro/internal/mock"
	"github.com/DuarteMRAlves/maestro/test/protobuf/integration"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"sync"
	"sync/atomic"
	"testing"
)

func TestLinear(t *testing.T) {
	var (
		source linearSource
		transf linearTransform
		sink   linearSink
	)

	max := 3
	collect := make([]*integration.LinearMessage, 0, max)
	done := make(chan struct{})

	sink.max = max
	sink.collect = &collect
	sink.done = done

	sourceName := createStageName(t, "source")
	transfName := createStageName(t, "transform")
	sinkName := createStageName(t, "sink")

	sourceAddr, sourceStart, sourceStop := createGrpcServer(
		t,
		func(registrar grpc.ServiceRegistrar) {
			integration.RegisterLinearSourceServer(registrar, &source)
		},
	)
	defer sourceStop()
	go sourceStart()

	transfAddr, transfStart, transfStop := createGrpcServer(
		t,
		func(registrar grpc.ServiceRegistrar) {
			integration.RegisterLinearTransformServer(registrar, &transf)
		},
	)
	defer transfStop()
	go transfStart()

	sinkAddr, sinkStart, sinkStop := createGrpcServer(
		t,
		func(registrar grpc.ServiceRegistrar) {
			integration.RegisterLinearSinkServer(registrar, &sink)
		},
	)
	defer sinkStop()
	go sinkStart()

	sourceCtx := createMethodContext(internal.NewAddress(sourceAddr.String()))
	transfCtx := createMethodContext(internal.NewAddress(transfAddr.String()))
	sinkCtx := createMethodContext(internal.NewAddress(sinkAddr.String()))

	sourceStage := internal.NewStage(sourceName, sourceCtx)
	transfStage := internal.NewStage(transfName, transfCtx)
	sinkStage := internal.NewStage(sinkName, sinkCtx)

	stages := map[internal.StageName]internal.Stage{
		sourceName: sourceStage,
		transfName: transfStage,
		sinkName:   sinkStage,
	}
	stageLoader := &mock.StageStorage{Stages: stages}

	sourceToTransformName := createLinkName(t, "link-source-transform")
	sourceToTransform := internal.NewLink(
		sourceToTransformName,
		internal.NewLinkEndpoint(sourceName, internal.MessageField{}),
		internal.NewLinkEndpoint(transfName, internal.MessageField{}),
	)

	transformToSinkName := createLinkName(t, "link-transform-sink")
	transformToSink := internal.NewLink(
		transformToSinkName,
		internal.NewLinkEndpoint(transfName, internal.MessageField{}),
		internal.NewLinkEndpoint(sinkName, internal.MessageField{}),
	)

	links := map[internal.LinkName]internal.Link{
		sourceToTransformName: sourceToTransform,
		transformToSinkName:   transformToSink,
	}
	linkLoader := &mock.LinkStorage{Links: links}

	executionBuilder := execute.NewBuilder(stageLoader, linkLoader, igrpc.ReflectionMethodLoader)

	orchestration := internal.NewOrchestration(
		createOrchName(t, "orchestration"),
		[]internal.StageName{sourceName, transfName, sinkName},
		[]internal.LinkName{sourceToTransformName, transformToSinkName},
	)

	e, err := executionBuilder(orchestration)
	if err != nil {
		t.Fatalf("build error: %s", err)
	}

	e.Start()
	<-done
	if err := e.Stop(); err != nil {
		cause, ok := errors.Unwrap(err).(interface {
			GRPCStatus() *status.Status
		})
		if !ok {
			t.Fatalf("stop error does not implement grpc interface")
		}
		st := cause.GRPCStatus()
		// The cancel can happen midways through a method call
		if diff := cmp.Diff(codes.Canceled, st.Code()); diff != "" {
			t.Fatalf("stop error code mismatch:\n%s", diff)
		}
	}
	expected := []*integration.LinearMessage{{Val: 3}, {Val: 6}, {Val: 9}}

	if diff := cmp.Diff(len(expected), len(collect)); diff != "" {
		t.Fatalf("mismatch on number of collected messages:\n%s", diff)
	}

	cmpOpts := cmpopts.IgnoreUnexported(integration.LinearMessage{})
	if diff := cmp.Diff(expected, collect, cmpOpts); diff != "" {
		t.Fatalf("mismatch on collected messages:\n%s", diff)
	}
}

func createGrpcServer(
	t *testing.T, registerSrv func(grpc.ServiceRegistrar),
) (net.Addr, func(), func()) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen: %s", err)
	}
	s := grpc.NewServer()
	registerSrv(s)
	reflection.Register(s)

	start := func() {
		if err := s.Serve(lis); err != nil {
			t.Fatalf("Failed to server: %s", err)
		}
	}
	stop := func() {
		s.Stop()
	}
	return lis.Addr(), start, stop
}

type linearSource struct {
	integration.LinearSourceServer
	counter int64
}

func (s *linearSource) Generate(
	_ context.Context, _ *emptypb.Empty,
) (*integration.LinearMessage, error) {
	val := atomic.AddInt64(&s.counter, 1)
	return &integration.LinearMessage{Val: val}, nil
}

type linearTransform struct {
	integration.LinearTransformServer
}

func (t *linearTransform) Process(
	_ context.Context, req *integration.LinearMessage,
) (*integration.LinearMessage, error) {
	return &integration.LinearMessage{Val: 3 * req.Val}, nil
}

type linearSink struct {
	integration.LinearSinkServer
	max     int
	collect *[]*integration.LinearMessage
	done    chan<- struct{}
	mu      sync.Mutex
}

func (s *linearSink) Collect(
	_ context.Context, req *integration.LinearMessage,
) (*emptypb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Receive while not at full capacity
	if len(*s.collect) < s.max {
		*s.collect = append(*s.collect, req)
	}
	// Notify when full. Remaining messages are discarded.
	if len(*s.collect) == s.max && s.done != nil {
		close(s.done)
		s.done = nil
	}
	return &emptypb.Empty{}, nil
}

func createOrchName(t *testing.T, name string) internal.OrchestrationName {
	orchName, err := internal.NewOrchestrationName(name)
	if err != nil {
		t.Fatalf("create orchestration name %s: %s", name, err)
	}
	return orchName
}

func createStageName(t *testing.T, name string) internal.StageName {
	stageName, err := internal.NewStageName(name)
	if err != nil {
		t.Fatalf("create stage name %s: %s", name, err)
	}
	return stageName
}

func createLinkName(t *testing.T, name string) internal.LinkName {
	linkName, err := internal.NewLinkName(name)
	if err != nil {
		t.Fatalf("create link name %s: %s", name, err)
	}
	return linkName
}

func createMethodContext(addr internal.Address) internal.MethodContext {
	var (
		emptyService internal.Service
		emptyMethod  internal.Method
	)
	return internal.NewMethodContext(addr, emptyService, emptyMethod)
}
