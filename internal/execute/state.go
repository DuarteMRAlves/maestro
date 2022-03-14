package execute

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal"
)

type id int

// State defines a structure to manage the flow of a computation. This State is
// created in a source stage with a unique id, that is transferred through the
// orchestration, identifying the flow so that messages in parallel branches can
// be synchronized.
type state struct {
	id  id
	msg internal.Message
}

var emptyState = newState(-1, nil)

func newState(id id, msg internal.Message) state {
	return state{
		id:  id,
		msg: msg,
	}
}

func updateStateMsg(s state, msg internal.Message) state {
	return newState(s.id, msg)
}

func (s state) String() string {
	return fmt.Sprintf("state{id:%d,msg:%v}", s.id, s.msg)
}
