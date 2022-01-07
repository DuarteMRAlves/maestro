package flow

// Flow defines a structure to manage the flow of a computation. This Flow is
// created in a source stage with a unique id, that is transferred through the
// orchestration, identifying the flow so that messages in parallel branches can
// be synchronized.
type Flow struct {
	id  Id
	msg interface{}
}

// Id is a unique id that identifies a computation flow.
type Id int

func New(id Id, msg interface{}) *Flow {
	return &Flow{
		id:  id,
		msg: msg,
	}
}

func (f *Flow) Id() Id {
	return f.id
}

func (f *Flow) Msg() interface{} {
	return f.msg
}

func (f *Flow) Update(msg interface{}) {
	f.msg = msg
}
