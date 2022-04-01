package grpcexecute

import (
	"context"
	"errors"
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/internal/execute"
	igrpc "github.com/DuarteMRAlves/maestro/internal/grpc"
	"github.com/DuarteMRAlves/maestro/internal/logs"
	"github.com/DuarteMRAlves/maestro/internal/mock"
	"github.com/DuarteMRAlves/maestro/internal/retry"
	"github.com/DuarteMRAlves/maestro/test/protobuf/integration"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestLinear(t *testing.T) {
	var (
		source  linearSource
		transf  linearTransform
		sink    linearSink
		backoff retry.ExponentialBackoff
	)

	max := 100
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

	r := igrpc.NewReflectionMethodLoader(5*time.Minute, backoff, logs.New(true))
	executionBuilder := execute.NewBuilder(
		stageLoader, linkLoader, r, logs.New(true),
	)

	pipeline := internal.NewPipeline(
		createPipelineName(t, "pipeline"),
		[]internal.StageName{sourceName, transfName, sinkName},
		[]internal.LinkName{sourceToTransformName, transformToSinkName},
	)

	e, err := executionBuilder(pipeline)
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

	if diff := cmp.Diff(max, len(collect)); diff != "" {
		t.Fatalf("mismatch on number of collected messages:\n%s", diff)
	}

	prev := int64(0)
	for i, msg := range collect {
		if prev >= msg.Val {
			t.Fatalf("wrong value order at %d, %d: values are %d, %d", i-1, i, prev, msg.Val)
		}
		if msg.Val%3 != 0 {
			t.Fatalf("value %d is not divisible by 3: %d", i, msg.Val)
		}
		prev = msg.Val
	}
}

type linearSource struct {
	integration.LinearSourceServer
	counter int64
}

func (s *linearSource) Generate(
	_ context.Context, _ *emptypb.Empty,
) (*integration.LinearMessage, error) {
	delay := time.Duration(rand.Int63n(5))
	time.Sleep(delay * time.Millisecond)
	val := atomic.AddInt64(&s.counter, 1)
	return &integration.LinearMessage{Val: val}, nil
}

type linearTransform struct {
	integration.LinearTransformServer
}

func (t *linearTransform) Process(
	_ context.Context, req *integration.LinearMessage,
) (*integration.LinearMessage, error) {
	delay := time.Duration(rand.Int63n(5))
	time.Sleep(delay * time.Millisecond)
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
	delay := time.Duration(rand.Int63n(5))
	time.Sleep(delay * time.Millisecond)
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
