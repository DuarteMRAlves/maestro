package execute

import (
	"fmt"

	"github.com/DuarteMRAlves/maestro/internal/message"
)

// state defines a structure to store the state of an pipeline.
type state struct {
	msg message.Instance
}

func newState(msg message.Instance) state {
	return state{msg: msg}
}

func (s state) String() string {
	return fmt.Sprintf("state{msg:%v}", s.msg)
}
