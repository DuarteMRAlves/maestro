package events

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
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
	ctx   PubSubContext
	token api.SubscriptionToken
	// hist stores past events such that new subscribers can retrieve them.
	hist []*api.Event
	// subs are the channels used to send messages to the subscribers.
	subs map[api.SubscriptionToken]chan<- *api.Event

	mu sync.Mutex
}

func NewPubSub(ctx PubSubContext) *PubSub {
	return &PubSub{
		ctx:   ctx,
		token: 0,
		hist:  make([]*api.Event, 0),
		subs:  make(map[api.SubscriptionToken]chan<- *api.Event, 0),
		mu:    sync.Mutex{},
	}
}

func (pb *PubSub) Subscribe() *api.Subscription {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	hist := make([]*api.Event, 0, len(pb.hist))
	for _, h := range pb.hist {
		event := &api.Event{}
		copyEvent(event, h)
		hist = append(hist, event)
	}

	token := pb.token
	pb.token++

	sub := make(chan *api.Event, pb.ctx.BuffSize)
	pb.subs[token] = sub

	return &api.Subscription{
		Token:  token,
		Hist:   hist,
		Future: sub,
	}
}

func (pb *PubSub) Unsubscribe(token api.SubscriptionToken) error {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	future, exists := pb.subs[token]
	if !exists {
		return errdefs.NotFoundWithMsg("Token not found: %d", token)
	}
	close(future)
	delete(pb.subs, token)
	return nil
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

func (pb *PubSub) Close() {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	for _, sub := range pb.subs {
		close(sub)
	}
}
