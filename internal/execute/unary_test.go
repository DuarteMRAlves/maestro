package execute

import (
	"context"
	"fmt"
	"testing"

	"github.com/DuarteMRAlves/maestro/internal/compiled"
	"github.com/DuarteMRAlves/maestro/internal/message"
	"github.com/DuarteMRAlves/maestro/internal/method"
	"github.com/google/go-cmp/cmp"
)

func TestUnaryStage_Run(t *testing.T) {
	var received []state

	stageDone := make(chan struct{})
	receiveDone := make(chan struct{})

	requests := []testUnaryMessage{{"val1"}, {"val2"}, {"val3"}}
	states := []state{
		newState(requests[0]),
		newState(requests[1]),
		newState(requests[2]),
	}

	input := make(chan state, len(requests))
	output := make(chan state, len(requests))

	name := createStageName(t, "test-stage")
	dialer := testDialer{}
	stage := newUnary(name, input, output, dialer, logger{debug: true})

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
		exp := state{
			msg: testUnaryMessage{
				val: fmt.Sprintf("val%dval%d", i+1, i+1),
			},
		}
		cmpOpts := cmp.AllowUnexported(state{}, testUnaryMessage{})
		if diff := cmp.Diff(exp, rcv, cmpOpts); diff != "" {
			t.Fatalf("mismatch on message %d:\n%s", i, diff)
		}
	}
}

func createStageName(t *testing.T, name string) compiled.StageName {
	stageName, err := compiled.NewStageName(name)
	if err != nil {
		t.Fatalf("create stage name %s: %s", name, err)
	}
	return stageName
}

type testUnaryMessage struct{ val string }

func (m testUnaryMessage) Set(_ message.Field, _ message.Instance) error {
	panic("Should not set field in unary test")
}

func (m testUnaryMessage) Get(_ message.Field) (message.Instance, error) {
	panic("Should not get field in unary test")
}

type testDialer struct{}

func (d testDialer) Dial() (method.Conn, error) { return testUnaryConn{}, nil }

type testUnaryConn struct{}

func (c testUnaryConn) Call(_ context.Context, req message.Instance) (
	message.Instance,
	error,
) {
	reqMsg, ok := req.(testUnaryMessage)
	if !ok {
		panic("request message is not testUnaryMessage")
	}
	replyVal := reqMsg.val + reqMsg.val
	return testUnaryMessage{replyVal}, nil
}

func (c testUnaryConn) Close() error { return nil }
