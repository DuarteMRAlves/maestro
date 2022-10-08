package grpcexecute

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/compiled"
	"github.com/DuarteMRAlves/maestro/internal/execute"
	"github.com/DuarteMRAlves/maestro/internal/grpcw"
	"github.com/DuarteMRAlves/maestro/internal/logs"
	"github.com/DuarteMRAlves/maestro/internal/retry"
	"github.com/DuarteMRAlves/maestro/test/protobuf/integration"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestSplitAndMerge(t *testing.T) {
	var (
		source  splitAndMergeSource
		transf  splitAndMergeTransform
		sink    splitAndMergeSink
		backoff retry.ExponentialBackoff
	)

	max := 100
	collect := make([]*integration.JoinMessage, 0, max)
	done := make(chan struct{})

	sink.max = max
	sink.collect = &collect
	sink.done = done

	sourceAddr, sourceStart, sourceStop := createGrpcServer(
		t,
		func(registrar grpc.ServiceRegistrar) {
			integration.RegisterSplitAndMergeSourceServer(registrar, &source)
		},
	)
	defer sourceStop()
	go sourceStart()

	transfAddr, transfStart, transfStop := createGrpcServer(
		t,
		func(registrar grpc.ServiceRegistrar) {
			integration.RegisterSplitAndMergeTransformServer(registrar, &transf)
		},
	)
	defer transfStop()
	go transfStart()

	sinkAddr, sinkStart, sinkStop := createGrpcServer(
		t,
		func(registrar grpc.ServiceRegistrar) {
			integration.RegisterSplitAndMergeSinkServer(registrar, &sink)
		},
	)
	defer sinkStop()
	go sinkStart()

	cfg := &api.Pipeline{
		Name: "pipeline",
		Stages: []*api.Stage{
			{Name: "source", Address: sourceAddr.String()},
			{Name: "transform", Address: transfAddr.String()},
			{Name: "sink", Address: sinkAddr.String()},
		},
		Links: []*api.Link{
			{
				Name:        "link-source-transform",
				SourceStage: "source",
				TargetStage: "transform",
			},
			{
				Name:        "link-source-sink",
				SourceStage: "source",
				TargetStage: "sink",
				TargetField: "original",
			},
			{
				Name:        "link-transform-sink",
				SourceStage: "transform",
				TargetStage: "sink",
				TargetField: "transformed",
			},
		},
	}

	resolver, err := grpcw.NewReflectionResolver(5*time.Minute, backoff, logs.New(true))
	if err != nil {
		t.Fatalf("create resolver: %v", err)
	}
	compilationCtx := compiled.NewContext(resolver)
	pipeline, err := compiled.New(compilationCtx, cfg)
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
		orig := msg.Original
		trans := msg.Transformed
		if diff := cmp.Diff(int64(i+1), orig.Val); diff != "" {
			t.Fatalf("mismatch at orig %d:\n%s", i, diff)
		}
		if diff := cmp.Diff(int64((i+1)*3), trans.Val); diff != "" {
			t.Fatalf("mismatch at trans %d:\n%s", i, diff)
		}
	}
}

type splitAndMergeSource struct {
	integration.SplitAndMergeSourceServer
	counter int64
}

func (s *splitAndMergeSource) Generate(
	_ context.Context, _ *emptypb.Empty,
) (*integration.SplitAndMergeMessage, error) {
	delay := time.Duration(rand.Int63n(5))
	time.Sleep(delay * time.Millisecond)
	val := atomic.AddInt64(&s.counter, 1)
	return &integration.SplitAndMergeMessage{Val: val}, nil
}

type splitAndMergeTransform struct {
	integration.SplitAndMergeTransformServer
}

func (t *splitAndMergeTransform) Process(
	_ context.Context, req *integration.SplitAndMergeMessage,
) (*integration.SplitAndMergeMessage, error) {
	delay := time.Duration(rand.Int63n(5))
	time.Sleep(delay * time.Millisecond)
	return &integration.SplitAndMergeMessage{Val: 3 * req.Val}, nil
}

type splitAndMergeSink struct {
	integration.SplitAndMergeSinkServer
	max     int
	collect *[]*integration.JoinMessage
	done    chan<- struct{}
	mu      sync.Mutex
}

func (s *splitAndMergeSink) Collect(
	_ context.Context, req *integration.JoinMessage,
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
