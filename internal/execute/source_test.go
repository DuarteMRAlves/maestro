package execute

import (
	"context"
	"testing"

	"github.com/DuarteMRAlves/maestro/internal/message"
	"github.com/google/go-cmp/cmp"
)

func TestOnlineSourceStage_Run(t *testing.T) {
	start := int32(1)
	numRequest := 10

	output := make(chan onlineState)
	builder := message.BuildFunc(func() message.Instance {
		return testSourceMessage{}
	})
	s := newOnlineSource(start, builder, output)

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
		exp := onlineState{id: id(i + 1), msg: testSourceMessage{}}

		cmpOpts := cmp.AllowUnexported(onlineState{}, testSourceMessage{})
		if diff := cmp.Diff(exp, g, cmpOpts); diff != "" {
			t.Fatalf("mismatch on message %d:\n%s", i, diff)
		}
	}
}

type testSourceMessage struct{}

func (m testSourceMessage) Set(_ message.Field, _ message.Instance) error {
	panic("Should not set field in source test")
}

func (m testSourceMessage) Get(_ message.Field) (message.Instance, error) {
	panic("Should not get field in source test")
}
