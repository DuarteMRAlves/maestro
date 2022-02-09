package exec

import (
	"github.com/DuarteMRAlves/maestro/internal/rpc"
	"github.com/DuarteMRAlves/maestro/internal/util"
	"github.com/DuarteMRAlves/maestro/tests/pb"
	"github.com/google/go-cmp/cmp/cmpopts"
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

func TestSourceStage_Run(t *testing.T) {
	var (
		state *State
		err   error
	)

	reqType := reflect.TypeOf(pb.Request{})

	reqDesc, err := desc.LoadMessageDescriptorForType(reqType)
	assert.NilError(t, err, "load desc Request")

	msg, err := rpc.NewMessage(reqDesc)
	assert.NilError(t, err, "message Request")

	ch := make(chan *State)
	st := NewSourceStage(1, ch, msg)
	assert.NilError(t, err, "build input")

	term := make(chan struct{})
	done := make(chan struct{})
	errs := make(chan error)

	go st.Run(&RunCfg{term: term, done: done, errs: errs})

	for i := 1; i < 10; i++ {
		state = <-ch
		assert.NilError(t, err, "next at iter %d", i)
		assert.Equal(t, Id(i), state.Id(), "id at iter %d", i)
		opt := cmpopts.IgnoreUnexported(dynamic.Message{})
		assert.DeepEqual(t, msg.NewEmpty(), state.Msg(), opt)
	}

	close(term)
	<-done
	close(errs)
	assert.Assert(t, len(errs) == 0)
}

func TestMergeStage_Run(t *testing.T) {
	fields := []string{"in1", "in2", "in3"}

	input1 := make(chan *State)
	input2 := make(chan *State)
	input3 := make(chan *State)
	inputs := []<-chan *State{input1, input2, input3}

	output := make(chan *State)

	outType := reflect.TypeOf(pb.MergeMessage{})
	outDesc, err := desc.LoadMessageDescriptorForType(outType)
	assert.NilError(t, err, "load desc MergeMessage")
	msg, err := rpc.NewMessage(outDesc)
	assert.NilError(t, err, "create message MergeMessage")

	s := NewMergeStage(fields, inputs, output, msg)

	term := make(chan struct{})
	done := make(chan struct{})
	errs := make(chan error)

	expected := []*State{
		NewState(3, testMergeMessage(3)),
		NewState(6, testMergeMessage(6)),
	}

	go func() {
		input1 <- NewState(1, &pb.MergeInner1{Val: 1})
		input1 <- NewState(2, &pb.MergeInner1{Val: 2})
		input1 <- NewState(3, &pb.MergeInner1{Val: 3})
		input1 <- NewState(6, &pb.MergeInner1{Val: 6})
	}()

	go func() {
		input2 <- NewState(2, &pb.MergeInner2{Val: 2})
		input2 <- NewState(3, &pb.MergeInner2{Val: 3})
		input2 <- NewState(5, &pb.MergeInner2{Val: 5})
		input2 <- NewState(6, &pb.MergeInner2{Val: 6})
	}()

	go func() {
		input3 <- NewState(1, &pb.MergeInner3{Val: 2})
		input3 <- NewState(3, &pb.MergeInner3{Val: 3})
		input3 <- NewState(5, &pb.MergeInner3{Val: 5})
		input3 <- NewState(6, &pb.MergeInner3{Val: 6})
	}()

	go s.Run(&RunCfg{term: term, done: done, errs: errs})

	for i, exp := range expected {
		out := <-output
		assert.NilError(t, err, "next at iter %d", i)
		assert.Equal(t, exp.id, out.Id(), "id at iter %d", i)
		expMsg, ok := exp.Msg().(*pb.MergeMessage)
		assert.Assert(t, ok, "cast for exp at iter %d", i)
		dynMsg, ok := out.Msg().(*dynamic.Message)
		assert.Assert(t, ok, "cast for out at iter %d", i)
		outMsg := &pb.MergeMessage{}
		err = dynMsg.ConvertTo(outMsg)
		assert.NilError(t, err, "convert dyn to out")
		assert.Equal(t, expMsg.In1.Val, outMsg.In1.Val)
		assert.Equal(t, expMsg.In2.Val, outMsg.In2.Val)
		assert.Equal(t, expMsg.In3.Val, outMsg.In3.Val)
	}

	close(term)
	<-done
	close(errs)
	assert.Assert(t, len(errs) == 0)
}

func testMergeMessage(val int32) *pb.MergeMessage {
	return &pb.MergeMessage{
		In1: &pb.MergeInner1{Val: val},
		In2: &pb.MergeInner2{Val: val},
		In3: &pb.MergeInner3{Val: val},
	}
}
