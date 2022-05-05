package execute

import (
	"context"
	"testing"

	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/internal/mock"
	"github.com/google/go-cmp/cmp"
)

func TestOnlineSourceStage_Run(t *testing.T) {
	start := int32(1)
	numRequest := 10

	output := make(chan onlineState)
	buildFunc := onlineSourceBuildFunc(start)
	s := buildFunc(mock.NewGen(), output)

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
		m := &mock.Message{Fields: map[internal.MessageField]interface{}{}}
		exp := onlineState{id: id(i + 1), msg: m}

		cmpOpts := cmp.AllowUnexported(onlineState{})
		if diff := cmp.Diff(exp, g, cmpOpts); diff != "" {
			t.Fatalf("mismatch on message %d:\n%s", i, diff)
		}
	}
}
