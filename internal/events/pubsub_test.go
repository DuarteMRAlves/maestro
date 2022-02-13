package events

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"gotest.tools/v3/assert"
	"testing"
	"time"
)

func TestPubSub_Publish(t *testing.T) {
	numEvents := 20
	prevTimestamp := time.Now()

	pubSub := NewPubSub(DefaultPubSubContext())
	hist1, sub1 := pubSub.RegisterSub()
	assert.Equal(t, 0, len(hist1))
	hist2, sub2 := pubSub.RegisterSub()
	assert.Equal(t, 0, len(hist2))

	go func() {
		for i := 0; i < numEvents; i++ {
			pubSub.Publish(fmt.Sprintf("Event-%d", i+1))
		}
	}()

	for i := 0; i < numEvents; i++ {
		expected := fmt.Sprintf("Event-%d", i+1)
		event1 := <-sub1
		event2 := <-sub2

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
	hist1, sub1 := pubSub.RegisterSub()
	assert.Equal(t, 0, len(hist1))

	go func() {
		for i := 0; i < firstEvents; i++ {
			pubSub.Publish(fmt.Sprintf("Event-First-%d", i+1))
		}
	}()

	collected := make([]*api.Event, 0, firstEvents)

	for i := 0; i < firstEvents; i++ {
		expected := fmt.Sprintf("Event-First-%d", i+1)
		event1 := <-sub1

		assert.Equal(t, expected, event1.Description)

		assert.Assert(t, firstTimestamp.Before(event1.Timestamp))
		firstTimestamp = event1.Timestamp

		collected = append(collected, event1)
	}

	hist2, sub2 := pubSub.RegisterSub()
	assert.Equal(t, firstEvents, len(hist2))
	for i := 0; i < firstEvents; i++ {
		hist := hist2[i]
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
		event1 := <-sub1
		event2 := <-sub2

		assert.Equal(t, expected, event1.Description)
		assert.Equal(t, expected, event2.Description)

		assert.Equal(t, event1.Timestamp, event2.Timestamp)
		assert.Assert(t, secondTimestamp.Before(event1.Timestamp))
		secondTimestamp = event1.Timestamp
	}

	pubSub.Close()
}
