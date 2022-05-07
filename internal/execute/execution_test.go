package execute

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/internal/mock"
	"github.com/google/go-cmp/cmp"
)

func TestOfflineExecution_Linear(t *testing.T) {
	fieldName := internal.NewMessageField("field")

	max := 100
	collect := make([]*mock.Message, 0, max)
	done := make(chan struct{})

	pipeline, stageLoader, linkLoader, methodLoader := setupLinear(t, fieldName, max, &collect, done)
	pipeline = internal.FromPipeline(pipeline, internal.WithOfflineExec())

	executionBuilder := NewBuilder(stageLoader, linkLoader, methodLoader, logger{debug: true})

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
		val, ok := msg.Fields[fieldName]
		if !ok {
			t.Fatalf("field %s does not exist in msg %d", fieldName, i)
		}
		curr, ok := val.(int64)
		if !ok {
			format := "type mismatch in value %d: expected int64, got %s"
			t.Fatalf(format, i, reflect.TypeOf(val))
		}
		exp := int64((i + 1) * 2)
		if diff := cmp.Diff(exp, curr); diff != "" {
			t.Fatalf("mismatch on value %d:\n%s", i, diff)
		}
	}
}

func TestOnlineExecution_Linear(t *testing.T) {
	fieldName := internal.NewMessageField("field")

	max := 100
	collect := make([]*mock.Message, 0, max)
	done := make(chan struct{})

	pipeline, stageLoader, linkLoader, methodLoader := setupLinear(t, fieldName, max, &collect, done)
	pipeline = internal.FromPipeline(pipeline, internal.WithOnlineExec())

	executionBuilder := NewBuilder(stageLoader, linkLoader, methodLoader, logger{debug: true})

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
		val, ok := msg.Fields[fieldName]
		if !ok {
			t.Fatalf("field %s does not exist in msg %d", fieldName, i)
		}
		curr, ok := val.(int64)
		if !ok {
			format := "type mismatch in value %d: expected int64, got %s"
			t.Fatalf(format, i, reflect.TypeOf(val))
		}
		if prev >= curr {
			t.Fatalf("wrong value order at %d, %d: values are %d, %d", i-1, i, prev, curr)
		}
		if curr%2 != 0 {
			t.Fatalf("value %d is not pair: %d", i, curr)
		}
		prev = curr
	}
}

func setupLinear(
	t *testing.T, fieldName internal.MessageField, max int, collect *[]*mock.Message, done chan struct{},
) (internal.Pipeline, *mock.StageStorage, *mock.LinkStorage, *mock.MethodLoader) {
	emptyDesc := mock.MessageDescriptor{Ident: "empty"}
	linearMsgDesc := mock.MessageDescriptor{Ident: "message"}

	sourceMethod := mock.Method{
		MethodClientBuilder: linearSourceClientBuilder(fieldName),
		In:                  emptyDesc,
		Out:                 linearMsgDesc,
	}

	transformMethod := mock.Method{
		MethodClientBuilder: linearTransformClientBuilder(fieldName),
		In:                  linearMsgDesc,
		Out:                 linearMsgDesc,
	}
	sinkMethod := mock.Method{
		MethodClientBuilder: linearSinkClientBuilder(max, collect, done),
		In:                  linearMsgDesc,
		Out:                 emptyDesc,
	}

	sourceName := createStageName(t, "source")
	transformName := createStageName(t, "transform")
	sinkName := createStageName(t, "sink")

	sourceAddr := internal.NewAddress("source")
	transformAddr := internal.NewAddress("transform")
	sinkAddr := internal.NewAddress("sink")

	sourceContext := createMethodContext(sourceAddr)
	transformContext := createMethodContext(transformAddr)
	sinkContext := createMethodContext(sinkAddr)

	sourceStage := internal.NewStage(sourceName, sourceContext)
	transformStage := internal.NewStage(transformName, transformContext)
	sinkStage := internal.NewStage(sinkName, sinkContext)

	stages := map[internal.StageName]internal.Stage{
		sourceName:    sourceStage,
		transformName: transformStage,
		sinkName:      sinkStage,
	}
	stageLoader := &mock.StageStorage{Stages: stages}

	sourceToTransformName := createLinkName(t, "link-source-transform")
	sourceToTransform := internal.NewLink(
		sourceToTransformName,
		internal.NewLinkEndpoint(sourceName, internal.MessageField{}),
		internal.NewLinkEndpoint(transformName, internal.MessageField{}),
	)

	transformToSinkName := createLinkName(t, "link-transform-sink")
	transformToSink := internal.NewLink(
		transformToSinkName,
		internal.NewLinkEndpoint(transformName, internal.MessageField{}),
		internal.NewLinkEndpoint(sinkName, internal.MessageField{}),
	)

	links := map[internal.LinkName]internal.Link{
		sourceToTransformName: sourceToTransform,
		transformToSinkName:   transformToSink,
	}
	linkLoader := &mock.LinkStorage{Links: links}

	methods := map[internal.MethodContext]internal.UnaryMethod{
		sourceContext:    sourceMethod,
		transformContext: transformMethod,
		sinkContext:      sinkMethod,
	}
	methodLoader := &mock.MethodLoader{Methods: methods}

	pipeline := internal.NewPipeline(
		createPipelineName(t, "pipeline"),
		internal.WithStages(sourceName, transformName, sinkName),
		internal.WithLinks(sourceToTransformName, transformToSinkName),
	)
	return pipeline, stageLoader, linkLoader, methodLoader
}

func linearSourceClientBuilder(field internal.MessageField) internal.UnaryClientBuilder {
	return func(_ internal.Address) (internal.UnaryClient, error) {
		c := &linearSourceClient{
			field:   field,
			counter: 0,
		}
		return c, nil
	}
}

type linearSourceClient struct {
	field   internal.MessageField
	counter int64
}

func (c *linearSourceClient) Call(_ context.Context, req internal.Message) (
	internal.Message,
	error,
) {
	reqMock, ok := req.(*mock.Message)
	if !ok {
		return nil, errors.New("source request message is not *mock.Message")
	}
	if len(reqMock.Fields) != 0 {
		return nil, errors.New("source request message is not empty")
	}
	val := atomic.AddInt64(&c.counter, 1)
	repFields := map[internal.MessageField]interface{}{c.field: val}
	return &mock.Message{Fields: repFields}, nil
}

func (c *linearSourceClient) Close() error { return nil }

func linearTransformClientBuilder(field internal.MessageField) internal.UnaryClientBuilder {
	return func(_ internal.Address) (internal.UnaryClient, error) {
		c := &linearTransformClient{field: field}
		return c, nil
	}
}

type linearTransformClient struct {
	field internal.MessageField
}

func (c *linearTransformClient) Call(_ context.Context, req internal.Message) (
	internal.Message,
	error,
) {
	reqMock, ok := req.(*mock.Message)
	if !ok {
		return nil, errors.New("transform request message is not *mock.Message")
	}
	val, ok := reqMock.Fields[c.field]
	if !ok {
		return nil, fmt.Errorf(
			"transform request message does not have %s field",
			c.field,
		)
	}
	valAsInt64, ok := val.(int64)
	if !ok {
		return nil, fmt.Errorf(
			"transform request message %s is not an int64",
			c.field,
		)
	}
	replyVal := 2 * valAsInt64
	repFields := map[internal.MessageField]interface{}{c.field: replyVal}
	return &mock.Message{Fields: repFields}, nil
}

func (c *linearTransformClient) Close() error { return nil }

func linearSinkClientBuilder(
	max int,
	collect *[]*mock.Message,
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
	collect *[]*mock.Message
	done    chan<- struct{}
	mu      sync.Mutex
}

func (c *linearSinkClient) Call(_ context.Context, req internal.Message) (
	internal.Message,
	error,
) {
	reqMock, ok := req.(*mock.Message)
	if !ok {
		return nil, errors.New("sink request message is not *mock.Message")
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
	return &mock.Message{}, nil
}

func (c *linearSinkClient) Close() error { return nil }

func TestOfflineExecution_SplitAndMerge(t *testing.T) {
	fieldName := internal.NewMessageField("field")
	origFld := internal.NewMessageField("original")
	transfFld := internal.NewMessageField("transform")

	max := 100
	collect := make([]*mock.Message, 0, max)
	done := make(chan struct{})

	pipeline, stageLoader, linkLoader, methodLoader := setupSplitAndMerge(
		t, fieldName, origFld, transfFld, max, &collect, done,
	)
	pipeline = internal.FromPipeline(pipeline, internal.WithOfflineExec())

	executionBuilder := NewBuilder(stageLoader, linkLoader, methodLoader, logger{debug: true})

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
		orig, ok := msg.Fields[origFld]
		if !ok {
			t.Fatalf("field %s does not exist in msg %d", origFld, i)
		}
		origMock, ok := orig.(*mock.Message)
		if !ok {
			t.Fatalf("orig %d is not a *mock.Message", i)
		}
		origVal, ok := origMock.Fields[fieldName]
		if !ok {
			t.Fatalf("field %s does not exist in orig %d", fieldName, i)
		}
		origCurr, ok := origVal.(int64)
		if !ok {
			format := "type mismatch in value %d: expected int64, got %s"
			t.Fatalf(format, i, reflect.TypeOf(origVal))
		}
		origExp := int64((i + 1))
		if diff := cmp.Diff(origExp, origCurr); diff != "" {
			t.Fatalf("mismatch on orig value %d:\n%s", i, diff)
		}

		transf, ok := msg.Fields[transfFld]
		if !ok {
			t.Fatalf("field %s does not exist in msg %d", transfFld, i)
		}
		transfMock, ok := transf.(*mock.Message)
		if !ok {
			t.Fatalf("transf %d is not a *mock.Message", i)
		}
		transfVal, ok := transfMock.Fields[fieldName]
		if !ok {
			t.Fatalf("field %s does not exist in transf %d", fieldName, i)
		}
		transfCurr, ok := transfVal.(int64)
		if !ok {
			format := "type mismatch in value %d: expected int64, got %s"
			t.Fatalf(format, i, reflect.TypeOf(transfVal))
		}
		transfExp := int64((i + 1) * 3)
		if diff := cmp.Diff(transfExp, transfCurr); diff != "" {
			t.Fatalf("mismatch on transf value %d:\n%s", i, diff)
		}
	}
}

func TestOnlineExecution_SplitAndMerge(t *testing.T) {
	fieldName := internal.NewMessageField("field")
	origFld := internal.NewMessageField("original")
	transfFld := internal.NewMessageField("transform")

	max := 100
	collect := make([]*mock.Message, 0, max)
	done := make(chan struct{})

	pipeline, stageLoader, linkLoader, methodLoader := setupSplitAndMerge(
		t, fieldName, origFld, transfFld, max, &collect, done,
	)
	pipeline = internal.FromPipeline(pipeline, internal.WithOnlineExec())

	executionBuilder := NewBuilder(stageLoader, linkLoader, methodLoader, logger{debug: true})

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
		orig, ok := msg.Fields[origFld]
		if !ok {
			t.Fatalf("field %s does not exist in msg %d", origFld, i)
		}
		origMock, ok := orig.(*mock.Message)
		if !ok {
			t.Fatalf("orig %d is not a *mock.Message", i)
		}
		origVal, ok := origMock.Fields[fieldName]
		if !ok {
			t.Fatalf("field %s does not exist in orig %d", fieldName, i)
		}
		origCurr, ok := origVal.(int64)
		if !ok {
			format := "type mismatch in value %d: expected int64, got %s"
			t.Fatalf(format, i, reflect.TypeOf(origVal))
		}
		if prev >= origCurr {
			t.Fatalf("wrong value order at %d, %d: values are %d, %d", i-1, i, prev, origCurr)
		}

		transf, ok := msg.Fields[transfFld]
		if !ok {
			t.Fatalf("field %s does not exist in msg %d", transfFld, i)
		}
		transfMock, ok := transf.(*mock.Message)
		if !ok {
			t.Fatalf("transf %d is not a *mock.Message", i)
		}
		transfVal, ok := transfMock.Fields[fieldName]
		if !ok {
			t.Fatalf("field %s does not exist in transf %d", fieldName, i)
		}
		transfCurr, ok := transfVal.(int64)
		if !ok {
			format := "type mismatch in value %d: expected int64, got %s"
			t.Fatalf(format, i, reflect.TypeOf(transfVal))
		}
		if transfCurr != 3*origCurr {
			t.Fatalf("transf != 3 * orig at %d: orig is %d and transf is %d", i, origCurr, transfCurr)
		}
		prev = origCurr
	}
}

func setupSplitAndMerge(
	t *testing.T,
	fieldName internal.MessageField,
	originalField internal.MessageField,
	transformField internal.MessageField,
	max int,
	collect *[]*mock.Message,
	done chan struct{},
) (internal.Pipeline, *mock.StageStorage, *mock.LinkStorage, *mock.MethodLoader) {
	emptyDesc := mock.MessageDescriptor{Ident: "empty"}
	singleMsgDesc := mock.MessageDescriptor{Ident: "single"}
	mergeMsgDesc := mock.MessageDescriptor{
		Ident: "merge",
		Fields: map[internal.MessageField]internal.MessageDesc{
			originalField:  singleMsgDesc,
			transformField: singleMsgDesc,
		},
	}

	sourceMethod := mock.Method{
		MethodClientBuilder: splitAndMergeSourceClientBuilder(fieldName),
		In:                  emptyDesc,
		Out:                 singleMsgDesc,
	}
	transformMethod := mock.Method{
		MethodClientBuilder: splitAndMergeTransformClientBuilder(fieldName),
		In:                  singleMsgDesc,
		Out:                 singleMsgDesc,
	}
	sinkMethod := mock.Method{
		MethodClientBuilder: splitAndMergeSinkClientBuilder(max, collect, done),
		In:                  mergeMsgDesc,
		Out:                 emptyDesc,
	}

	sourceName := createStageName(t, "source")
	transformName := createStageName(t, "transform")
	sinkName := createStageName(t, "sink")

	sourceAddr := internal.NewAddress("source")
	transformAddr := internal.NewAddress("transform")
	sinkAddr := internal.NewAddress("sink")

	sourceContext := createMethodContext(sourceAddr)
	transformContext := createMethodContext(transformAddr)
	sinkContext := createMethodContext(sinkAddr)

	sourceStage := internal.NewStage(sourceName, sourceContext)
	transformStage := internal.NewStage(transformName, transformContext)
	sinkStage := internal.NewStage(sinkName, sinkContext)

	stages := map[internal.StageName]internal.Stage{
		sourceName:    sourceStage,
		transformName: transformStage,
		sinkName:      sinkStage,
	}
	stageLoader := &mock.StageStorage{Stages: stages}

	sourceToTransformName := createLinkName(t, "link-source-transform")
	sourceToTransform := internal.NewLink(
		sourceToTransformName,
		internal.NewLinkEndpoint(sourceName, internal.MessageField{}),
		internal.NewLinkEndpoint(transformName, internal.MessageField{}),
	)

	sourceToSinkName := createLinkName(t, "link-source-sink")
	sourceToSink := internal.NewLink(
		sourceToSinkName,
		internal.NewLinkEndpoint(sourceName, internal.MessageField{}),
		internal.NewLinkEndpoint(sinkName, originalField),
	)

	transformToSinkName := createLinkName(t, "link-transform-sink")
	transformToSink := internal.NewLink(
		transformToSinkName,
		internal.NewLinkEndpoint(transformName, internal.MessageField{}),
		internal.NewLinkEndpoint(sinkName, transformField),
	)

	links := map[internal.LinkName]internal.Link{
		sourceToTransformName: sourceToTransform,
		sourceToSinkName:      sourceToSink,
		transformToSinkName:   transformToSink,
	}
	linkLoader := &mock.LinkStorage{Links: links}

	methods := map[internal.MethodContext]internal.UnaryMethod{
		sourceContext:    sourceMethod,
		transformContext: transformMethod,
		sinkContext:      sinkMethod,
	}
	methodLoader := &mock.MethodLoader{Methods: methods}

	pipeline := internal.NewPipeline(
		createPipelineName(t, "pipeline"),
		internal.WithStages(sourceName, transformName, sinkName),
		internal.WithLinks(sourceToTransformName, sourceToSinkName, transformToSinkName),
	)
	return pipeline, stageLoader, linkLoader, methodLoader
}

func splitAndMergeSourceClientBuilder(field internal.MessageField) internal.UnaryClientBuilder {
	return func(_ internal.Address) (internal.UnaryClient, error) {
		c := &splitAndMergeSourceClient{
			field:   field,
			counter: 0,
		}
		return c, nil
	}
}

type splitAndMergeSourceClient struct {
	field   internal.MessageField
	counter int64
}

func (c *splitAndMergeSourceClient) Call(
	_ context.Context, req internal.Message,
) (internal.Message, error) {
	reqMock, ok := req.(*mock.Message)
	if !ok {
		return nil, errors.New("source request message is not *mock.Message")
	}
	if len(reqMock.Fields) != 0 {
		return nil, errors.New("source request message is not empty")
	}
	val := atomic.AddInt64(&c.counter, 1)
	repFields := map[internal.MessageField]interface{}{c.field: val}
	return &mock.Message{Fields: repFields}, nil
}

func (c *splitAndMergeSourceClient) Close() error { return nil }

func splitAndMergeTransformClientBuilder(field internal.MessageField) internal.UnaryClientBuilder {
	return func(_ internal.Address) (internal.UnaryClient, error) {
		c := &splitAndMergeTransformClient{field: field}
		return c, nil
	}
}

type splitAndMergeTransformClient struct {
	field internal.MessageField
}

func (c *splitAndMergeTransformClient) Call(
	_ context.Context, req internal.Message,
) (internal.Message, error) {
	reqMock, ok := req.(*mock.Message)
	if !ok {
		return nil, errors.New("transform request message is not *mock.Message")
	}
	val, ok := reqMock.Fields[c.field]
	if !ok {
		return nil, fmt.Errorf(
			"transform request message does not have %s field",
			c.field,
		)
	}
	valAsInt64, ok := val.(int64)
	if !ok {
		return nil, fmt.Errorf(
			"transform request message %s is not an int64",
			c.field,
		)
	}
	replyVal := 3 * valAsInt64
	repFields := map[internal.MessageField]interface{}{c.field: replyVal}
	return &mock.Message{Fields: repFields}, nil
}

func (c *splitAndMergeTransformClient) Close() error { return nil }

func splitAndMergeSinkClientBuilder(
	max int, collect *[]*mock.Message, done chan<- struct{},
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
	collect *[]*mock.Message
	done    chan<- struct{}
	mu      sync.Mutex
}

func (c *splitAndMergeSinkClient) Call(_ context.Context, req internal.Message) (
	internal.Message,
	error,
) {
	reqMock, ok := req.(*mock.Message)
	if !ok {
		return nil, errors.New("sink request message is not *mock.Message")
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
	return &mock.Message{}, nil
}

func (c *splitAndMergeSinkClient) Close() error { return nil }

func TestOfflineExecution_Slow(t *testing.T) {
	fieldName := internal.NewMessageField("field")

	max := 100
	collect := make([]*mock.Message, 0, max)
	done := make(chan struct{})

	pipeline, stageLoader, linkLoader, methodLoader := setupSlow(t, fieldName, max, &collect, done)
	pipeline = internal.FromPipeline(pipeline, internal.WithOfflineExec())
	executionBuilder := NewBuilder(stageLoader, linkLoader, methodLoader, logger{debug: true})

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
		val, ok := msg.Fields[fieldName]
		if !ok {
			t.Fatalf("field %s does not exist in msg %d", fieldName, i)
		}
		curr, ok := val.(int64)
		if !ok {
			format := "type mismatch in value %d: expected int64, got %s"
			t.Fatalf(format, i, reflect.TypeOf(val))
		}
		exp := int64((i + 1) * 2)
		if diff := cmp.Diff(exp, curr); diff != "" {
			t.Fatalf("mismatch on value %d:\n%s", i, diff)
		}
	}
}

func TestOnlineExecution_Slow(t *testing.T) {
	fieldName := internal.NewMessageField("field")

	max := 100
	collect := make([]*mock.Message, 0, max)
	done := make(chan struct{})

	pipeline, stageLoader, linkLoader, methodLoader := setupSlow(t, fieldName, max, &collect, done)
	pipeline = internal.FromPipeline(pipeline, internal.WithOnlineExec())
	executionBuilder := NewBuilder(stageLoader, linkLoader, methodLoader, logger{debug: true})

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
		val, ok := msg.Fields[fieldName]
		if !ok {
			t.Fatalf("field %s does not exist in msg %d", fieldName, i)
		}
		curr, ok := val.(int64)
		if !ok {
			format := "type mismatch in value %d: expected int64, got %s"
			t.Fatalf(format, i, reflect.TypeOf(val))
		}
		if prev >= curr {
			t.Fatalf("wrong value order at %d, %d: values are %d, %d", i-1, i, prev, curr)
		}
		if curr%2 != 0 {
			t.Fatalf("value %d is not pair: %d", i, curr)
		}
		prev = curr
	}
}

func setupSlow(
	t *testing.T, fieldName internal.MessageField, max int, collect *[]*mock.Message, done chan struct{},
) (internal.Pipeline, *mock.StageStorage, *mock.LinkStorage, *mock.MethodLoader) {
	emptyDesc := mock.MessageDescriptor{Ident: "empty"}
	linearMsgDesc := mock.MessageDescriptor{Ident: "message"}

	sourceMethod := mock.Method{
		MethodClientBuilder: slowSourceClientBuilder(fieldName),
		In:                  emptyDesc,
		Out:                 linearMsgDesc,
	}

	transformMethod := mock.Method{
		MethodClientBuilder: slowTransformClientBuilder(fieldName, 1*time.Millisecond),
		In:                  linearMsgDesc,
		Out:                 linearMsgDesc,
	}
	sinkMethod := mock.Method{
		MethodClientBuilder: slowSinkClientBuilder(max, collect, done),
		In:                  linearMsgDesc,
		Out:                 emptyDesc,
	}

	sourceName := createStageName(t, "source")
	transformName := createStageName(t, "transform")
	sinkName := createStageName(t, "sink")

	sourceAddr := internal.NewAddress("source")
	transformAddr := internal.NewAddress("transform")
	sinkAddr := internal.NewAddress("sink")

	sourceContext := createMethodContext(sourceAddr)
	transformContext := createMethodContext(transformAddr)
	sinkContext := createMethodContext(sinkAddr)

	sourceStage := internal.NewStage(sourceName, sourceContext)
	transformStage := internal.NewStage(transformName, transformContext)
	sinkStage := internal.NewStage(sinkName, sinkContext)

	stages := map[internal.StageName]internal.Stage{
		sourceName:    sourceStage,
		transformName: transformStage,
		sinkName:      sinkStage,
	}
	stageLoader := &mock.StageStorage{Stages: stages}

	sourceToTransformName := createLinkName(t, "link-source-transform")
	sourceToTransform := internal.NewLink(
		sourceToTransformName,
		internal.NewLinkEndpoint(sourceName, internal.MessageField{}),
		internal.NewLinkEndpoint(transformName, internal.MessageField{}),
	)

	transformToSinkName := createLinkName(t, "link-transform-sink")
	transformToSink := internal.NewLink(
		transformToSinkName,
		internal.NewLinkEndpoint(transformName, internal.MessageField{}),
		internal.NewLinkEndpoint(sinkName, internal.MessageField{}),
	)

	links := map[internal.LinkName]internal.Link{
		sourceToTransformName: sourceToTransform,
		transformToSinkName:   transformToSink,
	}
	linkLoader := &mock.LinkStorage{Links: links}

	methods := map[internal.MethodContext]internal.UnaryMethod{
		sourceContext:    sourceMethod,
		transformContext: transformMethod,
		sinkContext:      sinkMethod,
	}
	methodLoader := &mock.MethodLoader{Methods: methods}

	pipeline := internal.NewPipeline(
		createPipelineName(t, "pipeline"),
		internal.WithStages(sourceName, transformName, sinkName),
		internal.WithLinks(sourceToTransformName, transformToSinkName),
	)
	return pipeline, stageLoader, linkLoader, methodLoader
}

func slowSourceClientBuilder(field internal.MessageField) internal.UnaryClientBuilder {
	return func(_ internal.Address) (internal.UnaryClient, error) {
		c := &slowSourceClient{
			field:   field,
			counter: 0,
		}
		return c, nil
	}
}

type slowSourceClient struct {
	field   internal.MessageField
	counter int64
}

func (c *slowSourceClient) Call(_ context.Context, req internal.Message) (
	internal.Message,
	error,
) {
	reqMock, ok := req.(*mock.Message)
	if !ok {
		return nil, errors.New("source request message is not *mock.Message")
	}
	if len(reqMock.Fields) != 0 {
		return nil, errors.New("source request message is not empty")
	}
	val := atomic.AddInt64(&c.counter, 1)
	repFields := map[internal.MessageField]interface{}{c.field: val}
	return &mock.Message{Fields: repFields}, nil
}

func (c *slowSourceClient) Close() error { return nil }

func slowTransformClientBuilder(
	field internal.MessageField, sleep time.Duration,
) internal.UnaryClientBuilder {
	return func(_ internal.Address) (internal.UnaryClient, error) {
		c := &slowTransformClient{field: field, sleep: sleep}
		return c, nil
	}
}

type slowTransformClient struct {
	field internal.MessageField
	sleep time.Duration
}

func (c *slowTransformClient) Call(_ context.Context, req internal.Message) (
	internal.Message,
	error,
) {
	time.Sleep(c.sleep)
	reqMock, ok := req.(*mock.Message)
	if !ok {
		return nil, errors.New("transform request message is not *mock.Message")
	}
	val, ok := reqMock.Fields[c.field]
	if !ok {
		return nil, fmt.Errorf(
			"transform request message does not have %s field",
			c.field,
		)
	}
	valAsInt64, ok := val.(int64)
	if !ok {
		return nil, fmt.Errorf(
			"transform request message %s is not an int64",
			c.field,
		)
	}
	replyVal := 2 * valAsInt64
	repFields := map[internal.MessageField]interface{}{c.field: replyVal}
	return &mock.Message{Fields: repFields}, nil
}

func (c *slowTransformClient) Close() error { return nil }

func slowSinkClientBuilder(
	max int,
	collect *[]*mock.Message,
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
	collect *[]*mock.Message
	done    chan<- struct{}
	mu      sync.Mutex
}

func (c *slowSinkClient) Call(_ context.Context, req internal.Message) (
	internal.Message,
	error,
) {
	reqMock, ok := req.(*mock.Message)
	if !ok {
		return nil, errors.New("sink request message is not *mock.Message")
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
	return &mock.Message{}, nil
}

func (c *slowSinkClient) Close() error { return nil }

func createPipelineName(t *testing.T, name string) internal.PipelineName {
	pipelineName, err := internal.NewPipelineName(name)
	if err != nil {
		t.Fatalf("create pipeline name %s: %s", name, err)
	}
	return pipelineName
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
