package execute

import (
	"context"
	"fmt"
	"testing"

	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/google/go-cmp/cmp"
)

func TestOfflineUnaryStage_Run(t *testing.T) {
	var received []offlineState

	stageDone := make(chan struct{})
	receiveDone := make(chan struct{})

	requests := []testUnaryMessage{{"val1"}, {"val2"}, {"val3"}}
	states := []offlineState{
		newOfflineState(requests[0]),
		newOfflineState(requests[1]),
		newOfflineState(requests[2]),
	}

	input := make(chan offlineState, len(requests))
	output := make(chan offlineState, len(requests))

	name := createStageName(t, "test-stage")
	address := internal.NewAddress("some-address")
	clientBuilder := testUnaryClientBuilder()
	stage := newOfflineUnary(name, input, output, address, clientBuilder, logger{debug: true})

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
			msg: testUnaryMessage{
				val: fmt.Sprintf("val%dval%d", i+1, i+1),
			},
		}
		cmpOpts := cmp.AllowUnexported(offlineState{}, testUnaryMessage{})
		if diff := cmp.Diff(exp, rcv, cmpOpts); diff != "" {
			t.Fatalf("mismatch on message %d:\n%s", i, diff)
		}
	}
}

func TestOnlineUnaryStage_Run(t *testing.T) {
	var received []onlineState

	stageDone := make(chan struct{})
	receiveDone := make(chan struct{})

	requests := []testUnaryMessage{{"val1"}, {"val2"}, {"val3"}}
	states := []onlineState{
		newOnlineState(1, requests[0]),
		newOnlineState(3, requests[1]),
		newOnlineState(5, requests[2]),
	}

	input := make(chan onlineState, len(requests))
	output := make(chan onlineState, len(requests))

	name := createStageName(t, "test-stage")
	address := internal.NewAddress("some-address")
	clientBuilder := testUnaryClientBuilder()
	stage := newOnlineUnary(name, input, output, address, clientBuilder, logger{debug: true})

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
			msg: testUnaryMessage{
				val: fmt.Sprintf("val%dval%d", i+1, i+1),
			},
		}
		cmpOpts := cmp.AllowUnexported(onlineState{}, testUnaryMessage{})
		if diff := cmp.Diff(exp, rcv, cmpOpts); diff != "" {
			t.Fatalf("mismatch on message %d:\n%s", i, diff)
		}
	}
}

type testUnaryMessage struct{ val string }

func (m testUnaryMessage) SetField(_ internal.MessageField, _ internal.Message) error {
	panic("Should not set field in unary test")
}

func (m testUnaryMessage) GetField(_ internal.MessageField) (internal.Message, error) {
	panic("Should not get field in unary test")
}

func testUnaryClientBuilder() internal.UnaryClientBuilder {
	return func(_ internal.Address) (internal.UnaryClient, error) {
		return testUnaryClient{}, nil
	}
}

type testUnaryClient struct{}

func (c testUnaryClient) Call(_ context.Context, req internal.Message) (
	internal.Message,
	error,
) {
	reqMsg, ok := req.(testUnaryMessage)
	if !ok {
		panic("request message is not testUnaryMessage")
	}
	replyVal := reqMsg.val + reqMsg.val
	return testUnaryMessage{replyVal}, nil
}

func (c testUnaryClient) Close() error { return nil }
