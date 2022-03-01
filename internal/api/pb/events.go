package pb

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/events"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MarshalEvent(pbEvent *pb.Event, event *events.Event) {
	pbEvent.Description = event.Description
	pbEvent.Timestamp = timestamppb.New(event.Timestamp)
}
