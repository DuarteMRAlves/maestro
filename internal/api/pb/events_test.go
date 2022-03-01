package pb

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/events"
	"gotest.tools/v3/assert"
	"testing"
	"time"
)

func TestEvent(t *testing.T) {
	var marshalled pb.Event

	orig := &events.Event{
		Description: "Event Description",
		Timestamp:   time.Now(),
	}

	MarshalEvent(&marshalled, orig)
	assert.Equal(t, orig.Description, marshalled.Description)
	assert.Assert(t, orig.Timestamp.Equal(marshalled.Timestamp.AsTime()))
}
