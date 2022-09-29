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

func TestOfflineCycle(t *testing.T) {
	var (
		counter cycleCounter
		sum     cycleSum
		inc     cycleInc
		backoff retry.ExponentialBackoff
	)

	max := 100
	collect := make([]*integration.CycleSumMessage, 0, max)
	done := make(chan struct{})

	sum.max = max
	sum.collect = &collect
	sum.done = done

	counterAddr, counterStart, counterStop := createGrpcServer(
		t,
		func(registrar grpc.ServiceRegistrar) {
			integration.RegisterCycleCounterServer(registrar, &counter)
		},
	)
	defer counterStop()
	go counterStart()

	sumAddr, sumStart, sumStop := createGrpcServer(
		t,
		func(registrar grpc.ServiceRegistrar) {
			integration.RegisterCycleSumServer(registrar, &sum)
		},
	)
	defer sumStop()
	go sumStart()

	incAddr, incStart, incStop := createGrpcServer(
		t,
		func(registrar grpc.ServiceRegistrar) {
			integration.RegisterCycleIncServer(registrar, &inc)
		},
	)
	defer incStop()
	go incStart()

	cfg := &api.Pipeline{
		Name: "pipeline",
		Stages: []*api.Stage{
			{Name: "counter", Address: counterAddr.String()},
			{Name: "sum", Address: sumAddr.String()},
			{Name: "inc", Address: incAddr.String()},
		},
		Links: []*api.Link{
			{
				Name:        "link-counter-sum",
				SourceStage: "counter",
				TargetStage: "sum",
				TargetField: "counter",
			},
			{
				Name:             "link-inc-sum",
				SourceStage:      "inc",
				TargetStage:      "sum",
				TargetField:      "inc",
				NumEmptyMessages: 1,
			},
			{
				Name:        "link-sum-inc",
				SourceStage: "sum",
				TargetStage: "inc",
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
	incVal := int64(0)
	for i, msg := range collect {
		counterVal := int64(i + 1)
		if diff := cmp.Diff(counterVal, msg.Counter.Val); diff != "" {
			t.Fatalf("mismatch on counter value %d:\n%s", i, diff)
		}
		if diff := cmp.Diff(incVal, msg.Inc.Val); diff != "" {
			t.Fatalf("mismatch on inc value %d:\n%s", i, diff)
		}
		incVal = counterVal + incVal + 1
	}
}

type cycleCounter struct {
	integration.CycleCounterServer
	counter int64
}

func (s *cycleCounter) Generate(
	_ context.Context, _ *emptypb.Empty,
) (*integration.CycleValMessage, error) {
	delay := time.Duration(rand.Int63n(5))
	time.Sleep(delay * time.Millisecond)
	val := atomic.AddInt64(&s.counter, 1)
	return &integration.CycleValMessage{Val: val}, nil
}

type cycleSum struct {
	integration.CycleSumServer
	max     int
	collect *[]*integration.CycleSumMessage
	done    chan<- struct{}
	mu      sync.Mutex
}

func (s *cycleSum) Sum(
	_ context.Context, req *integration.CycleSumMessage,
) (*integration.CycleValMessage, error) {
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
	val := req.Inc.Val + req.Counter.Val
	return &integration.CycleValMessage{Val: val}, nil
}

type cycleInc struct{ integration.CycleIncServer }

func (t *cycleInc) Inc(
	_ context.Context, req *integration.CycleValMessage,
) (*integration.CycleValMessage, error) {
	delay := time.Duration(rand.Int63n(5))
	time.Sleep(delay * time.Millisecond)
	return &integration.CycleValMessage{Val: req.Val + 1}, nil
}
