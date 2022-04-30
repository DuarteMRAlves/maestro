package execute

import (
	"context"
	"testing"

	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/internal/mock"
	"github.com/google/go-cmp/cmp"
)

func TestOfflineMergeStage_Run(t *testing.T) {
	inner1 := internal.NewMessageField("inner1")
	inner2 := internal.NewMessageField("inner2")
	inner3 := internal.NewMessageField("inner3")
	inner := []internal.MessageField{inner1, inner2, inner3}

	f1 := internal.NewMessageField("f1")
	f2 := internal.NewMessageField("f2")
	f3 := internal.NewMessageField("f3")
	fields := []internal.MessageField{f1, f2, f3}

	input1 := make(chan offlineState)
	defer close(input1)
	input2 := make(chan offlineState)
	defer close(input2)
	input3 := make(chan offlineState)
	defer close(input3)
	inputs := []<-chan offlineState{input1, input2, input3}

	output := make(chan offlineState)

	s := newOfflineMergeStage(fields, inputs, output, mock.NewGen())

	expected := []offlineState{
		newOfflineState(testMergeOuterMessage(inner, fields, 1)),
		newOfflineState(testMergeOuterMessage(inner, fields, 2)),
		newOfflineState(testMergeOuterMessage(inner, fields, 3)),
		newOfflineState(testMergeOuterMessage(inner, fields, 6)),
	}

	go func() {
		input1 <- newOfflineState(testMergeInnerMessage(inner[0], 1))
		input1 <- newOfflineState(testMergeInnerMessage(inner[0], 2))
		input1 <- newOfflineState(testMergeInnerMessage(inner[0], 3))
		input1 <- newOfflineState(testMergeInnerMessage(inner[0], 6))
	}()

	go func() {
		input2 <- newOfflineState(testMergeInnerMessage(inner[1], 1))
		input2 <- newOfflineState(testMergeInnerMessage(inner[1], 2))
		input2 <- newOfflineState(testMergeInnerMessage(inner[1], 3))
		input2 <- newOfflineState(testMergeInnerMessage(inner[1], 6))
	}()

	go func() {
		input3 <- newOfflineState(testMergeInnerMessage(inner[2], 1))
		input3 <- newOfflineState(testMergeInnerMessage(inner[2], 2))
		input3 <- newOfflineState(testMergeInnerMessage(inner[2], 3))
		input3 <- newOfflineState(testMergeInnerMessage(inner[2], 6))
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
		cmpOpts := cmp.AllowUnexported(offlineState{})
		if diff := cmp.Diff(exp, out, cmpOpts); diff != "" {
			t.Fatalf("mismatch on message %d:\n%s", i, diff)
		}
	}
	cancel()
	<-done
}

func TestOnlineMergeStage_Run(t *testing.T) {
	inner1 := internal.NewMessageField("inner1")
	inner2 := internal.NewMessageField("inner2")
	inner3 := internal.NewMessageField("inner3")
	inner := []internal.MessageField{inner1, inner2, inner3}

	f1 := internal.NewMessageField("f1")
	f2 := internal.NewMessageField("f2")
	f3 := internal.NewMessageField("f3")
	fields := []internal.MessageField{f1, f2, f3}

	input1 := make(chan onlineState)
	defer close(input1)
	input2 := make(chan onlineState)
	defer close(input2)
	input3 := make(chan onlineState)
	defer close(input3)
	inputs := []<-chan onlineState{input1, input2, input3}

	output := make(chan onlineState)

	s := newOnlineMergeStage(fields, inputs, output, mock.NewGen())

	expected := []onlineState{
		newOnlineState(3, testMergeOuterMessage(inner, fields, 3)),
		newOnlineState(6, testMergeOuterMessage(inner, fields, 6)),
	}

	go func() {
		input1 <- newOnlineState(1, testMergeInnerMessage(inner[0], 1))
		input1 <- newOnlineState(2, testMergeInnerMessage(inner[0], 2))
		input1 <- newOnlineState(3, testMergeInnerMessage(inner[0], 3))
		input1 <- newOnlineState(6, testMergeInnerMessage(inner[0], 6))
	}()

	go func() {
		input2 <- newOnlineState(2, testMergeInnerMessage(inner[1], 2))
		input2 <- newOnlineState(3, testMergeInnerMessage(inner[1], 3))
		input2 <- newOnlineState(5, testMergeInnerMessage(inner[1], 5))
		input2 <- newOnlineState(6, testMergeInnerMessage(inner[1], 6))
	}()

	go func() {
		input3 <- newOnlineState(1, testMergeInnerMessage(inner[2], 2))
		input3 <- newOnlineState(3, testMergeInnerMessage(inner[2], 3))
		input3 <- newOnlineState(5, testMergeInnerMessage(inner[2], 5))
		input3 <- newOnlineState(6, testMergeInnerMessage(inner[2], 6))
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
		cmpOpts := cmp.AllowUnexported(onlineState{})
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
		innerMsg := testMergeInnerMessage(inner[i], val)
		msgFields[f] = innerMsg
	}
	msg := &mock.Message{Fields: msgFields}
	return msg
}

func testMergeInnerMessage(field internal.MessageField, val int32) *mock.Message {
	fields := map[internal.MessageField]interface{}{field: val}
	return &mock.Message{Fields: fields}
}
