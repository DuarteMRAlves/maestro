package execute

import (
	"context"
	"fmt"
	"testing"

	"github.com/DuarteMRAlves/maestro/internal/message"
	"github.com/google/go-cmp/cmp"
)

func TestOfflineSplitStage_Run(t *testing.T) {
	// Send full message through second output
	fields := []message.Field{"inner1", "", "inner3"}

	input := make(chan offlineState)
	defer close(input)

	output1 := make(chan offlineState)
	output2 := make(chan offlineState)
	output3 := make(chan offlineState)

	outputs := []chan<- offlineState{output1, output2, output3}

	s := newOfflineSplit(fields, input, outputs)

	inputs := []*testSplitOuterMessage{
		{&testSplitInnerMessage{1}, &testSplitInnerMessage{2}, &testSplitInnerMessage{3}},
		{&testSplitInnerMessage{4}, &testSplitInnerMessage{5}, &testSplitInnerMessage{6}},
		{&testSplitInnerMessage{7}, &testSplitInnerMessage{8}, &testSplitInnerMessage{9}},
	}

	expected1 := []offlineState{
		newOfflineState(inputs[0].inner1),
		newOfflineState(inputs[1].inner1),
		newOfflineState(inputs[2].inner1),
	}
	expected2 := []offlineState{
		newOfflineState(inputs[0]),
		newOfflineState(inputs[1]),
		newOfflineState(inputs[2]),
	}
	expected3 := []offlineState{
		newOfflineState(inputs[0].inner3),
		newOfflineState(inputs[1].inner3),
		newOfflineState(inputs[2].inner3),
	}

	go func() {
		input <- newOfflineState(inputs[0])
		input <- newOfflineState(inputs[1])
		input <- newOfflineState(inputs[2])
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

	cmpOpts := cmp.AllowUnexported(
		offlineState{}, testSplitInnerMessage{}, testSplitOuterMessage{},
	)
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
	// Send full message through second output
	fields := []message.Field{"inner1", "", "inner3"}

	input := make(chan onlineState)
	defer close(input)

	output1 := make(chan onlineState)
	output2 := make(chan onlineState)
	output3 := make(chan onlineState)

	outputs := []chan<- onlineState{output1, output2, output3}

	s := newOnlineSplit(fields, input, outputs)

	inputs := []*testSplitOuterMessage{
		{&testSplitInnerMessage{1}, &testSplitInnerMessage{2}, &testSplitInnerMessage{3}},
		{&testSplitInnerMessage{4}, &testSplitInnerMessage{5}, &testSplitInnerMessage{6}},
		{&testSplitInnerMessage{7}, &testSplitInnerMessage{8}, &testSplitInnerMessage{9}},
	}

	expected1 := []onlineState{
		newOnlineState(1, inputs[0].inner1),
		newOnlineState(3, inputs[1].inner1),
		newOnlineState(5, inputs[2].inner1),
	}
	expected2 := []onlineState{
		newOnlineState(1, inputs[0]),
		newOnlineState(3, inputs[1]),
		newOnlineState(5, inputs[2]),
	}
	expected3 := []onlineState{
		newOnlineState(1, inputs[0].inner3),
		newOnlineState(3, inputs[1].inner3),
		newOnlineState(5, inputs[2].inner3),
	}

	go func() {
		input <- newOnlineState(1, inputs[0])
		input <- newOnlineState(3, inputs[1])
		input <- newOnlineState(5, inputs[2])
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

	cmpOpts := cmp.AllowUnexported(
		onlineState{}, testSplitInnerMessage{}, testSplitOuterMessage{},
	)
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

type testSplitInnerMessage struct{ val int32 }

func (m *testSplitInnerMessage) Set(_ message.Field, _ message.Instance) error {
	panic("Should not set field for inner message in split test")
}

func (m *testSplitInnerMessage) Get(_ message.Field) (message.Instance, error) {
	panic("Should not get field for inner message in split test")
}

type testSplitOuterMessage struct {
	inner1 *testSplitInnerMessage
	inner2 *testSplitInnerMessage
	inner3 *testSplitInnerMessage
}

func (m *testSplitOuterMessage) Set(f message.Field, v message.Instance) error {
	panic("Should not set field for outer message in split test")
}

func (m *testSplitOuterMessage) Get(f message.Field) (message.Instance, error) {
	switch f {
	case "inner1":
		return m.inner1, nil
	case "inner2":
		return m.inner2, nil
	case "inner3":
		return m.inner3, nil
	default:
		msg := fmt.Sprintf("Set field for merge outer message received unknown field: %s", f)
		panic(msg)
	}
}
