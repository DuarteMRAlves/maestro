package exec

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/rpc"
	"io"
)

// State defines a structure to manage the flow of a computation. This State is
// created in a source stage with a unique id, that is transferred through the
// orchestration, identifying the flow so that messages in parallel branches can
// be synchronized.
type State struct {
	id  Id
	msg rpc.DynMessage
	err error
}

// Id is a unique id that identifies a computation flow.
type Id int

func NewState(id Id, msg rpc.DynMessage) *State {
	return &State{
		id:  id,
		msg: msg,
	}
}

func NewEOFState(id Id) *State {
	return &State{
		id:  id,
		msg: nil,
		err: io.EOF,
	}
}

func (s *State) Id() Id {
	return s.id
}

func (s *State) Msg() rpc.DynMessage {
	return s.msg
}

func (s *State) Err() error {
	return s.err
}

func (s *State) Update(msg rpc.DynMessage) {
	s.msg = msg
}

func (s *State) String() string {
	return fmt.Sprintf("State{Id:%d,Msg:%v,Err:%v}", s.id, s.msg, s.err)
}
