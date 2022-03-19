package execute

import (
	"context"
	"errors"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/internal/mock"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestUnaryStage_Run(t *testing.T) {
	var received []state

	stageDone := make(chan struct{})
	receiveDone := make(chan struct{})

	fieldName := internal.NewMessageField("field")

	requests := testRequests(fieldName)
	states := []state{
		newState(1, requests[0]),
		newState(3, requests[1]),
		newState(5, requests[2]),
	}

	input := make(chan state, len(requests))
	output := make(chan state, len(requests))

	address := internal.NewAddress("some-address")
	clientBuilder := testUnaryClientBuilder(fieldName)
	stage := newUnaryStage(input, output, address, clientBuilder)

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
		exp := state{
			id: in.id,
			msg: &mock.Message{
				Fields: map[internal.MessageField]interface{}{
					fieldName: fmt.Sprintf("val%dval%d", i+1, i+1),
				},
			},
		}
		cmpOpts := cmp.AllowUnexported(state{})
		if diff := cmp.Diff(exp, rcv, cmpOpts); diff != "" {
			t.Fatalf("mismatch on message %d:\n%s", i, diff)
		}
	}
}

func testRequests(field internal.MessageField) []*mock.Message {
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
		return unaryClient{field: field}, nil
	}
}

type unaryClient struct {
	field internal.MessageField
}

func (c unaryClient) Call(_ context.Context, req internal.Message) (
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

func (c unaryClient) Close() error { return nil }

func TestSourceStage_Run(t *testing.T) {
	start := int32(1)
	numRequest := 10

	output := make(chan state)
	s := newSourceStage(start, mock.NewGen(), output)

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

	generated := make([]state, 0, numRequest)
	for i := 0; i < numRequest; i++ {
		generated = append(generated, <-output)
	}
	cancel()
	<-done

	for i, g := range generated {
		m := &mock.Message{Fields: map[internal.MessageField]interface{}{}}
		exp := state{id: id(i + 1), msg: m}

		cmpOpts := cmp.AllowUnexported(state{})
		if diff := cmp.Diff(exp, g, cmpOpts); diff != "" {
			t.Fatalf("mismatch on message %d:\n%s", i, diff)
		}
	}
}

func TestMergeStage_Run(t *testing.T) {
	inner1 := internal.NewMessageField("inner1")
	inner2 := internal.NewMessageField("inner2")
	inner3 := internal.NewMessageField("inner3")
	inner := []internal.MessageField{inner1, inner2, inner3}

	f1 := internal.NewMessageField("f1")
	f2 := internal.NewMessageField("f2")
	f3 := internal.NewMessageField("f3")
	fields := []internal.MessageField{f1, f2, f3}

	input1 := make(chan state)
	defer close(input1)
	input2 := make(chan state)
	defer close(input2)
	input3 := make(chan state)
	defer close(input3)
	inputs := []<-chan state{input1, input2, input3}

	output := make(chan state)

	s := newMergeStage(fields, inputs, output, mock.NewGen())

	expected := []state{
		newState(3, testMergeOuterMessage(inner, fields, 3)),
		newState(6, testMergeOuterMessage(inner, fields, 6)),
	}

	go func() {
		input1 <- newState(1, testInnerMessage(inner[0], 1))
		input1 <- newState(2, testInnerMessage(inner[0], 2))
		input1 <- newState(3, testInnerMessage(inner[0], 3))
		input1 <- newState(6, testInnerMessage(inner[0], 6))
	}()

	go func() {
		input2 <- newState(2, testInnerMessage(inner[1], 2))
		input2 <- newState(3, testInnerMessage(inner[1], 3))
		input2 <- newState(5, testInnerMessage(inner[1], 5))
		input2 <- newState(6, testInnerMessage(inner[1], 6))
	}()

	go func() {
		input3 <- newState(1, testInnerMessage(inner[2], 2))
		input3 <- newState(3, testInnerMessage(inner[2], 3))
		input3 <- newState(5, testInnerMessage(inner[2], 5))
		input3 <- newState(6, testInnerMessage(inner[2], 6))
	}()

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

	for i, exp := range expected {
		out := <-output
		cmpOpts := cmp.AllowUnexported(state{})
		if diff := cmp.Diff(exp, out, cmpOpts); diff != "" {
			t.Fatalf("mismatch on message %d:\n%s", i, diff)
		}
	}
	cancel()
	<-done
}

func testMergeOuterMessage(
	inner, fields []internal.MessageField,
	val int32,
) *mock.Message {
	msgFields := map[internal.MessageField]interface{}{}
	for i, f := range fields {
		innerMsg := testInnerMessage(inner[i], val)
		msgFields[f] = innerMsg
	}
	msg := &mock.Message{Fields: msgFields}
	return msg
}

func TestSplitStage_Run(t *testing.T) {
	inner1 := internal.NewMessageField("inner1")
	inner3 := internal.NewMessageField("inner3")
	inner := []internal.MessageField{inner1, inner3}

	f1 := internal.NewMessageField("f1")
	f2 := internal.MessageField{}
	f3 := internal.NewMessageField("f3")

	fields := []internal.MessageField{f1, f2, f3}

	input := make(chan state)

	output1 := make(chan state)
	output2 := make(chan state)
	output3 := make(chan state)

	outputs := []chan<- state{output1, output2, output3}

	s := newSplitStage(fields, input, outputs)

	expected1 := []state{
		newState(id(1), testInnerMessage(inner[0], 1)),
		newState(id(3), testInnerMessage(inner[0], 3)),
		newState(id(5), testInnerMessage(inner[0], 5)),
	}
	expected2 := []state{
		newState(id(1), testSplitOuterMessage(inner, fields, 1)),
		newState(id(3), testSplitOuterMessage(inner, fields, 3)),
		newState(id(5), testSplitOuterMessage(inner, fields, 5)),
	}
	expected3 := []state{
		newState(id(1), testInnerMessage(inner[1], 1)),
		newState(id(3), testInnerMessage(inner[1], 3)),
		newState(id(5), testInnerMessage(inner[1], 5)),
	}

	go func() {
		input <- newState(id(1), testSplitOuterMessage(inner, fields, 1))
		input <- newState(id(3), testSplitOuterMessage(inner, fields, 3))
		input <- newState(id(5), testSplitOuterMessage(inner, fields, 5))
	}()

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

	cmpOpts := cmp.AllowUnexported(state{})
	for i := 0; i < len(expected1); i++ {
		exp1 := expected1[i]
		out1 := <-output1
		if diff := cmp.Diff(exp1, out1, cmpOpts); diff != "" {
			t.Fatalf("mismatch on state 1 %d:\n%s", i, diff)
		}

		exp2 := expected2[i]
		out2 := <-output2
		if diff := cmp.Diff(exp2, out2, cmpOpts); diff != "" {
			t.Fatalf("mismatch on state 2 %d:\n%s", i, diff)
		}

		exp3 := expected3[i]
		out3 := <-output3
		if diff := cmp.Diff(exp3, out3, cmpOpts); diff != "" {
			t.Fatalf("mismatch on state 3 %d:\n%s", i, diff)
		}
	}
	cancel()
	<-done
}

func testSplitOuterMessage(
	inner []internal.MessageField, fields []internal.MessageField,
	val int32,
) *mock.Message {
	msgFields := map[internal.MessageField]interface{}{}
	innerIdx := 0
	for _, f := range fields {
		if !f.IsEmpty() {
			innerMsg := testInnerMessage(inner[innerIdx], val)
			msgFields[f] = innerMsg
			innerIdx++
		}
	}
	msg := &mock.Message{Fields: msgFields}
	return msg
}

func testInnerMessage(field internal.MessageField, val int32) *mock.Message {
	fields := map[internal.MessageField]interface{}{field: val}
	return &mock.Message{Fields: fields}
}
