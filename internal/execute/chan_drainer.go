package execute

import (
	"context"
	"time"
)

func newChanDrainer(sleep time.Duration, channels ...chan state) func(context.Context) {
	highs := make([]int, 0, len(channels))
	for _, ch := range channels {
		highs = append(highs, cap(ch)*7/10)
	}
	return func(ctx context.Context) {
		for {
			for i, ch := range channels {
				if len(ch) >= highs[i] {
					drainChan(ch)
				}
			}
			timer := time.NewTimer(sleep)
			select {
			case <-timer.C:
			case <-ctx.Done():
				return
			}
		}
	}
}

func drainChan(ch chan state) {
	maxDrain := len(ch)
	for i := 0; i < maxDrain; i++ {
		select {
		case <-ch:
		default:
			return
		}
	}
}
