package execution

import (
	"github.com/DuarteMRAlves/maestro/internal/rpc"
	"github.com/DuarteMRAlves/maestro/tests/pb"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"gotest.tools/v3/assert"
	"reflect"
	"testing"
)

func TestSourceInput_Next(t *testing.T) {
	var (
		state *State
		err   error
	)

	reqType := reflect.TypeOf(pb.Request{})

	reqDesc, err := desc.LoadMessageDescriptorForType(reqType)
	assert.NilError(t, err, "load desc Request")

	msg, err := rpc.NewMessage(reqDesc)
	assert.NilError(t, err, "message Request")

	input, err := NewInputBuilder().WithMessage(msg).Build()
	assert.NilError(t, err, "build input")

	sourceInput, ok := input.(*SourceInput)
	assert.Assert(t, ok, "source input cast")

	for i := 1; i < 10; i++ {
		state = <-sourceInput.Chan()
		assert.NilError(t, err, "next at iter %d", i)
		assert.Equal(t, Id(i), state.Id(), "id at iter %d", i)
		opt := cmpopts.IgnoreUnexported(dynamic.Message{})
		assert.DeepEqual(t, msg.NewEmpty(), state.Msg(), opt)
	}
}
