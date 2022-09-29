package execute

import (
	"fmt"

	"github.com/DuarteMRAlves/maestro/internal/message"
)

// offlineState defines a structure to store the state of an offline pipeline.
type offlineState struct {
	msg message.Instance
}

func newOfflineState(msg message.Instance) offlineState {
	return offlineState{msg: msg}
}

func (s offlineState) String() string {
	return fmt.Sprintf("offlineState{msg:%v}", s.msg)
}
