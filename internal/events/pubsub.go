package events

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"sync"
	"time"
)

type GenToken func() api.SubscriptionToken
type CreateChan func() chan *api.Event
type SendEvent func(chan<- *api.Event, *api.Event)

// PubSub handles the distribution of events for multiple subscribers with
// multiple publishers.
type PubSub interface {
	// Subscribe returns a new subscription with past events and a channel to
	// listen to new events.
	Subscribe() *api.Subscription
	// Unsubscribe stops sending new events for a subscription, closing the
	// respective channel.
	Unsubscribe(api.SubscriptionToken) error
	// Publish publishes a new event with the received description and the
	// current timestamp.
	Publish(string)
	// Close shuts down the PubSub system, closing all channels.
	Close()
}

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
type pubSub struct {
	genToken   GenToken
	createChan CreateChan
	sendEvent  SendEvent
	// hist stores past events such that new subscribers can retrieve them.
	hist []*api.Event
	// subs are the channels used to send messages to the subscribers.
	subs map[api.SubscriptionToken]chan<- *api.Event

	mu sync.Mutex
}

func NewPubSub(ctx PubSubContext) PubSub {
	var sendEvent SendEvent

	createChan := func() chan *api.Event {
		return make(chan *api.Event, ctx.BuffSize)
	}
	if ctx.Timeout > 0 {
		sendEvent = SendWithTimeout(ctx.Timeout)
	} else {
		sendEvent = SendWithoutTimeout()
	}
	return &pubSub{
		genToken:   IncrementalGenToken(0),
		createChan: createChan,
		sendEvent:  sendEvent,
		hist:       make([]*api.Event, 0),
		subs:       make(map[api.SubscriptionToken]chan<- *api.Event, 0),
	}
}

func (pb *pubSub) Subscribe() *api.Subscription {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	hist := make([]*api.Event, 0, len(pb.hist))
	for _, h := range pb.hist {
		event := &api.Event{}
		copyEvent(event, h)
		hist = append(hist, event)
	}

	token := pb.genToken()
	sub := pb.createChan()
	pb.subs[token] = sub

	return &api.Subscription{
		Token:  token,
		Hist:   hist,
		Future: sub,
	}
}

func (pb *pubSub) Unsubscribe(token api.SubscriptionToken) error {
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

func (pb *pubSub) Publish(description string) {
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

func SendWithTimeout(timeout time.Duration) SendEvent {
	return func(sub chan<- *api.Event, event *api.Event) {
		timer := time.NewTimer(timeout)
		defer timer.Stop()
		select {
		case sub <- event:
		case <-timer.C:
		}
	}
}

func SendWithoutTimeout() SendEvent {
	return func(sub chan<- *api.Event, event *api.Event) {
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

func (pb *pubSub) Close() {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	for _, sub := range pb.subs {
		close(sub)
	}
}

func IncrementalGenToken(start api.SubscriptionToken) GenToken {
	c := tokenCounter{curr: start}
	return c.Next
}

type tokenCounter struct {
	curr api.SubscriptionToken
}

func (c *tokenCounter) Next() api.SubscriptionToken {
	t := c.curr
	c.curr++
	return t
}
