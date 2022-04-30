package execute

import (
	"context"
	"testing"

	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/internal/mock"
	"github.com/google/go-cmp/cmp"
)

func TestOfflineSplitStage_Run(t *testing.T) {
	inner1 := internal.NewMessageField("inner1")
	inner3 := internal.NewMessageField("inner3")
	inner := []internal.MessageField{inner1, inner3}

	f1 := internal.NewMessageField("f1")
	f2 := internal.MessageField{}
	f3 := internal.NewMessageField("f3")

	fields := []internal.MessageField{f1, f2, f3}

	input := make(chan offlineState)

	output1 := make(chan offlineState)
	output2 := make(chan offlineState)
	output3 := make(chan offlineState)

	outputs := []chan<- offlineState{output1, output2, output3}

	s := newOfflineSplitStage(fields, input, outputs)

	expected1 := []offlineState{
		newOfflineState(testSplitInnerMessage(inner[0], 1)),
		newOfflineState(testSplitInnerMessage(inner[0], 3)),
		newOfflineState(testSplitInnerMessage(inner[0], 5)),
	}
	expected2 := []offlineState{
		newOfflineState(testSplitOuterMessage(inner, fields, 1)),
		newOfflineState(testSplitOuterMessage(inner, fields, 3)),
		newOfflineState(testSplitOuterMessage(inner, fields, 5)),
	}
	expected3 := []offlineState{
		newOfflineState(testSplitInnerMessage(inner[1], 1)),
		newOfflineState(testSplitInnerMessage(inner[1], 3)),
		newOfflineState(testSplitInnerMessage(inner[1], 5)),
	}

	go func() {
		input <- newOfflineState(testSplitOuterMessage(inner, fields, 1))
		input <- newOfflineState(testSplitOuterMessage(inner, fields, 3))
		input <- newOfflineState(testSplitOuterMessage(inner, fields, 5))
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

	cmpOpts := cmp.AllowUnexported(offlineState{})
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

func TestOnlineSplitStage_Run(t *testing.T) {
	inner1 := internal.NewMessageField("inner1")
	inner3 := internal.NewMessageField("inner3")
	inner := []internal.MessageField{inner1, inner3}

	f1 := internal.NewMessageField("f1")
	f2 := internal.MessageField{}
	f3 := internal.NewMessageField("f3")

	fields := []internal.MessageField{f1, f2, f3}

	input := make(chan onlineState)

	output1 := make(chan onlineState)
	output2 := make(chan onlineState)
	output3 := make(chan onlineState)

	outputs := []chan<- onlineState{output1, output2, output3}

	s := newOnlineSplitStage(fields, input, outputs)

	expected1 := []onlineState{
		newOnlineState(id(1), testSplitInnerMessage(inner[0], 1)),
		newOnlineState(id(3), testSplitInnerMessage(inner[0], 3)),
		newOnlineState(id(5), testSplitInnerMessage(inner[0], 5)),
	}
	expected2 := []onlineState{
		newOnlineState(id(1), testSplitOuterMessage(inner, fields, 1)),
		newOnlineState(id(3), testSplitOuterMessage(inner, fields, 3)),
		newOnlineState(id(5), testSplitOuterMessage(inner, fields, 5)),
	}
	expected3 := []onlineState{
		newOnlineState(id(1), testSplitInnerMessage(inner[1], 1)),
		newOnlineState(id(3), testSplitInnerMessage(inner[1], 3)),
		newOnlineState(id(5), testSplitInnerMessage(inner[1], 5)),
	}

	go func() {
		input <- newOnlineState(id(1), testSplitOuterMessage(inner, fields, 1))
		input <- newOnlineState(id(3), testSplitOuterMessage(inner, fields, 3))
		input <- newOnlineState(id(5), testSplitOuterMessage(inner, fields, 5))
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

	cmpOpts := cmp.AllowUnexported(onlineState{})
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
			innerMsg := testSplitInnerMessage(inner[innerIdx], val)
			msgFields[f] = innerMsg
			innerIdx++
		}
	}
	msg := &mock.Message{Fields: msgFields}
	return msg
}

func testSplitInnerMessage(field internal.MessageField, val int32) *mock.Message {
	fields := map[internal.MessageField]interface{}{field: val}
	return &mock.Message{Fields: fields}
}
