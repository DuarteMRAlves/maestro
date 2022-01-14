package flow

import (
	"github.com/DuarteMRAlves/maestro/internal/link"
	"github.com/DuarteMRAlves/maestro/internal/queue"
)

// Flow is a connection between two stages where data is transferred.
type Flow struct {
	Link  *link.Link
	Queue queue.Ring
}

func NewFlow(l *link.Link) (*Flow, error) {
	q, err := queue.NewRing(1)
	if err != nil {
		return nil, err
	}
	f := &Flow{
		Link:  l,
		Queue: q,
	}
	return f, nil
}
