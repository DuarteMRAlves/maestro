package events

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
	"sync"
	"time"
)

// PubSubContext specifies configurations for PubSub.
type PubSubContext struct {
	// Timeout is the timeout that the PubSub system waits when sending a
	// message.
	Timeout time.Duration
	// BuffSize specifies the size of the creates sub channels.
	BuffSize int
}

func DefaultPubSubContext() PubSubContext {
	return PubSubContext{Timeout: time.Second, BuffSize: 10}
}

// PubSub handles the distribution of events for multiple subscribers with
// multiple publishers.
type PubSub struct {
	ctx PubSubContext
	// hist stores past events such that new subscribers can retrieve them.
	hist []*api.Event
	// subs are the channels used to send messages to the subscribers.
	subs []chan<- *api.Event

	mu sync.Mutex
}

func NewPubSub(ctx PubSubContext) *PubSub {
	return &PubSub{
		ctx:  ctx,
		hist: make([]*api.Event, 0),
		subs: make([]chan<- *api.Event, 0),
		mu:   sync.Mutex{},
	}
}

func (pb *PubSub) RegisterSub() ([]*api.Event, <-chan *api.Event) {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	hist := make([]*api.Event, 0, len(pb.hist))
	for _, h := range pb.hist {
		event := &api.Event{}
		copyEvent(event, h)
		hist = append(hist, event)
	}

	sub := make(chan *api.Event, pb.ctx.BuffSize)
	pb.subs = append(pb.subs, sub)

	return hist, sub
}

func (pb *PubSub) Publish(description string) {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	ts := time.Now()
	event := &api.Event{
		Description: description,
		Timestamp:   ts,
	}

	pb.hist = append(pb.hist, event)
	for _, sub := range pb.subs {
		send := &api.Event{}
		copyEvent(send, event)
		pb.sendEvent(sub, send)
	}
}

func (pb *PubSub) sendEvent(sub chan<- *api.Event, event *api.Event) {
	if pb.ctx.Timeout > 0 {
		timer := time.NewTimer(pb.ctx.Timeout)
		defer timer.Stop()
		select {
		case sub <- event:
		case <-timer.C:
		}
	} else {
		select {
		case sub <- event:
		default:
			// No timeout
		}
	}
}

func copyEvent(dst *api.Event, src *api.Event) {
	dst.Description = src.Description
	dst.Timestamp = src.Timestamp
}
