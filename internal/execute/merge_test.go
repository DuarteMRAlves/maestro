package execute

import (
	"context"
	"fmt"
	"testing"

	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/google/go-cmp/cmp"
)

func TestOfflineMergeStage_Run(t *testing.T) {
	inner1 := internal.NewMessageField("inner1")
	inner2 := internal.NewMessageField("inner2")
	inner3 := internal.NewMessageField("inner3")

	fields := []internal.MessageField{inner1, inner2, inner3}

	input1 := make(chan offlineState)
	defer close(input1)
	input2 := make(chan offlineState)
	defer close(input2)
	input3 := make(chan offlineState)
	defer close(input3)
	inputs := []<-chan offlineState{input1, input2, input3}

	output := make(chan offlineState)

	gen := func() internal.Message { return &testMergeOuterMessage{} }

	s := newOfflineMerge(fields, inputs, output, gen)

	inputs1 := []*testMergeInnerMessage{{1}, {4}, {7}, {10}}
	inputs2 := []*testMergeInnerMessage{{2}, {5}, {8}, {11}}
	inputs3 := []*testMergeInnerMessage{{3}, {6}, {9}, {12}}

	expected := []offlineState{
		newOfflineState(&testMergeOuterMessage{inputs1[0], inputs2[0], inputs3[0]}),
		newOfflineState(&testMergeOuterMessage{inputs1[1], inputs2[1], inputs3[1]}),
		newOfflineState(&testMergeOuterMessage{inputs1[2], inputs2[2], inputs3[2]}),
		newOfflineState(&testMergeOuterMessage{inputs1[3], inputs2[3], inputs3[3]}),
	}

	go func() {
		input1 <- newOfflineState(inputs1[0])
		input1 <- newOfflineState(inputs1[1])
		input1 <- newOfflineState(inputs1[2])
		input1 <- newOfflineState(inputs1[3])
	}()

	go func() {
		input2 <- newOfflineState(inputs2[0])
		input2 <- newOfflineState(inputs2[1])
		input2 <- newOfflineState(inputs2[2])
		input2 <- newOfflineState(inputs2[3])
	}()

	go func() {
		input3 <- newOfflineState(inputs3[0])
		input3 <- newOfflineState(inputs3[1])
		input3 <- newOfflineState(inputs3[2])
		input3 <- newOfflineState(inputs3[3])
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
		cmpOpts := cmp.AllowUnexported(
			offlineState{}, testMergeInnerMessage{}, testMergeOuterMessage{},
		)
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

	fields := []internal.MessageField{inner1, inner2, inner3}

	input1 := make(chan onlineState)
	defer close(input1)
	input2 := make(chan onlineState)
	defer close(input2)
	input3 := make(chan onlineState)
	defer close(input3)
	inputs := []<-chan onlineState{input1, input2, input3}

	output := make(chan onlineState)

	gen := func() internal.Message { return &testMergeOuterMessage{} }

	s := newOnlineMerge(fields, inputs, output, gen)

	inputs1 := []*testMergeInnerMessage{{1}, {4}, {7}, {10}}
	inputs2 := []*testMergeInnerMessage{{2}, {5}, {8}, {11}}
	inputs3 := []*testMergeInnerMessage{{3}, {6}, {9}, {12}}

	expected := []onlineState{
		// The messages with id 3 are at the indexes 2, 1, 1 for the inputs arrays.
		newOnlineState(3, &testMergeOuterMessage{inputs1[2], inputs2[1], inputs3[1]}),
		newOnlineState(6, &testMergeOuterMessage{inputs1[3], inputs2[3], inputs3[3]}),
	}

	go func() {
		input1 <- newOnlineState(1, inputs1[0])
		input1 <- newOnlineState(2, inputs1[1])
		input1 <- newOnlineState(3, inputs1[2])
		input1 <- newOnlineState(6, inputs1[3])
	}()

	go func() {
		input2 <- newOnlineState(2, inputs2[0])
		input2 <- newOnlineState(3, inputs2[1])
		input2 <- newOnlineState(5, inputs2[2])
		input2 <- newOnlineState(6, inputs2[3])
	}()

	go func() {
		input3 <- newOnlineState(1, inputs3[0])
		input3 <- newOnlineState(3, inputs3[1])
		input3 <- newOnlineState(5, inputs3[2])
		input3 <- newOnlineState(6, inputs3[3])
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
		cmpOpts := cmp.AllowUnexported(
			onlineState{}, testMergeInnerMessage{}, testMergeOuterMessage{},
		)
		if diff := cmp.Diff(exp, out, cmpOpts); diff != "" {
			t.Fatalf("mismatch on message %d:\n%s", i, diff)
		}
	}
	cancel()
	<-done
}

type testMergeInnerMessage struct{ val int32 }

func (m *testMergeInnerMessage) SetField(_ internal.MessageField, _ internal.Message) error {
	panic("Should not set field for inner message in merge test")
}

func (m *testMergeInnerMessage) GetField(_ internal.MessageField) (internal.Message, error) {
	panic("Should not get field for inner message in merge test")
}

type testMergeOuterMessage struct {
	inner1 *testMergeInnerMessage
	inner2 *testMergeInnerMessage
	inner3 *testMergeInnerMessage
}

func (m *testMergeOuterMessage) SetField(f internal.MessageField, v internal.Message) error {
	inner, ok := v.(*testMergeInnerMessage)
	if !ok {
		panic("Set field for merge outer message did not receive inner message")
	}
	switch f.Unwrap() {
	case "inner1":
		m.inner1 = inner
	case "inner2":
		m.inner2 = inner
	case "inner3":
		m.inner3 = inner
	default:
		msg := fmt.Sprintf("Set field for merge outer message received unknown field: %s", f.Unwrap())
		panic(msg)
	}
	return nil
}

func (m *testMergeOuterMessage) GetField(_ internal.MessageField) (internal.Message, error) {
	panic("Should not get field for outer message in merge test")
}
