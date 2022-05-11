package grpcexecute

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

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
)

func TestOfflineSplitAndMerge(t *testing.T) {
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

	sourceName := createStageName(t, "source")
	transfName := createStageName(t, "transform")
	sinkName := createStageName(t, "sink")

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

	sourceToSinkName := createLinkName(t, "link-source-sink")
	sourceToSink := internal.NewLink(
		sourceToSinkName,
		internal.NewLinkEndpoint(sourceName, internal.MessageField{}),
		internal.NewLinkEndpoint(sinkName, internal.NewMessageField("original")),
	)

	transformToSinkName := createLinkName(t, "link-transform-sink")
	transformToSink := internal.NewLink(
		transformToSinkName,
		internal.NewLinkEndpoint(transfName, internal.MessageField{}),
		internal.NewLinkEndpoint(sinkName, internal.NewMessageField("transformed")),
	)

	links := map[internal.LinkName]internal.Link{
		sourceToTransformName: sourceToTransform,
		sourceToSinkName:      sourceToSink,
		transformToSinkName:   transformToSink,
	}
	linkLoader := &mock.LinkStorage{Links: links}

	r := igrpc.NewReflectionMethodLoader(5*time.Minute, backoff, logs.New(true))
	executionBuilder := execute.NewBuilder(
		stageLoader, linkLoader, r, logs.New(true),
	)

	pipeline := internal.NewPipeline(
		createPipelineName(t, "pipeline"),
		internal.WithStages(sourceName, transfName, sinkName),
		internal.WithLinks(sourceToTransformName, sourceToSinkName, transformToSinkName),
		internal.WithOfflineExec(),
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

func TestOnlineSplitAndMerge(t *testing.T) {
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

	sourceName := createStageName(t, "source")
	transfName := createStageName(t, "transform")
	sinkName := createStageName(t, "sink")

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

	sourceToSinkName := createLinkName(t, "link-source-sink")
	sourceToSink := internal.NewLink(
		sourceToSinkName,
		internal.NewLinkEndpoint(sourceName, internal.MessageField{}),
		internal.NewLinkEndpoint(sinkName, internal.NewMessageField("original")),
	)

	transformToSinkName := createLinkName(t, "link-transform-sink")
	transformToSink := internal.NewLink(
		transformToSinkName,
		internal.NewLinkEndpoint(transfName, internal.MessageField{}),
		internal.NewLinkEndpoint(sinkName, internal.NewMessageField("transformed")),
	)

	links := map[internal.LinkName]internal.Link{
		sourceToTransformName: sourceToTransform,
		sourceToSinkName:      sourceToSink,
		transformToSinkName:   transformToSink,
	}
	linkLoader := &mock.LinkStorage{Links: links}

	r := igrpc.NewReflectionMethodLoader(5*time.Minute, backoff, logs.New(true))
	executionBuilder := execute.NewBuilder(
		stageLoader, linkLoader, r, logs.New(true),
	)

	pipeline := internal.NewPipeline(
		createPipelineName(t, "pipeline"),
		internal.WithStages(sourceName, transfName, sinkName),
		internal.WithLinks(sourceToTransformName, sourceToSinkName, transformToSinkName),
		internal.WithOnlineExec(),
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
		orig := msg.Original
		trans := msg.Transformed
		if prev >= orig.Val {
			t.Fatalf("wrong value order at %d, %d: values are %d, %d", i-1, i, prev, orig.Val)
		}
		if trans.Val != 3*orig.Val {
			t.Fatalf("transformed != 3 * original at %d: orig is %d and transf is %d", i, orig.Val, trans.Val)
		}
		prev = orig.Val
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
