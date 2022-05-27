package retry

import (
	"time"
)

type Retryable func() bool

type BackoffStrategy interface {
	Next() time.Duration
}

func WhileTrue(f Retryable, strat BackoffStrategy) {
	for {
		retry := f()
		if !retry {
			return
		}
		time.Sleep(strat.Next())
	}
}

const (
	defaultInitBackoff = 1 * time.Millisecond
	defaultFact        = 2
)

type ExponentialBackoff struct {
	curr time.Duration
	fact int
}

func NewExponentialBackoff(initBackoff time.Duration, fact int) *ExponentialBackoff {
	return &ExponentialBackoff{
		curr: initBackoff,
		fact: fact,
	}
}

func (b *ExponentialBackoff) Next() time.Duration {
	if b.curr == 0 {
		b.curr = defaultInitBackoff
	}
	if b.fact == 0 {
		b.fact = defaultFact
	}
	backoff := b.curr
	b.curr = b.curr * time.Duration(b.fact)
	return backoff
}
