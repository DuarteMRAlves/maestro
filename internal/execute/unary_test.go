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

func TestOfflineUnaryStage_Run(t *testing.T) {
	var received []offlineState

	stageDone := make(chan struct{})
	receiveDone := make(chan struct{})

	fieldName := internal.NewMessageField("field")

	requests := testUnaryRequests(fieldName)
	states := []offlineState{
		newOfflineState(requests[0]),
		newOfflineState(requests[1]),
		newOfflineState(requests[2]),
	}

	input := make(chan offlineState, len(requests))
	output := make(chan offlineState, len(requests))

	name := createStageName(t, "test-stage")
	address := internal.NewAddress("some-address")
	clientBuilder := testUnaryClientBuilder(fieldName)
	logger := logs.New(true)
	buildFunc := offlineUnaryBuildFunc()
	stage := buildFunc(name, input, output, address, clientBuilder, logger)

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
		exp := offlineState{
			msg: &mock.Message{
				Fields: map[internal.MessageField]interface{}{
					fieldName: fmt.Sprintf("val%dval%d", i+1, i+1),
				},
			},
		}
		cmpOpts := cmp.AllowUnexported(offlineState{})
		if diff := cmp.Diff(exp, rcv, cmpOpts); diff != "" {
			t.Fatalf("mismatch on message %d:\n%s", i, diff)
		}
	}
}

func TestOnlineUnaryStage_Run(t *testing.T) {
	var received []onlineState

	stageDone := make(chan struct{})
	receiveDone := make(chan struct{})

	fieldName := internal.NewMessageField("field")

	requests := testUnaryRequests(fieldName)
	states := []onlineState{
		newOnlineState(1, requests[0]),
		newOnlineState(3, requests[1]),
		newOnlineState(5, requests[2]),
	}

	input := make(chan onlineState, len(requests))
	output := make(chan onlineState, len(requests))

	name := createStageName(t, "test-stage")
	address := internal.NewAddress("some-address")
	clientBuilder := testUnaryClientBuilder(fieldName)
	logger := logs.New(true)
	buildFunc := onlineUnaryBuildFunc()
	stage := buildFunc(name, input, output, address, clientBuilder, logger)

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

func testUnaryRequests(field internal.MessageField) []*mock.Message {
	fields1 := map[internal.MessageField]interface{}{field: "val1"}
	msg1 := &mock.Message{Fields: fields1}

	fields2 := map[internal.MessageField]interface{}{field: "val2"}
	msg2 := &mock.Message{Fields: fields2}

	fields3 := map[internal.MessageField]interface{}{field: "val3"}
	msg3 := &mock.Message{Fields: fields3}

	return []*mock.Message{msg1, msg2, msg3}
}

func testUnaryClientBuilder(field internal.MessageField) internal.UnaryClientBuilder {
	return func(_ internal.Address) (internal.UnaryClient, error) {
		return testUnaryClient{field: field}, nil
	}
}

type testUnaryClient struct {
	field internal.MessageField
}

func (c testUnaryClient) Call(_ context.Context, req internal.Message) (
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

func (c testUnaryClient) Close() error { return nil }
