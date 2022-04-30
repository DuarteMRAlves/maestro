package execute

import (
	"fmt"

	"github.com/DuarteMRAlves/maestro/internal"
)

// offlineState defines a structure to store the state of an offline pipeline.
type offlineState struct {
	msg internal.Message
}

func newOfflineState(msg internal.Message) offlineState {
	return offlineState{msg: msg}
}

func (s offlineState) String() string {
	return fmt.Sprintf("offlineState{msg:%v}", s.msg)
}
