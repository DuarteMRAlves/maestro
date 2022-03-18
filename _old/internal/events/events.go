package events

import (
	"time"
)

// Subscription allows for a subscriber to see past events and wait for future
// events sent by a PubSub system.
type Subscription struct {
	// Token identifies the subscription. This token is used to unsubscribe.
	Token SubscriptionToken
	// Hist stores past events.
	Hist []*Event
	// Future is a channel where next events will be sent to.
	Future <-chan *Event
}

// SubscriptionToken uniquely identifies a Subscription.
type SubscriptionToken int

// Event is a report of an event that happened withing the maestro server. The
// event contains informative description that can be displayed to a user.
type Event struct {
	// Description is a short human-readable description for a event.
	Description string
	// Timestamp reports the timestamp at which the event was registered.
	Timestamp time.Time
}
