package execution

import (
	"github.com/DuarteMRAlves/maestro/internal/rpc"
	"github.com/DuarteMRAlves/maestro/internal/util"
	"github.com/DuarteMRAlves/maestro/tests/pb"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"gotest.tools/v3/assert"
	"reflect"
	"testing"
)

func TestUnaryStage_RunAndEOF(t *testing.T) {
	lis := util.NewTestListener(t)
	addr := lis.Addr().String()
	server := util.StartTestServer(t, lis, true, true)
	defer server.Stop()

	rpc := testRpc(t)
	msgs := testRequests()

	states := []*State{
		NewState(1, msgs[0]),
		NewState(3, msgs[1]),
	}
	input := make(chan *State, len(states)+1)
	input <- states[0]
	input <- states[1]
	input <- NewEOFState(4)

	output := make(chan *State, len(states))

	cfg := &StageCfg{
		Address: addr,
		Rpc:     rpc,
		Input:   input,
		Output:  output,
	}

	s, err := NewStage(cfg)
	assert.NilError(t, err, "create stage error")

	term := make(chan struct{})
	errs := make(chan error)
	done := make(chan struct{})
	runCfg := &RunCfg{
		term: term,
		errs: errs,
		done: done,
	}
	defer close(term)
	defer close(errs)
	go s.Run(runCfg)

	<-done
	close(input)
	close(output)

	rcvStates := collectState(output)

	assert.Equal(t, len(states), len(rcvStates), "correct number of replies")
	for i, in := range states {
		out := rcvStates[i]
		assert.Equal(t, in.Id(), out.Id(), "correct received id")

		req, ok := in.Msg().(*pb.Request)
		assert.Assert(t, ok, "request type assertion")

		dynRep, ok := out.Msg().(*dynamic.Message)
		assert.Assert(t, ok, "reply type assertion")
		rep := &pb.Reply{}
		err = dynRep.ConvertTo(rep)
		assert.NilError(t, err, "convert dynamic to Reply")

		util.AssertUnaryRequest(t, req, rep)
	}
	assert.Equal(t, 0, len(errs), "No errors")
}

func collectState(output <-chan *State) []*State {
	rcvStates := make([]*State, 0)
	collect := make(chan struct{})
	go func() {
		for s := range output {
			rcvStates = append(rcvStates, s)
		}
		collect <- struct{}{}
	}()
	<-collect
	return rcvStates
}

func testRpc(t *testing.T) *rpc.MockRPC {
	return &rpc.MockRPC{
		Name_:  "Unary",
		FQN:    "pb.TestService/Unary",
		Invoke: "pb.TestService/Unary",
		In:     requestMessage(t),
		Out:    replyMessage(t),
		Unary:  true,
	}
}

func testRequests() []*pb.Request {
	return []*pb.Request{
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
}

func requestMessage(t *testing.T) rpc.Message {
	reqType := reflect.TypeOf(pb.Request{})

	reqDesc, err := desc.LoadMessageDescriptorForType(reqType)
	assert.NilError(t, err, "load desc Request")

	msg, err := rpc.NewMessage(reqDesc)
	assert.NilError(t, err, "message Request")

	return msg
}

func replyMessage(t *testing.T) rpc.Message {
	repType := reflect.TypeOf(pb.Reply{})

	repDesc, err := desc.LoadMessageDescriptorForType(repType)
	assert.NilError(t, err, "load desc Reply")

	msg, err := rpc.NewMessage(repDesc)
	assert.NilError(t, err, "message Reply")

	return msg
}
