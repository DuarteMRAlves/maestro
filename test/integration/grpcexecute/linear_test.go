package grpcexecute

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/DuarteMRAlves/maestro/internal/compiled"
	"github.com/DuarteMRAlves/maestro/internal/execute"
	igrpc "github.com/DuarteMRAlves/maestro/internal/grpc"
	"github.com/DuarteMRAlves/maestro/internal/logs"
	"github.com/DuarteMRAlves/maestro/internal/retry"
	"github.com/DuarteMRAlves/maestro/internal/spec"
	"github.com/DuarteMRAlves/maestro/test/protobuf/integration"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestOfflineLinear(t *testing.T) {
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

	pipelineSpec := &spec.Pipeline{
		Name: "pipeline",
		Mode: spec.OfflineExecution,
		Stages: []*spec.Stage{
			{
				Name:          "source",
				MethodContext: spec.MethodContext{Address: sourceAddr.String()},
			},
			{
				Name:          "transform",
				MethodContext: spec.MethodContext{Address: transfAddr.String()},
			},
			{
				Name:          "sink",
				MethodContext: spec.MethodContext{Address: sinkAddr.String()},
			},
		},
		Links: []*spec.Link{
			{
				Name:        "link-source-transform",
				SourceStage: "source",
				TargetStage: "transform",
			},
			{
				Name:        "link-transform-sink",
				SourceStage: "transform",
				TargetStage: "sink",
			},
		},
	}

	methodLoader := igrpc.NewReflectionMethodLoader(5*time.Minute, backoff, logs.New(true))
	compilationCtx := compiled.NewContext(methodLoader)
	pipeline, err := compiled.New(compilationCtx, pipelineSpec)
	if err != nil {
		t.Fatalf("compile error: %s", err)
	}

	executionBuilder := execute.NewBuilder(logs.New(true))
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

	for i, msg := range collect {
		if diff := cmp.Diff(int64((i+1)*3), msg.Val); diff != "" {
			t.Fatalf("mismatch at msg %d:\n%s", i, diff)
		}
	}
}

func TestOnlineLinear(t *testing.T) {
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

	pipelineSpec := &spec.Pipeline{
		Name: "pipeline",
		Mode: spec.OnlineExecution,
		Stages: []*spec.Stage{
			{
				Name:          "source",
				MethodContext: spec.MethodContext{Address: sourceAddr.String()},
			},
			{
				Name:          "transform",
				MethodContext: spec.MethodContext{Address: transfAddr.String()},
			},
			{
				Name:          "sink",
				MethodContext: spec.MethodContext{Address: sinkAddr.String()},
			},
		},
		Links: []*spec.Link{
			{
				Name:        "link-source-transform",
				SourceStage: "source",
				TargetStage: "transform",
			},
			{
				Name:        "link-transform-sink",
				SourceStage: "transform",
				TargetStage: "sink",
			},
		},
	}

	methodLoader := igrpc.NewReflectionMethodLoader(5*time.Minute, backoff, logs.New(true))
	compilationCtx := compiled.NewContext(methodLoader)
	pipeline, err := compiled.New(compilationCtx, pipelineSpec)
	if err != nil {
		t.Fatalf("compile error: %s", err)
	}

	executionBuilder := execute.NewBuilder(logs.New(true))
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
