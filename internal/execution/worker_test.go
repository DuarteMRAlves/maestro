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

func TestUnaryWorker_RunAndEOF(t *testing.T) {
	lis := util.NewTestListener(t)
	addr := lis.Addr().String()
	server := util.StartTestServer(t, lis, true, true)
	defer server.Stop()

	rpc := &rpc.MockRPC{
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
	states := []*State{
		NewState(1, msgs[0]),
		NewState(3, msgs[1]),
	}
	input := NewMockInput(append(states, NewEOFState(4)), func() {})
	output := NewMockOutput()
	done := make(chan bool)

	cfg := &WorkerCfg{
		Address: addr,
		Rpc:     rpc,
		Input:   input,
		Output:  output,
		Done:    done,
	}

	w, err := NewWorker(cfg)
	assert.NilError(t, err, "create worker error")

	term := make(chan struct{})
	defer close(term)
	go w.Run(term)

	<-done
	input.Close()

	assert.Equal(
		t,
		len(states),
		len(output.States),
		"correct number of replies",
	)
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

		util.AssertUnaryRequest(t, req, rep)
	}
}

func TestUnaryWorker_RunAndCtxDone(t *testing.T) {
	lis := util.NewTestListener(t)
	addr := lis.Addr().String()
	server := util.StartTestServer(t, lis, true, true)
	defer server.Stop()

	rpc := &rpc.MockRPC{
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
	states := []*State{NewState(1, msgs[0]), NewState(3, msgs[1])}

	inputSync := make(chan bool)
	defer close(inputSync)

	input := NewMockInput(
		states,
		func() {
			inputSync <- true
		},
	)
	output := NewMockOutput()
	done := make(chan bool)
	defer close(done)

	cfg := &WorkerCfg{
		Address: addr,
		Rpc:     rpc,
		Input:   input,
		Output:  output,
		Done:    done,
	}

	w, err := NewWorker(cfg)
	assert.NilError(t, err, "create worker error")
	term := make(chan struct{})
	defer close(term)
	go w.Run(term)

	// Wait for last input to be sent. Now the input function will block.
	<-inputSync

	term <- struct{}{}

	<-done
	// Release input
	input.Close()

	assert.Equal(
		t,
		len(states),
		len(output.States),
		"correct number of replies",
	)
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

		util.AssertUnaryRequest(t, req, rep)
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
