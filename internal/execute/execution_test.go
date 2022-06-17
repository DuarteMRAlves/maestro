package execute

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/DuarteMRAlves/maestro/internal/compiled"
	"github.com/google/go-cmp/cmp"
)

func TestOfflineExecution_Linear(t *testing.T) {
	max := 100
	collect := make([]*testValMsg, 0, max)
	done := make(chan struct{})

	pipelineCfg, methodLoader := setupLinear(t, max, &collect, done)
	pipelineCfg.Mode = compiled.OfflineExecution

	compilationCtx := compiled.NewContext(methodLoader)
	pipeline, err := compiled.New(compilationCtx, pipelineCfg)
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

	pipelineCfg, methodLoader := setupLinear(t, max, &collect, done)
	pipelineCfg.Mode = compiled.OnlineExecution

	compilationCtx := compiled.NewContext(methodLoader)
	pipeline, err := compiled.New(compilationCtx, pipelineCfg)
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
) (*compiled.PipelineConfig, compiled.MethodLoader) {
	cfg := &compiled.PipelineConfig{
		Name: "pipeline",
		Stages: []*compiled.StageConfig{
			{
				Name:     "source",
				MethodID: testMethodID("source"),
			},
			{
				Name:     "transform",
				MethodID: testMethodID("transform"),
			},
			{
				Name:     "sink",
				MethodID: testMethodID("sink"),
			},
		},
		Links: []*compiled.LinkConfig{
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

	sourceMethod := testMethod{
		D:   linearSourceDialFunc(),
		In:  testEmptyDesc{},
		Out: testValDesc{},
	}

	transformMethod := testMethod{
		D:   linearTransformDialFunc(),
		In:  testValDesc{},
		Out: testValDesc{},
	}
	sinkMethod := testMethod{
		D:   linearSinkDialFunc(max, collect, done),
		In:  testValDesc{},
		Out: testEmptyDesc{},
	}

	methods := map[compiled.MethodID]compiled.MethodDesc{
		testMethodID("source"):    sourceMethod,
		testMethodID("transform"): transformMethod,
		testMethodID("sink"):      sinkMethod,
	}
	methodLoader := func(mid compiled.MethodID) (compiled.MethodDesc, error) {
		m, ok := methods[mid]
		if !ok {
			panic(fmt.Sprintf("No such method: %s", mid))
		}
		return m, nil
	}

	return cfg, compiled.MethodLoaderFunc(methodLoader)
}

func linearSourceDialFunc() compiled.DialFunc {
	return func() (compiled.Conn, error) {
		return &linearSourceConn{counter: 0}, nil
	}
}

type linearSourceConn struct{ counter int64 }

func (c *linearSourceConn) Call(_ context.Context, req compiled.Message) (
	compiled.Message,
	error,
) {
	_, ok := req.(*testEmptyMsg)
	if !ok {
		panic("source request message is not testEmptyMsg")
	}
	val := atomic.AddInt64(&c.counter, 1)
	return &testValMsg{Val: val}, nil
}

func (c *linearSourceConn) Close() error { return nil }

func linearTransformDialFunc() compiled.DialFunc {
	return func() (compiled.Conn, error) {
		return &linearTransformConn{}, nil
	}
}

type linearTransformConn struct{}

func (c *linearTransformConn) Call(_ context.Context, req compiled.Message) (
	compiled.Message,
	error,
) {
	reqMsg, ok := req.(*testValMsg)
	if !ok {
		panic("transform request message is not testValMsg")
	}
	return &testValMsg{Val: 2 * reqMsg.Val}, nil
}

func (c *linearTransformConn) Close() error { return nil }

func linearSinkDialFunc(
	max int, collect *[]*testValMsg, done chan<- struct{},
) compiled.DialFunc {
	return func() (compiled.Conn, error) {
		return &linearSinkConn{
			max:     max,
			collect: collect,
			done:    done,
			mu:      sync.Mutex{},
		}, nil
	}
}

type linearSinkConn struct {
	max     int
	collect *[]*testValMsg
	done    chan<- struct{}
	mu      sync.Mutex
}

func (c *linearSinkConn) Call(_ context.Context, req compiled.Message) (
	compiled.Message,
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

func (c *linearSinkConn) Close() error { return nil }

func TestOfflineExecution_SplitAndMerge(t *testing.T) {
	max := 100
	collect := make([]*testTwoValMsg, 0, max)
	done := make(chan struct{})

	pipelineCfg, methodLoader := setupSplitAndMerge(t, max, &collect, done)
	pipelineCfg.Mode = compiled.OfflineExecution

	compilationCtx := compiled.NewContext(methodLoader)
	pipeline, err := compiled.New(compilationCtx, pipelineCfg)
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

	pipelineCfg, methodLoader := setupSplitAndMerge(t, max, &collect, done)
	pipelineCfg.Mode = compiled.OnlineExecution

	compilationCtx := compiled.NewContext(methodLoader)
	pipeline, err := compiled.New(compilationCtx, pipelineCfg)
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
) (*compiled.PipelineConfig, compiled.MethodLoader) {
	pipelineCfg := &compiled.PipelineConfig{
		Name: "pipeline",
		Stages: []*compiled.StageConfig{
			{
				Name:     "source",
				MethodID: testMethodID("source"),
			},
			{
				Name:     "transform",
				MethodID: testMethodID("transform"),
			},
			{
				Name:     "sink",
				MethodID: testMethodID("sink"),
			},
		},
		Links: []*compiled.LinkConfig{
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

	sourceMethod := testMethod{
		D:   splitAndMergeSourceDialFunc(),
		In:  testEmptyDesc{},
		Out: testValDesc{},
	}
	transformMethod := testMethod{
		D:   splitAndMergeTransformDialFunc(),
		In:  testValDesc{},
		Out: testValDesc{},
	}
	sinkMethod := testMethod{
		D:   splitAndMergeSinkDialFunc(max, collect, done),
		In:  testTwoValDesc{},
		Out: testEmptyDesc{},
	}

	methods := map[compiled.MethodID]compiled.MethodDesc{
		testMethodID("source"):    sourceMethod,
		testMethodID("transform"): transformMethod,
		testMethodID("sink"):      sinkMethod,
	}

	methodLoader := func(mid compiled.MethodID) (compiled.MethodDesc, error) {
		m, ok := methods[mid]
		if !ok {
			panic(fmt.Sprintf("No such method: %s", mid))
		}
		return m, nil
	}

	return pipelineCfg, compiled.MethodLoaderFunc(methodLoader)
}

func splitAndMergeSourceDialFunc() compiled.DialFunc {
	return func() (compiled.Conn, error) {
		return &splitAndMergeSourceConn{counter: 0}, nil
	}
}

type splitAndMergeSourceConn struct{ counter int64 }

func (c *splitAndMergeSourceConn) Call(
	_ context.Context, req compiled.Message,
) (compiled.Message, error) {
	_, ok := req.(*testEmptyMsg)
	if !ok {
		panic("source request message is not testEmptyMsg")
	}
	return &testValMsg{Val: atomic.AddInt64(&c.counter, 1)}, nil
}

func (c *splitAndMergeSourceConn) Close() error { return nil }

func splitAndMergeTransformDialFunc() compiled.DialFunc {
	return func() (compiled.Conn, error) {
		return &splitAndMergeTransformConn{}, nil
	}
}

type splitAndMergeTransformConn struct{}

func (c *splitAndMergeTransformConn) Call(
	_ context.Context, req compiled.Message,
) (compiled.Message, error) {
	reqMsg, ok := req.(*testValMsg)
	if !ok {
		panic("transform request message is not testValMsg")
	}
	return &testValMsg{Val: 3 * reqMsg.Val}, nil
}

func (c *splitAndMergeTransformConn) Close() error { return nil }

func splitAndMergeSinkDialFunc(
	max int, collect *[]*testTwoValMsg, done chan<- struct{},
) compiled.DialFunc {
	return func() (compiled.Conn, error) {
		return &splitAndMergeSinkConn{
			max:     max,
			collect: collect,
			done:    done,
			mu:      sync.Mutex{},
		}, nil
	}
}

type splitAndMergeSinkConn struct {
	max     int
	collect *[]*testTwoValMsg
	done    chan<- struct{}
	mu      sync.Mutex
}

func (c *splitAndMergeSinkConn) Call(_ context.Context, req compiled.Message) (
	compiled.Message,
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

func (c *splitAndMergeSinkConn) Close() error { return nil }

func TestOfflineExecution_Slow(t *testing.T) {
	max := 100
	collect := make([]*testValMsg, 0, max)
	done := make(chan struct{})

	pipelineCfg, methodLoader := setupSlow(t, max, &collect, done)
	pipelineCfg.Mode = compiled.OfflineExecution

	compilationCtx := compiled.NewContext(methodLoader)
	pipeline, err := compiled.New(compilationCtx, pipelineCfg)
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

	pipelineCfg, methodLoader := setupSlow(t, max, &collect, done)
	pipelineCfg.Mode = compiled.OnlineExecution

	compilationCtx := compiled.NewContext(methodLoader)
	pipeline, err := compiled.New(compilationCtx, pipelineCfg)
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
) (*compiled.PipelineConfig, compiled.MethodLoader) {
	pipelineCfg := &compiled.PipelineConfig{
		Name: "pipeline",
		Stages: []*compiled.StageConfig{
			{
				Name:     "source",
				MethodID: testMethodID("source"),
			},
			{
				Name:     "transform",
				MethodID: testMethodID("transform"),
			},
			{
				Name:     "sink",
				MethodID: testMethodID("sink"),
			},
		},
		Links: []*compiled.LinkConfig{
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

	sourceMethod := testMethod{
		D:   slowSourceDialFunc(),
		In:  testEmptyDesc{},
		Out: testValDesc{},
	}

	transformMethod := testMethod{
		D:   slowTransformDialFunc(1 * time.Millisecond),
		In:  testValDesc{},
		Out: testValDesc{},
	}
	sinkMethod := testMethod{
		D:   slowSinkDialFunc(max, collect, done),
		In:  testValDesc{},
		Out: testEmptyDesc{},
	}

	methods := map[compiled.MethodID]compiled.MethodDesc{
		testMethodID("source"):    sourceMethod,
		testMethodID("transform"): transformMethod,
		testMethodID("sink"):      sinkMethod,
	}

	methodLoader := func(mid compiled.MethodID) (compiled.MethodDesc, error) {
		m, ok := methods[mid]
		if !ok {
			panic(fmt.Sprintf("No such method: %s", mid))
		}
		return m, nil
	}

	return pipelineCfg, compiled.MethodLoaderFunc(methodLoader)
}

func slowSourceDialFunc() compiled.DialFunc {
	return func() (compiled.Conn, error) {
		return &slowSourceConn{counter: 0}, nil
	}
}

type slowSourceConn struct{ counter int64 }

func (c *slowSourceConn) Call(_ context.Context, req compiled.Message) (
	compiled.Message,
	error,
) {
	_, ok := req.(*testEmptyMsg)
	if !ok {
		panic("source request message is not testEmptyMsg")
	}
	val := atomic.AddInt64(&c.counter, 1)
	return &testValMsg{Val: val}, nil
}

func (c *slowSourceConn) Close() error { return nil }

func slowTransformDialFunc(sleep time.Duration) compiled.DialFunc {
	return func() (compiled.Conn, error) {
		return &slowTransformConn{sleep: sleep}, nil
	}
}

type slowTransformConn struct{ sleep time.Duration }

func (c *slowTransformConn) Call(_ context.Context, req compiled.Message) (
	compiled.Message,
	error,
) {
	time.Sleep(c.sleep)
	reqMsg, ok := req.(*testValMsg)
	if !ok {
		panic("transform request message is not testValMsg")
	}
	return &testValMsg{Val: 2 * reqMsg.Val}, nil
}

func (c *slowTransformConn) Close() error { return nil }

func slowSinkDialFunc(
	max int,
	collect *[]*testValMsg,
	done chan<- struct{},
) compiled.DialFunc {
	return func() (compiled.Conn, error) {
		return &slowSinkConn{
			max:     max,
			collect: collect,
			done:    done,
			mu:      sync.Mutex{},
		}, nil
	}
}

type slowSinkConn struct {
	max     int
	collect *[]*testValMsg
	done    chan<- struct{}
	mu      sync.Mutex
}

func (c *slowSinkConn) Call(_ context.Context, req compiled.Message) (
	compiled.Message,
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

func (c *slowSinkConn) Close() error { return nil }

type testMethodID string

func (mid testMethodID) String() string {
	return string(mid)
}

type testMethod struct {
	D   compiled.Dialer
	In  compiled.MessageDesc
	Out compiled.MessageDesc
}

func (m testMethod) Dial() (compiled.Conn, error) {
	return m.D.Dial()
}

func (m testMethod) Input() compiled.MessageDesc {
	return m.In
}

func (m testMethod) Output() compiled.MessageDesc {
	return m.Out
}

type testEmptyMsg struct{}

func (m *testEmptyMsg) SetField(_ compiled.MessageField, _ compiled.Message) error {
	panic("Should not set field in empty message")
}

func (m *testEmptyMsg) GetField(_ compiled.MessageField) (compiled.Message, error) {
	panic("Should not get field in empty message")
}

type testEmptyDesc struct{}

func (d testEmptyDesc) Compatible(other compiled.MessageDesc) bool {
	_, ok := other.(testEmptyDesc)
	return ok
}

func (d testEmptyDesc) EmptyGen() compiled.EmptyMessageGen {
	return func() compiled.Message { return &testEmptyMsg{} }
}

func (d testEmptyDesc) GetField(f compiled.MessageField) (compiled.MessageDesc, error) {
	panic("method get field should not be called for testEmptyDesc")
}

type testValMsg struct{ Val int64 }

func (m *testValMsg) SetField(_ compiled.MessageField, _ compiled.Message) error {
	panic("Should not set field in val message")
}

func (m *testValMsg) GetField(_ compiled.MessageField) (compiled.Message, error) {
	panic("Should not get field in val message")
}

type testValDesc struct{}

func (d testValDesc) Compatible(other compiled.MessageDesc) bool {
	_, ok := other.(testValDesc)
	return ok
}

func (d testValDesc) EmptyGen() compiled.EmptyMessageGen {
	return func() compiled.Message { return &testValMsg{} }
}

func (d testValDesc) GetField(f compiled.MessageField) (compiled.MessageDesc, error) {
	panic("method get field should not be called for testValDesc")
}

type testTwoValMsg struct {
	Orig   *testValMsg
	Transf *testValMsg
}

func (m *testTwoValMsg) SetField(f compiled.MessageField, v compiled.Message) error {
	inner, ok := v.(*testValMsg)
	if !ok {
		panic("v is not *testValMsg")
	}
	switch f {
	case "Orig":
		m.Orig = inner
	case "Transf":
		m.Transf = inner
	default:
		panic(fmt.Sprintf("Unknown field for testTwoValMsg: %s", f))
	}
	return nil
}

func (m *testTwoValMsg) GetField(_ compiled.MessageField) (compiled.Message, error) {
	panic("Should not get field in two val message")
}

type testTwoValDesc struct{}

func (d testTwoValDesc) Compatible(other compiled.MessageDesc) bool {
	_, ok := other.(testTwoValDesc)
	return ok
}

func (d testTwoValDesc) EmptyGen() compiled.EmptyMessageGen {
	return func() compiled.Message { return &testTwoValMsg{} }
}

func (d testTwoValDesc) GetField(f compiled.MessageField) (compiled.MessageDesc, error) {
	switch f {
	case "Orig", "Transf":
		return testValDesc{}, nil
	default:
		panic(fmt.Sprintf("Unknown field for testTwoValDesc: %s", f))
	}
}
