package mock

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal"
)

type MessageDescriptor struct {
	Ident           string
	MatchCompatible []string
	Fields          map[internal.MessageField]internal.MessageDesc
}

func (d MessageDescriptor) Compatible(other internal.MessageDesc) bool {
	otherMock, ok := other.(MessageDescriptor)
	if !ok {
		return false
	}
	if d.Ident == otherMock.Ident {
		return true
	}
	for _, c := range d.MatchCompatible {
		if c == otherMock.Ident {
			return true
		}
	}
	return false
}

func (d MessageDescriptor) EmptyGen() internal.EmptyMessageGen {
	return func() internal.Message {
		return &Message{Fields: map[internal.MessageField]interface{}{}}
	}
}

func (d MessageDescriptor) GetField(field internal.MessageField) (
	internal.MessageDesc,
	error,
) {
	inner, ok := d.Fields[field]
	if !ok {
		return nil, fmt.Errorf("field not found %s", field)
	}
	return inner, nil
}
