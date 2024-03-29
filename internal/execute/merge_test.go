package execute

import (
	"context"
	"fmt"
	"testing"

	"github.com/DuarteMRAlves/maestro/internal/message"
	"github.com/google/go-cmp/cmp"
)

func TestMergeStage_Run(t *testing.T) {
	fields := []message.Field{"inner1", "inner2", "inner3"}

	input1 := make(chan state)
	defer close(input1)
	input2 := make(chan state)
	defer close(input2)
	input3 := make(chan state)
	defer close(input3)
	inputs := []<-chan state{input1, input2, input3}

	output := make(chan state)

	builder := message.BuildFunc(func() message.Instance { return &testMergeOuterMessage{} })

	s := newMerge(fields, inputs, output, builder)

	inputs1 := []*testMergeInnerMessage{{1}, {4}, {7}, {10}}
	inputs2 := []*testMergeInnerMessage{{2}, {5}, {8}, {11}}
	inputs3 := []*testMergeInnerMessage{{3}, {6}, {9}, {12}}

	expected := []state{
		newState(&testMergeOuterMessage{inputs1[0], inputs2[0], inputs3[0]}),
		newState(&testMergeOuterMessage{inputs1[1], inputs2[1], inputs3[1]}),
		newState(&testMergeOuterMessage{inputs1[2], inputs2[2], inputs3[2]}),
		newState(&testMergeOuterMessage{inputs1[3], inputs2[3], inputs3[3]}),
	}

	go func() {
		input1 <- newState(inputs1[0])
		input1 <- newState(inputs1[1])
		input1 <- newState(inputs1[2])
		input1 <- newState(inputs1[3])
	}()

	go func() {
		input2 <- newState(inputs2[0])
		input2 <- newState(inputs2[1])
		input2 <- newState(inputs2[2])
		input2 <- newState(inputs2[3])
	}()

	go func() {
		input3 <- newState(inputs3[0])
		input3 <- newState(inputs3[1])
		input3 <- newState(inputs3[2])
		input3 <- newState(inputs3[3])
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
			state{}, testMergeInnerMessage{}, testMergeOuterMessage{},
		)
		if diff := cmp.Diff(exp, out, cmpOpts); diff != "" {
			t.Fatalf("mismatch on message %d:\n%s", i, diff)
		}
	}
	cancel()
	<-done
}

type testMergeInnerMessage struct{ val int32 }

func (m *testMergeInnerMessage) Set(_ message.Field, _ message.Instance) error {
	panic("Should not set field for inner message in merge test")
}

func (m *testMergeInnerMessage) Get(_ message.Field) (message.Instance, error) {
	panic("Should not get field for inner message in merge test")
}

type testMergeOuterMessage struct {
	inner1 *testMergeInnerMessage
	inner2 *testMergeInnerMessage
	inner3 *testMergeInnerMessage
}

func (m *testMergeOuterMessage) Set(f message.Field, v message.Instance) error {
	inner, ok := v.(*testMergeInnerMessage)
	if !ok {
		panic("Set field for merge outer message did not receive inner message")
	}
	switch f {
	case "inner1":
		m.inner1 = inner
	case "inner2":
		m.inner2 = inner
	case "inner3":
		m.inner3 = inner
	default:
		msg := fmt.Sprintf("Set field for merge outer message received unknown field: %s", f)
		panic(msg)
	}
	return nil
}

func (m *testMergeOuterMessage) Get(_ message.Field) (message.Instance, error) {
	panic("Should not get field for outer message in merge test")
}
