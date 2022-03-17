package execute

import (
	"context"
	"errors"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/internal/mock"
	"gotest.tools/v3/assert"
	"sync"
	"sync/atomic"
	"testing"
)

func TestExecution_Linear(t *testing.T) {
	fieldName := internal.NewMessageField("field")

	max := 3
	collect := make([]*mock.Message, 0, max)
	done := make(chan struct{})

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
		MethodClientBuilder: linearSinkClientBuilder(max, &collect, done),
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
		internal.NewLinkEndpoint(sourceName, internal.NewEmptyMessageField()),
		internal.NewLinkEndpoint(
			transformName,
			internal.NewEmptyMessageField(),
		),
	)

	transformToSinkName := createLinkName(t, "link-transform-sink")
	transformToSink := internal.NewLink(
		transformToSinkName,
		internal.NewLinkEndpoint(
			transformName,
			internal.NewEmptyMessageField(),
		),
		internal.NewLinkEndpoint(sinkName, internal.NewEmptyMessageField()),
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

	executionBuilder := NewBuilder(stageLoader, linkLoader, methodLoader)

	orchestration := internal.NewOrchestration(
		createOrchName(t, "orchestration"),
		[]internal.StageName{sourceName, transformName, sinkName},
		[]internal.LinkName{sourceToTransformName, transformToSinkName},
	)

	e, err := executionBuilder(orchestration)
	assert.NilError(t, err, "build error")

	e.Start()
	<-done
	assert.NilError(t, e.Stop(), "stop error")
	assert.Equal(t, 3, len(collect), "invalid length")

	for i, msg := range collect {
		expected := int64(i)
		val, ok := msg.Fields[fieldName]
		assert.Assert(t, ok, "message does not have field %s", fieldName)
		valAsInt64, ok := val.(int64)
		assert.Assert(t, ok, "val is not an int64")
		assert.Equal(t, valAsInt64, (expected+1)*2)
	}
}

func createOrchName(t *testing.T, name string) internal.OrchestrationName {
	orchName, err := internal.NewOrchestrationName(name)
	assert.NilError(t, err, "create stage name %s", name)
	return orchName
}

func createStageName(t *testing.T, name string) internal.StageName {
	stageName, err := internal.NewStageName(name)
	assert.NilError(t, err, "create stage name %s", name)
	return stageName
}

func createLinkName(t *testing.T, name string) internal.LinkName {
	linkName, err := internal.NewLinkName(name)
	assert.NilError(t, err, "create link name %s", name)
	return linkName
}

func createMethodContext(addr internal.Address) internal.MethodContext {
	var (
		emptyService internal.OptionalService
		emptyMethod  internal.OptionalMethod
	)
	return internal.NewMethodContext(addr, emptyService, emptyMethod)
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

func (s *linearSinkClient) Call(_ context.Context, req internal.Message) (
	internal.Message,
	error,
) {
	reqMock, ok := req.(*mock.Message)
	if !ok {
		return nil, errors.New("sink request message is not *mock.Message")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	// Receive while not at full capacity
	if len(*s.collect) < s.max {
		*s.collect = append(*s.collect, reqMock)
	}
	// Notify when full. Remaining messages are discarded.
	if len(*s.collect) == s.max && s.done != nil {
		close(s.done)
		s.done = nil
	}
	return &mock.Message{}, nil
}

func (c *linearSinkClient) Close() error { return nil }
