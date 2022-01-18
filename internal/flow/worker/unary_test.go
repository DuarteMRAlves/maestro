package worker

import (
	"github.com/DuarteMRAlves/maestro/internal/flow/state"
	"github.com/DuarteMRAlves/maestro/internal/reflection"
	"github.com/DuarteMRAlves/maestro/internal/testutil"
	mockflow "github.com/DuarteMRAlves/maestro/internal/testutil/mock/flow"
	mockreflection "github.com/DuarteMRAlves/maestro/internal/testutil/mock/reflection"
	"github.com/DuarteMRAlves/maestro/tests/pb"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"gotest.tools/v3/assert"
	"reflect"
	"testing"
)

func TestUnaryWorker_Run(t *testing.T) {
	lis := testutil.ListenAvailablePort(t)
	addr := lis.Addr().String()
	server := testutil.StartTestServer(t, lis, true, true)
	defer server.Stop()

	rpc := &mockreflection.RPC{
		Name_: "Unary",
		FQN:   "pb.TestService/Unary",
		In:    requestMessage(t),
		Out:   replyMessage(t),
		Unary: true,
	}
	msgs := []*pb.Request{
		{
			StringField:   "string-1",
			RepeatedField: []int64{1, 2, 3, 4},
			RepeatedInnerMsg: []*pb.InnerMessage{
				{
					RepeatedString: []string{"hello", "world", "1"},
				},
				{
					RepeatedString: []string{"other", "message", "2"},
				},
			},
		},
		{
			StringField:   "string-2",
			RepeatedField: []int64{1, 2, 3, 4},
			RepeatedInnerMsg: []*pb.InnerMessage{
				{
					RepeatedString: []string{"hello", "world", "2"},
				},
				{
					RepeatedString: []string{"other", "message", "2"},
				},
			},
		},
	}
	states := []*state.State{state.New(1, msgs[0]), state.New(3, msgs[1])}
	input := mockflow.NewInput(states)
	output := mockflow.NewOutput()
	done := make(chan bool)

	cfg := &Cfg{
		Address: addr,
		Rpc:     rpc,
		Input:   input,
		Output:  output,
		Done:    done,
	}

	w, err := NewWorker(cfg)
	assert.NilError(t, err, "create worker error")

	go w.Run()

	<-done

	assert.Equal(t, len(states), len(output.States), "correct number of replies")
	for i, in := range states {
		out := output.States[i]
		assert.Equal(t, in.Id(), out.Id(), "correct received id")

		req, ok := in.Msg().(*pb.Request)
		assert.Assert(t, ok, "request type assertion")

		dynRep, ok := out.Msg().(*dynamic.Message)
		assert.Assert(t, ok, "reply type assertion")
		rep := &pb.Reply{}
		err = dynRep.ConvertTo(rep)
		assert.NilError(t, err, "convert dynamic to Reply")

		testutil.AssertUnaryRequest(t, req, rep)
	}
}

func requestMessage(t *testing.T) reflection.Message {
	reqType := reflect.TypeOf(pb.Request{})

	reqDesc, err := desc.LoadMessageDescriptorForType(reqType)
	assert.NilError(t, err, "load desc Request")

	msg, err := reflection.NewMessage(reqDesc)
	assert.NilError(t, err, "message Request")

	return msg
}

func replyMessage(t *testing.T) reflection.Message {
	repType := reflect.TypeOf(pb.Reply{})

	repDesc, err := desc.LoadMessageDescriptorForType(repType)
	assert.NilError(t, err, "load desc Reply")

	msg, err := reflection.NewMessage(repDesc)
	assert.NilError(t, err, "message Reply")

	return msg
}
