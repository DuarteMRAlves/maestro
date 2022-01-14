package state

// State defines a structure to manage the flow of a computation. This State is
// created in a source stage with a unique id, that is transferred through the
// orchestration, identifying the flow so that messages in parallel branches can
// be synchronized.
type State struct {
	id  Id
	msg interface{}
}

// Id is a unique id that identifies a computation flow.
type Id int

func New(id Id, msg interface{}) *State {
	return &State{
		id:  id,
		msg: msg,
	}
}

func (f *State) Id() Id {
	return f.id
}

func (f *State) Msg() interface{} {
	return f.msg
}

func (f *State) Update(msg interface{}) {
	f.msg = msg
}
