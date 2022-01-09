package flow

import (
	"github.com/DuarteMRAlves/maestro/internal/link"
	"github.com/DuarteMRAlves/maestro/internal/queue"
)

// Flow is a connection between two stages where data is transferred.
type Flow struct {
	link  *link.Link
	queue queue.Ring
}

func newFlow(l *link.Link) (*Flow, error) {
	q, err := queue.NewRing(1)
	if err != nil {
		return nil, err
	}
	f := &Flow{
		link:  l,
		queue: q,
	}
	return f, nil
}
