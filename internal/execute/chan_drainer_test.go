package execute

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"
)

func TestChanDrainer(t *testing.T) {
	nChans := 10
	capacity := 10
	chans := make([]chan state, 0, nChans)
	defer func() {
		for _, ch := range chans {
			close(ch)
		}
	}()
	for i := 0; i < nChans; i++ {
		ch := make(chan state, capacity)
		chans = append(chans, ch)
	}

	for i := 0; i < nChans; i++ {
		for j := 0; j < i; j++ {
			chans[i] <- state{id: id(j)}
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	drainer := newChanDrainer(1*time.Millisecond, chans...)
	done := make(chan struct{})
	go func() {
		drainer(ctx)
		close(done)
	}()
	// Hack as the drainer should have run at least once when this finishes.
	time.Sleep(5 * time.Millisecond)
	cancel()
	<-done

	exp := make([]chan state, 0, nChans)
	defer func() {
		for _, ch := range exp {
			close(ch)
		}
	}()
	for i := 0; i < nChans; i++ {
		ch := make(chan state, capacity)
		exp = append(exp, ch)
	}
	for i := 0; i < nChans; i++ {
		max := i
		if i >= capacity*7/10 {
			max = 0
		}
		for j := 0; j < max; j++ {
			exp[i] <- state{id: id(j)}
		}
	}
	for i, ch := range chans {
		expCh := exp[i]
		if diff := cmp.Diff(len(expCh), len(ch)); diff != "" {
			t.Fatalf("channels length %d mismatch:\n%s", i, diff)
		}
		for j := 0; j < len(ch); j++ {
			expSt := <-expCh
			chSt := <-ch
			if diff := cmp.Diff(expSt.id, chSt.id); diff != "" {
				t.Fatalf("channels %d mismatch at state %d:\n%s", i, j, diff)
			}
		}
	}
}
