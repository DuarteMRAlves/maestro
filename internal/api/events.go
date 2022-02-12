package api

import "time"

// Event is a report of an event that happened withing the maestro server. The
// event contains informative description that can be displayed to a user.
type Event struct {
	// Description is a short human-readable description for a event.
	Description string
	// Timestamp reports the timestamp at which the event was registered.
	Timestamp time.Time
}
