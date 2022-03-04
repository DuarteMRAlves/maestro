package events

import (
	"fmt"
	"gotest.tools/v3/assert"
	"testing"
	"time"
)

func TestPubSub_Publish(t *testing.T) {
	numEvents := 20
	prevTimestamp := time.Now()

	pubSub := NewPubSub(DefaultPubSubContext())
	sub1 := pubSub.Subscribe()
	assert.Equal(t, 0, len(sub1.Hist))
	sub2 := pubSub.Subscribe()
	assert.Equal(t, 0, len(sub2.Hist))

	go func() {
		for i := 0; i < numEvents; i++ {
			pubSub.Publish(fmt.Sprintf("Event-%d", i+1))
		}
	}()

	for i := 0; i < numEvents; i++ {
		expected := fmt.Sprintf("Event-%d", i+1)
		event1 := <-sub1.Future
		event2 := <-sub2.Future

		assert.Equal(t, expected, event1.Description)
		assert.Equal(t, expected, event2.Description)

		assert.Equal(t, event1.Timestamp, event2.Timestamp)
		assert.Assert(t, prevTimestamp.Before(event1.Timestamp))
		prevTimestamp = event1.Timestamp
	}

	pubSub.Close()
}

func TestPubSub_Publish_History(t *testing.T) {
	firstEvents := 20
	firstTimestamp := time.Now()

	pubSub := NewPubSub(DefaultPubSubContext())
	sub1 := pubSub.Subscribe()
	assert.Equal(t, 0, len(sub1.Hist))

	go func() {
		for i := 0; i < firstEvents; i++ {
			pubSub.Publish(fmt.Sprintf("Event-First-%d", i+1))
		}
	}()

	collected := make([]*Event, 0, firstEvents)

	for i := 0; i < firstEvents; i++ {
		expected := fmt.Sprintf("Event-First-%d", i+1)
		event1 := <-sub1.Future

		assert.Equal(t, expected, event1.Description)

		assert.Assert(t, firstTimestamp.Before(event1.Timestamp))
		firstTimestamp = event1.Timestamp

		collected = append(collected, event1)
	}

	sub2 := pubSub.Subscribe()
	assert.Equal(t, firstEvents, len(sub2.Hist))
	for i := 0; i < firstEvents; i++ {
		hist := sub2.Hist[i]
		coll := collected[i]

		assert.Equal(t, coll.Description, hist.Description)
		assert.Assert(t, hist.Timestamp.Equal(coll.Timestamp))
	}

	secondEvents := 30
	secondTimestamp := time.Now()

	go func() {
		for i := 0; i < secondEvents; i++ {
			pubSub.Publish(fmt.Sprintf("Second-Event-%d", i+1))
		}
	}()

	for i := 0; i < secondEvents; i++ {
		expected := fmt.Sprintf("Second-Event-%d", i+1)
		event1 := <-sub1.Future
		event2 := <-sub2.Future

		assert.Equal(t, expected, event1.Description)
		assert.Equal(t, expected, event2.Description)

		assert.Equal(t, event1.Timestamp, event2.Timestamp)
		assert.Assert(t, secondTimestamp.Before(event1.Timestamp))
		secondTimestamp = event1.Timestamp
	}

	pubSub.Close()
}

func TestPubSub_Unsubscribe(t *testing.T) {
	firstEvents := 20
	firstTimestamp := time.Now()

	pubSub := NewPubSub(DefaultPubSubContext())
	defer pubSub.Close()

	sub1 := pubSub.Subscribe()
	assert.Equal(t, 0, len(sub1.Hist))
	sub2 := pubSub.Subscribe()
	assert.Equal(t, 0, len(sub2.Hist))

	go func() {
		for i := 0; i < firstEvents; i++ {
			pubSub.Publish(fmt.Sprintf("Event-First-%d", i+1))
		}
	}()

	for i := 0; i < firstEvents; i++ {
		expected := fmt.Sprintf("Event-First-%d", i+1)
		event1 := <-sub1.Future
		event2 := <-sub2.Future

		assert.Equal(t, expected, event1.Description)
		assert.Equal(t, expected, event2.Description)

		assert.Equal(t, event1.Timestamp, event2.Timestamp)
		assert.Assert(t, firstTimestamp.Before(event1.Timestamp))
		firstTimestamp = event1.Timestamp
	}

	err := pubSub.Unsubscribe(sub1.Token)
	assert.NilError(t, err, "unsubscribe error")

	event1, open := <-sub1.Future
	assert.Assert(t, !open, "sub1 future is closed")
	assert.Assert(t, event1 == nil, "event is nil after closed")

	secondEvents := 30
	secondTimestamp := time.Now()

	go func() {
		for i := 0; i < secondEvents; i++ {
			pubSub.Publish(fmt.Sprintf("Second-Event-%d", i+1))
		}
	}()

	for i := 0; i < secondEvents; i++ {
		expected := fmt.Sprintf("Second-Event-%d", i+1)
		event2 := <-sub2.Future

		assert.Equal(t, expected, event2.Description)

		assert.Assert(t, secondTimestamp.Before(event2.Timestamp))
		secondTimestamp = event2.Timestamp
	}
}
