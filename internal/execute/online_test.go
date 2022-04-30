package execute

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/internal/logs"
	"github.com/DuarteMRAlves/maestro/internal/mock"
	"github.com/google/go-cmp/cmp"
)

func TestOnlineUnaryStage_Run(t *testing.T) {
	var received []onlineState

	stageDone := make(chan struct{})
	receiveDone := make(chan struct{})

	fieldName := internal.NewMessageField("field")

	requests := testOnlineRequests(fieldName)
	states := []onlineState{
		newOnlineState(1, requests[0]),
		newOnlineState(3, requests[1]),
		newOnlineState(5, requests[2]),
	}

	input := make(chan onlineState, len(requests))
	output := make(chan onlineState, len(requests))

	name := createStageName(t, "test-stage")
	address := internal.NewAddress("some-address")
	clientBuilder := testOnlineUnaryClientBuilder(fieldName)
	logger := logs.New(true)
	stage := newOnlineUnaryStage(name, input, output, address, clientBuilder, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := stage.Run(ctx); err != nil {
			t.Errorf("run error: %s", err)
			return
		}
		close(stageDone)
	}()

	go func() {
		for i := 0; i < len(states); i++ {
			c := <-output
			received = append(received, c)
		}
		close(receiveDone)
	}()

	input <- states[0]
	input <- states[1]
	input <- states[2]
	<-receiveDone
	cancel()
	<-stageDone
	close(input)

	if diff := cmp.Diff(len(states), len(received)); diff != "" {
		t.Fatalf("mismatch on number of received states:\n%s", diff)
	}
	for i, rcv := range received {
		in := states[i]
		exp := onlineState{
			id: in.id,
			msg: &mock.Message{
				Fields: map[internal.MessageField]interface{}{
					fieldName: fmt.Sprintf("val%dval%d", i+1, i+1),
				},
			},
		}
		cmpOpts := cmp.AllowUnexported(onlineState{})
		if diff := cmp.Diff(exp, rcv, cmpOpts); diff != "" {
			t.Fatalf("mismatch on message %d:\n%s", i, diff)
		}
	}
}

func testOnlineRequests(field internal.MessageField) []*mock.Message {
	fields1 := map[internal.MessageField]interface{}{field: "val1"}
	msg1 := &mock.Message{Fields: fields1}

	fields2 := map[internal.MessageField]interface{}{field: "val2"}
	msg2 := &mock.Message{Fields: fields2}

	fields3 := map[internal.MessageField]interface{}{field: "val3"}
	msg3 := &mock.Message{Fields: fields3}

	return []*mock.Message{msg1, msg2, msg3}
}

func testOnlineUnaryClientBuilder(field internal.MessageField) internal.UnaryClientBuilder {
	return func(_ internal.Address) (internal.UnaryClient, error) {
		return testOnlineUnaryClient{field: field}, nil
	}
}

type testOnlineUnaryClient struct {
	field internal.MessageField
}

func (c testOnlineUnaryClient) Call(_ context.Context, req internal.Message) (
	internal.Message,
	error,
) {
	reqMock, ok := req.(*mock.Message)
	if !ok {
		return nil, errors.New("request message is not *mock.Message")
	}
	val1, ok := reqMock.Fields[c.field]
	if !ok {
		return nil, errors.New("request message does not have field1 field")
	}
	val1AsString, ok := val1.(string)
	if !ok {
		return nil, errors.New("request message field1 is not a string")
	}
	replyField := val1AsString + val1AsString
	repFields := map[internal.MessageField]interface{}{c.field: replyField}
	repMock := &mock.Message{Fields: repFields}
	return repMock, nil
}

func (c testOnlineUnaryClient) Close() error { return nil }

func TestOnlineSourceStage_Run(t *testing.T) {
	start := int32(1)
	numRequest := 10

	output := make(chan onlineState)
	s := newOnlineSourceStage(start, mock.NewGen(), output)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})

	go func() {
		if err := s.Run(ctx); err != nil {
			t.Errorf("run error: %s", err)
			return
		}
		close(done)
	}()

	generated := make([]onlineState, 0, numRequest)
	for i := 0; i < numRequest; i++ {
		generated = append(generated, <-output)
	}
	cancel()
	<-done

	for i, g := range generated {
		m := &mock.Message{Fields: map[internal.MessageField]interface{}{}}
		exp := onlineState{id: id(i + 1), msg: m}

		cmpOpts := cmp.AllowUnexported(onlineState{})
		if diff := cmp.Diff(exp, g, cmpOpts); diff != "" {
			t.Fatalf("mismatch on message %d:\n%s", i, diff)
		}
	}
}
