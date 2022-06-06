package execute

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/internal/compiled"
	"github.com/DuarteMRAlves/maestro/internal/mock"
	"github.com/DuarteMRAlves/maestro/internal/spec"
	"github.com/google/go-cmp/cmp"
)

func TestOfflineExecution_Linear(t *testing.T) {
	max := 100
	collect := make([]*testValMsg, 0, max)
	done := make(chan struct{})

	pipelineSpec, methodLoader := setupLinear(t, max, &collect, done)
	pipelineSpec.Mode = spec.OfflineExecution

	compilationCtx := compiled.NewContext(methodLoader)
	pipeline, err := compiled.New(compilationCtx, pipelineSpec)
	if err != nil {
		t.Fatalf("compile error: %s", err)
	}

	executionBuilder := NewBuilder(logger{debug: true})
	e, err := executionBuilder(pipeline)
	if err != nil {
		t.Fatalf("build error: %s", err)
	}

	e.Start()
	<-done
	if err := e.Stop(); err != nil {
		t.Fatalf("stop error: %s", err)
	}
	if diff := cmp.Diff(max, len(collect)); diff != "" {
		t.Fatalf("mismatch on number of collected messages:\n%s", diff)
	}

	for i, msg := range collect {
		if diff := cmp.Diff(int64((i+1)*2), msg.Val); diff != "" {
			t.Fatalf("mismatch on value %d:\n%s", i, diff)
		}
	}
}

func TestOnlineExecution_Linear(t *testing.T) {
	max := 100
	collect := make([]*testValMsg, 0, max)
	done := make(chan struct{})

	pipelineSpec, methodLoader := setupLinear(t, max, &collect, done)
	pipelineSpec.Mode = spec.OnlineExecution

	compilationCtx := compiled.NewContext(methodLoader)
	pipeline, err := compiled.New(compilationCtx, pipelineSpec)
	if err != nil {
		t.Fatalf("compile error: %s", err)
	}

	executionBuilder := NewBuilder(logger{debug: true})

	e, err := executionBuilder(pipeline)
	if err != nil {
		t.Fatalf("build error: %s", err)
	}

	e.Start()
	<-done
	if err := e.Stop(); err != nil {
		t.Fatalf("stop error: %s", err)
	}
	if diff := cmp.Diff(max, len(collect)); diff != "" {
		t.Fatalf("mismatch on number of collected messages:\n%s", diff)
	}

	prev := int64(0)
	for i, msg := range collect {
		if prev >= msg.Val {
			t.Fatalf("wrong value order at %d, %d: values are %d, %d", i-1, i, prev, msg.Val)
		}
		if msg.Val%2 != 0 {
			t.Fatalf("value %d is not pair: %d", i, msg.Val)
		}
		prev = msg.Val
	}
}

func setupLinear(
	t *testing.T, max int, collect *[]*testValMsg, done chan struct{},
) (*spec.Pipeline, compiled.MethodLoader) {
	pipelineSpec := &spec.Pipeline{
		Name: "pipeline",
		Stages: []*spec.Stage{
			{
				Name:          "source",
				MethodContext: spec.MethodContext{Address: "source"},
			},
			{
				Name:          "transform",
				MethodContext: spec.MethodContext{Address: "transform"},
			},
			{
				Name:          "sink",
				MethodContext: spec.MethodContext{Address: "sink"},
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

	sourceContext := createMethodContext(internal.NewAddress("source"))
	transformContext := createMethodContext(internal.NewAddress("transform"))
	sinkContext := createMethodContext(internal.NewAddress("sink"))

	sourceMethod := mock.Method{
		MethodClientBuilder: linearSourceClientBuilder(),
		In:                  testEmptyDesc{},
		Out:                 testValDesc{},
	}

	transformMethod := mock.Method{
		MethodClientBuilder: linearTransformClientBuilder(),
		In:                  testValDesc{},
		Out:                 testValDesc{},
	}
	sinkMethod := mock.Method{
		MethodClientBuilder: linearSinkClientBuilder(max, collect, done),
		In:                  testValDesc{},
		Out:                 testEmptyDesc{},
	}

	methods := map[internal.MethodContext]internal.UnaryMethod{
		sourceContext:    sourceMethod,
		transformContext: transformMethod,
		sinkContext:      sinkMethod,
	}
	methodLoader := func(methodCtx internal.MethodContext) (internal.UnaryMethod, error) {
		m, ok := methods[methodCtx]
		if !ok {
			panic(fmt.Sprintf("No such method: %s", methodCtx))
		}
		return m, nil
	}

	return pipelineSpec, compiled.MethodLoaderFunc(methodLoader)
}

func linearSourceClientBuilder() internal.UnaryClientBuilder {
	return func(_ internal.Address) (internal.UnaryClient, error) {
		return &linearSourceClient{counter: 0}, nil
	}
}

type linearSourceClient struct{ counter int64 }

func (c *linearSourceClient) Call(_ context.Context, req internal.Message) (
	internal.Message,
	error,
) {
	_, ok := req.(*testEmptyMsg)
	if !ok {
		panic("source request message is not testEmptyMsg")
	}
	val := atomic.AddInt64(&c.counter, 1)
	return &testValMsg{Val: val}, nil
}

func (c *linearSourceClient) Close() error { return nil }

func linearTransformClientBuilder() internal.UnaryClientBuilder {
	return func(_ internal.Address) (internal.UnaryClient, error) {
		c := &linearTransformClient{}
		return c, nil
	}
}

type linearTransformClient struct{}

func (c *linearTransformClient) Call(_ context.Context, req internal.Message) (
	internal.Message,
	error,
) {
	reqMsg, ok := req.(*testValMsg)
	if !ok {
		panic("transform request message is not testValMsg")
	}
	return &testValMsg{Val: 2 * reqMsg.Val}, nil
}

func (c *linearTransformClient) Close() error { return nil }

func linearSinkClientBuilder(
	max int,
	collect *[]*testValMsg,
	done chan<- struct{},
) internal.UnaryClientBuilder {
	return func(_ internal.Address) (internal.UnaryClient, error) {
		c := &linearSinkClient{
			max:     max,
			collect: collect,
			done:    done,
			mu:      sync.Mutex{},
		}
		return c, nil
	}
}

type linearSinkClient struct {
	max     int
	collect *[]*testValMsg
	done    chan<- struct{}
	mu      sync.Mutex
}

func (c *linearSinkClient) Call(_ context.Context, req internal.Message) (
	internal.Message,
	error,
) {
	reqMsg, ok := req.(*testValMsg)
	if !ok {
		panic("sink request message is not testValMsg")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	// Receive while not at full capacity
	if len(*c.collect) < c.max {
		*c.collect = append(*c.collect, reqMsg)
	}
	// Notify when full. Remaining messages are discarded.
	if len(*c.collect) == c.max && c.done != nil {
		close(c.done)
		c.done = nil
	}
	return &testEmptyMsg{}, nil
}

func (c *linearSinkClient) Close() error { return nil }

func TestOfflineExecution_SplitAndMerge(t *testing.T) {
	max := 100
	collect := make([]*testTwoValMsg, 0, max)
	done := make(chan struct{})

	pipelineSpec, methodLoader := setupSplitAndMerge(t, max, &collect, done)
	pipelineSpec.Mode = spec.OfflineExecution

	compilationCtx := compiled.NewContext(methodLoader)
	pipeline, err := compiled.New(compilationCtx, pipelineSpec)
	if err != nil {
		t.Fatalf("compile error: %s", err)
	}

	executionBuilder := NewBuilder(logger{debug: true})
	e, err := executionBuilder(pipeline)
	if err != nil {
		t.Fatalf("build error: %s", err)
	}

	e.Start()
	<-done
	if err := e.Stop(); err != nil {
		t.Fatalf("stop error: %s", err)
	}
	if diff := cmp.Diff(max, len(collect)); diff != "" {
		t.Fatalf("mismatch on number of collected messages:\n%s", diff)
	}

	for i, msg := range collect {
		if diff := cmp.Diff(int64((i + 1)), msg.Orig.Val); diff != "" {
			t.Fatalf("mismatch on orig value %d:\n%s", i, diff)
		}

		if diff := cmp.Diff(int64((i+1)*3), msg.Transf.Val); diff != "" {
			t.Fatalf("mismatch on transf value %d:\n%s", i, diff)
		}
	}
}

func TestOnlineExecution_SplitAndMerge(t *testing.T) {
	max := 100
	collect := make([]*testTwoValMsg, 0, max)
	done := make(chan struct{})

	pipelineSpec, methodLoader := setupSplitAndMerge(t, max, &collect, done)
	pipelineSpec.Mode = spec.OnlineExecution

	compilationCtx := compiled.NewContext(methodLoader)
	pipeline, err := compiled.New(compilationCtx, pipelineSpec)
	if err != nil {
		t.Fatalf("compile error: %s", err)
	}

	executionBuilder := NewBuilder(logger{debug: true})
	e, err := executionBuilder(pipeline)
	if err != nil {
		t.Fatalf("build error: %s", err)
	}

	e.Start()
	<-done
	if err := e.Stop(); err != nil {
		t.Fatalf("stop error: %s", err)
	}
	if diff := cmp.Diff(max, len(collect)); diff != "" {
		t.Fatalf("mismatch on number of collected messages:\n%s", diff)
	}

	prev := int64(0)
	for i, msg := range collect {
		origVal := msg.Orig.Val
		if prev >= origVal {
			t.Fatalf("wrong value order at %d, %d: values are %d, %d", i-1, i, prev, origVal)
		}
		transfVal := msg.Transf.Val
		if transfVal != 3*origVal {
			t.Fatalf("transf != 3 * orig at %d: orig is %d and transf is %d", i, origVal, transfVal)
		}
		prev = origVal
	}
}

func setupSplitAndMerge(
	t *testing.T, max int, collect *[]*testTwoValMsg, done chan struct{},
) (*spec.Pipeline, compiled.MethodLoader) {
	pipelineSpec := &spec.Pipeline{
		Name: "pipeline",
		Stages: []*spec.Stage{
			{
				Name:          "source",
				MethodContext: spec.MethodContext{Address: "source"},
			},
			{
				Name:          "transform",
				MethodContext: spec.MethodContext{Address: "transform"},
			},
			{
				Name:          "sink",
				MethodContext: spec.MethodContext{Address: "sink"},
			},
		},
		Links: []*spec.Link{
			{
				Name:        "link-source-transform",
				SourceStage: "source",
				TargetStage: "transform",
			},
			{
				Name:        "link-source-sink",
				SourceStage: "source",
				TargetStage: "sink",
				TargetField: "Orig",
			},
			{
				Name:        "link-transform-sink",
				SourceStage: "transform",
				TargetStage: "sink",
				TargetField: "Transf",
			},
		},
	}

	sourceContext := createMethodContext(internal.NewAddress("source"))
	transformContext := createMethodContext(internal.NewAddress("transform"))
	sinkContext := createMethodContext(internal.NewAddress("sink"))

	sourceMethod := mock.Method{
		MethodClientBuilder: splitAndMergeSourceClientBuilder(),
		In:                  testEmptyDesc{},
		Out:                 testValDesc{},
	}
	transformMethod := mock.Method{
		MethodClientBuilder: splitAndMergeTransformClientBuilder(),
		In:                  testValDesc{},
		Out:                 testValDesc{},
	}
	sinkMethod := mock.Method{
		MethodClientBuilder: splitAndMergeSinkClientBuilder(max, collect, done),
		In:                  testTwoValDesc{},
		Out:                 testEmptyDesc{},
	}

	methods := map[internal.MethodContext]internal.UnaryMethod{
		sourceContext:    sourceMethod,
		transformContext: transformMethod,
		sinkContext:      sinkMethod,
	}

	methodLoader := func(methodCtx internal.MethodContext) (internal.UnaryMethod, error) {
		m, ok := methods[methodCtx]
		if !ok {
			panic(fmt.Sprintf("No such method: %s", methodCtx))
		}
		return m, nil
	}

	return pipelineSpec, compiled.MethodLoaderFunc(methodLoader)
}

func splitAndMergeSourceClientBuilder() internal.UnaryClientBuilder {
	return func(_ internal.Address) (internal.UnaryClient, error) {
		return &splitAndMergeSourceClient{counter: 0}, nil
	}
}

type splitAndMergeSourceClient struct{ counter int64 }

func (c *splitAndMergeSourceClient) Call(
	_ context.Context, req internal.Message,
) (internal.Message, error) {
	_, ok := req.(*testEmptyMsg)
	if !ok {
		panic("source request message is not testEmptyMsg")
	}
	return &testValMsg{Val: atomic.AddInt64(&c.counter, 1)}, nil
}

func (c *splitAndMergeSourceClient) Close() error { return nil }

func splitAndMergeTransformClientBuilder() internal.UnaryClientBuilder {
	return func(_ internal.Address) (internal.UnaryClient, error) {
		return &splitAndMergeTransformClient{}, nil
	}
}

type splitAndMergeTransformClient struct{}

func (c *splitAndMergeTransformClient) Call(
	_ context.Context, req internal.Message,
) (internal.Message, error) {
	reqMsg, ok := req.(*testValMsg)
	if !ok {
		panic("transform request message is not testValMsg")
	}
	return &testValMsg{Val: 3 * reqMsg.Val}, nil
}

func (c *splitAndMergeTransformClient) Close() error { return nil }

func splitAndMergeSinkClientBuilder(
	max int, collect *[]*testTwoValMsg, done chan<- struct{},
) internal.UnaryClientBuilder {
	return func(_ internal.Address) (internal.UnaryClient, error) {
		c := &splitAndMergeSinkClient{
			max:     max,
			collect: collect,
			done:    done,
			mu:      sync.Mutex{},
		}
		return c, nil
	}
}

type splitAndMergeSinkClient struct {
	max     int
	collect *[]*testTwoValMsg
	done    chan<- struct{}
	mu      sync.Mutex
}

func (c *splitAndMergeSinkClient) Call(_ context.Context, req internal.Message) (
	internal.Message,
	error,
) {
	reqMock, ok := req.(*testTwoValMsg)
	if !ok {
		panic("sink request message is not *testTwoValMsg")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	// Receive while not at full capacity
	if len(*c.collect) < c.max {
		*c.collect = append(*c.collect, reqMock)
	}
	// Notify when full. Remaining messages are discarded.
	if len(*c.collect) == c.max && c.done != nil {
		close(c.done)
		c.done = nil
	}
	return &testEmptyMsg{}, nil
}

func (c *splitAndMergeSinkClient) Close() error { return nil }

func TestOfflineExecution_Slow(t *testing.T) {
	max := 100
	collect := make([]*testValMsg, 0, max)
	done := make(chan struct{})

	pipelineSpec, methodLoader := setupSlow(t, max, &collect, done)
	pipelineSpec.Mode = spec.OfflineExecution

	compilationCtx := compiled.NewContext(methodLoader)
	pipeline, err := compiled.New(compilationCtx, pipelineSpec)
	if err != nil {
		t.Fatalf("compile error: %s", err)
	}

	executionBuilder := NewBuilder(logger{debug: true})
	e, err := executionBuilder(pipeline)
	if err != nil {
		t.Fatalf("build error: %s", err)
	}

	e.Start()
	<-done
	if err := e.Stop(); err != nil {
		t.Fatalf("stop error: %s", err)
	}
	if diff := cmp.Diff(max, len(collect)); diff != "" {
		t.Fatalf("mismatch on number of collected messages:\n%s", diff)
	}

	for i, msg := range collect {
		if diff := cmp.Diff(int64((i+1)*2), msg.Val); diff != "" {
			t.Fatalf("mismatch on value %d:\n%s", i, diff)
		}
	}
}

func TestOnlineExecution_Slow(t *testing.T) {
	max := 100
	collect := make([]*testValMsg, 0, max)
	done := make(chan struct{})

	pipelineSpec, methodLoader := setupSlow(t, max, &collect, done)
	pipelineSpec.Mode = spec.OnlineExecution

	compilationCtx := compiled.NewContext(methodLoader)
	pipeline, err := compiled.New(compilationCtx, pipelineSpec)
	if err != nil {
		t.Fatalf("compile error: %s", err)
	}

	executionBuilder := NewBuilder(logger{debug: true})
	e, err := executionBuilder(pipeline)
	if err != nil {
		t.Fatalf("build error: %s", err)
	}

	e.Start()
	<-done
	if err := e.Stop(); err != nil {
		t.Fatalf("stop error: %s", err)
	}
	if diff := cmp.Diff(max, len(collect)); diff != "" {
		t.Fatalf("mismatch on number of collected messages:\n%s", diff)
	}

	prev := int64(0)
	for i, msg := range collect {
		if prev >= msg.Val {
			t.Fatalf("wrong value order at %d, %d: values are %d, %d", i-1, i, prev, msg.Val)
		}
		if msg.Val%2 != 0 {
			t.Fatalf("value %d is not pair: %d", i, msg.Val)
		}
		prev = msg.Val
	}
}

func setupSlow(
	t *testing.T, max int, collect *[]*testValMsg, done chan struct{},
) (*spec.Pipeline, compiled.MethodLoader) {
	pipelineSpec := &spec.Pipeline{
		Name: "pipeline",
		Stages: []*spec.Stage{
			{
				Name:          "source",
				MethodContext: spec.MethodContext{Address: "source"},
			},
			{
				Name:          "transform",
				MethodContext: spec.MethodContext{Address: "transform"},
			},
			{
				Name:          "sink",
				MethodContext: spec.MethodContext{Address: "sink"},
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

	sourceContext := createMethodContext(internal.NewAddress("source"))
	transformContext := createMethodContext(internal.NewAddress("transform"))
	sinkContext := createMethodContext(internal.NewAddress("sink"))

	sourceMethod := mock.Method{
		MethodClientBuilder: slowSourceClientBuilder(),
		In:                  testEmptyDesc{},
		Out:                 testValDesc{},
	}

	transformMethod := mock.Method{
		MethodClientBuilder: slowTransformClientBuilder(1 * time.Millisecond),
		In:                  testValDesc{},
		Out:                 testValDesc{},
	}
	sinkMethod := mock.Method{
		MethodClientBuilder: slowSinkClientBuilder(max, collect, done),
		In:                  testValDesc{},
		Out:                 testEmptyDesc{},
	}

	methods := map[internal.MethodContext]internal.UnaryMethod{
		sourceContext:    sourceMethod,
		transformContext: transformMethod,
		sinkContext:      sinkMethod,
	}
	methodLoader := func(methodCtx internal.MethodContext) (internal.UnaryMethod, error) {
		m, ok := methods[methodCtx]
		if !ok {
			panic(fmt.Sprintf("No such method: %s", methodCtx))
		}
		return m, nil
	}

	return pipelineSpec, compiled.MethodLoaderFunc(methodLoader)
}

func slowSourceClientBuilder() internal.UnaryClientBuilder {
	return func(_ internal.Address) (internal.UnaryClient, error) {
		return &slowSourceClient{counter: 0}, nil
	}
}

type slowSourceClient struct{ counter int64 }

func (c *slowSourceClient) Call(_ context.Context, req internal.Message) (
	internal.Message,
	error,
) {
	_, ok := req.(*testEmptyMsg)
	if !ok {
		panic("source request message is not testEmptyMsg")
	}
	val := atomic.AddInt64(&c.counter, 1)
	return &testValMsg{Val: val}, nil
}

func (c *slowSourceClient) Close() error { return nil }

func slowTransformClientBuilder(sleep time.Duration) internal.UnaryClientBuilder {
	return func(_ internal.Address) (internal.UnaryClient, error) {
		return &slowTransformClient{sleep: sleep}, nil
	}
}

type slowTransformClient struct{ sleep time.Duration }

func (c *slowTransformClient) Call(_ context.Context, req internal.Message) (
	internal.Message,
	error,
) {
	time.Sleep(c.sleep)
	reqMsg, ok := req.(*testValMsg)
	if !ok {
		panic("transform request message is not testValMsg")
	}
	return &testValMsg{Val: 2 * reqMsg.Val}, nil
}

func (c *slowTransformClient) Close() error { return nil }

func slowSinkClientBuilder(
	max int,
	collect *[]*testValMsg,
	done chan<- struct{},
) internal.UnaryClientBuilder {
	return func(_ internal.Address) (internal.UnaryClient, error) {
		c := &slowSinkClient{
			max:     max,
			collect: collect,
			done:    done,
			mu:      sync.Mutex{},
		}
		return c, nil
	}
}

type slowSinkClient struct {
	max     int
	collect *[]*testValMsg
	done    chan<- struct{}
	mu      sync.Mutex
}

func (c *slowSinkClient) Call(_ context.Context, req internal.Message) (
	internal.Message,
	error,
) {
	reqMsg, ok := req.(*testValMsg)
	if !ok {
		panic("sink request message is not testValMsg")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	// Receive while not at full capacity
	if len(*c.collect) < c.max {
		*c.collect = append(*c.collect, reqMsg)
	}
	// Notify when full. Remaining messages are discarded.
	if len(*c.collect) == c.max && c.done != nil {
		close(c.done)
		c.done = nil
	}
	return &testEmptyMsg{}, nil
}

func (c *slowSinkClient) Close() error { return nil }

func createMethodContext(addr internal.Address) internal.MethodContext {
	var (
		emptyService internal.Service
		emptyMethod  internal.Method
	)
	return internal.NewMethodContext(addr, emptyService, emptyMethod)
}

type testEmptyMsg struct{}

func (m *testEmptyMsg) SetField(_ internal.MessageField, _ internal.Message) error {
	panic("Should not set field in empty message")
}

func (m *testEmptyMsg) GetField(_ internal.MessageField) (internal.Message, error) {
	panic("Should not get field in empty message")
}

type testEmptyDesc struct{}

func (d testEmptyDesc) Compatible(other internal.MessageDesc) bool {
	_, ok := other.(testEmptyDesc)
	return ok
}

func (d testEmptyDesc) EmptyGen() internal.EmptyMessageGen {
	return func() internal.Message { return &testEmptyMsg{} }
}

func (d testEmptyDesc) GetField(f internal.MessageField) (internal.MessageDesc, error) {
	panic("method get field should not be called for testEmptyDesc")
}

type testValMsg struct{ Val int64 }

func (m *testValMsg) SetField(_ internal.MessageField, _ internal.Message) error {
	panic("Should not set field in val message")
}

func (m *testValMsg) GetField(_ internal.MessageField) (internal.Message, error) {
	panic("Should not get field in val message")
}

type testValDesc struct{}

func (d testValDesc) Compatible(other internal.MessageDesc) bool {
	_, ok := other.(testValDesc)
	return ok
}

func (d testValDesc) EmptyGen() internal.EmptyMessageGen {
	return func() internal.Message { return &testValMsg{} }
}

func (d testValDesc) GetField(f internal.MessageField) (internal.MessageDesc, error) {
	panic("method get field should not be called for testValDesc")
}

type testTwoValMsg struct {
	Orig   *testValMsg
	Transf *testValMsg
}

func (m *testTwoValMsg) SetField(f internal.MessageField, v internal.Message) error {
	inner, ok := v.(*testValMsg)
	if !ok {
		panic("v is not *testValMsg")
	}
	switch f.Unwrap() {
	case "Orig":
		m.Orig = inner
	case "Transf":
		m.Transf = inner
	default:
		panic(fmt.Sprintf("Unknown field for testTwoValMsg: %s", f.Unwrap()))
	}
	return nil
}

func (m *testTwoValMsg) GetField(_ internal.MessageField) (internal.Message, error) {
	panic("Should not get field in two val message")
}

type testTwoValDesc struct{}

func (d testTwoValDesc) Compatible(other internal.MessageDesc) bool {
	_, ok := other.(testTwoValDesc)
	return ok
}

func (d testTwoValDesc) EmptyGen() internal.EmptyMessageGen {
	return func() internal.Message { return &testTwoValMsg{} }
}

func (d testTwoValDesc) GetField(f internal.MessageField) (internal.MessageDesc, error) {
	switch f.Unwrap() {
	case "Orig", "Transf":
		return testValDesc{}, nil
	default:
		panic(fmt.Sprintf("Unknown field for testTwoValDesc: %s", f.Unwrap()))
	}
}
