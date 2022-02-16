package exec

import (
	"context"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/events"
	"github.com/DuarteMRAlves/maestro/internal/rpc"
	"github.com/DuarteMRAlves/maestro/tests/pb"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gotest.tools/v3/assert"
	"net"
	"reflect"
	"testing"
)

func TestUnaryStage_RunAndEOF(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:0")
	assert.NilError(t, err, "failed to listen")
	addr := lis.Addr().String()
	server := startTestServer(t, lis, true, true)
	defer server.Stop()

	rpc := testRpc(t)
	msgs := testRequests(t)

	states := []*State{
		NewState(1, msgs[0]),
		NewState(3, msgs[1]),
	}
	input := make(chan *State, len(states)+1)
	input <- states[0]
	input <- states[1]
	input <- NewEOFState(4)

	output := make(chan *State, len(states))

	s, err := NewRpcStage(addr, rpc, input, output)
	assert.NilError(t, err, "create stage error")

	term := make(chan struct{})
	errs := make(chan error)
	done := make(chan struct{})
	runCfg := &RunCfg{
		pubSub: &events.MockPubSub{},
		term:   term,
		errs:   errs,
		done:   done,
	}
	defer close(term)
	defer close(errs)
	go s.Run(runCfg)

	<-done
	close(input)
	s.Close()

	rcvStates := collectState(output)

	assert.Equal(t, len(states), len(rcvStates), "correct number of replies")
	for i, in := range states {
		out := rcvStates[i]
		assert.Equal(t, in.Id(), out.Id(), "correct received id")

		dynReq, ok := in.Msg().GrpcMsg().(*dynamic.Message)
		assert.Assert(t, ok, "request type assertion")
		req := &pb.Request{}
		err = dynReq.ConvertTo(req)
		assert.NilError(t, err, "convert dynamic to Request")

		dynRep, ok := out.Msg().GrpcMsg().(*dynamic.Message)
		assert.Assert(t, ok, "reply type assertion")
		rep := &pb.Reply{}
		err = dynRep.ConvertTo(rep)
		assert.NilError(t, err, "convert dynamic to Reply")

		assertUnaryRequest(t, req, rep)
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

func testRequests(t *testing.T) []rpc.DynMessage {
	msg1, err := rpc.DynMessageFromProto(
		&pb.Request{
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
	)
	assert.NilError(t, err, "create message 1")
	msg2, err := rpc.DynMessageFromProto(
		&pb.Request{
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
	)
	assert.NilError(t, err, "create message 2")

	return []rpc.DynMessage{msg1, msg2}
}

func requestMessage(t *testing.T) rpc.MessageDesc {
	reqType := reflect.TypeOf(pb.Request{})

	reqDesc, err := desc.LoadMessageDescriptorForType(reqType)
	assert.NilError(t, err, "load desc Request")

	msg, err := rpc.NewMessage(reqDesc)
	assert.NilError(t, err, "message Request")

	return msg
}

func replyMessage(t *testing.T) rpc.MessageDesc {
	repType := reflect.TypeOf(pb.Reply{})

	repDesc, err := desc.LoadMessageDescriptorForType(repType)
	assert.NilError(t, err, "load desc Reply")

	msg, err := rpc.NewMessage(repDesc)
	assert.NilError(t, err, "message Reply")

	return msg
}

type testService struct {
	pb.UnimplementedTestServiceServer
	pb.UnimplementedExtraServiceServer
}

func (s *testService) Unary(
	ctx context.Context,
	request *pb.Request,
) (*pb.Reply, error) {

	if request.StringField == "error" {
		return nil, fmt.Errorf("dummy error")
	}
	return replyFromRequest(request), nil
}

func startTestServer(
	t *testing.T,
	lis net.Listener,
	registerTest bool,
	registerExtra bool,
) *grpc.Server {
	testServer := grpc.NewServer()
	if registerTest {
		pb.RegisterTestServiceServer(testServer, &testService{})
	}
	if registerExtra {
		pb.RegisterExtraServiceServer(testServer, &testService{})
	}

	reflection.Register(testServer)

	go func() {
		err := testServer.Serve(lis)
		assert.NilError(t, err, "test server error")
	}()
	return testServer
}

func assertUnaryRequest(t *testing.T, req *pb.Request, rep *pb.Reply) {
	expected := replyFromRequest(req)
	opts := cmpopts.IgnoreUnexported(pb.Reply{}, pb.InnerMessage{})
	assert.DeepEqual(t, expected, rep, opts)
}

func replyFromRequest(request *pb.Request) *pb.Reply {
	doubleField := float64(len(request.StringField))
	for _, val := range request.RepeatedField {
		doubleField += float64(val)
	}

	innerMsg := &pb.InnerMessage{RepeatedString: []string{}}
	for _, inner := range request.RepeatedInnerMsg {
		repeatedString := ""
		for _, str := range inner.RepeatedString {
			repeatedString += str
		}
		innerMsg.RepeatedString = append(
			innerMsg.RepeatedString,
			repeatedString,
		)
	}
	return &pb.Reply{
		DoubleField: doubleField,
		InnerMsg:    innerMsg,
	}
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

	go st.Run(
		&RunCfg{
			pubSub: &events.MockPubSub{},
			term:   term,
			done:   done,
			errs:   errs,
		},
	)

	for i := 1; i < 10; i++ {
		state = <-ch
		assert.NilError(t, err, "next at iter %d", i)
		assert.Equal(t, Id(i), state.Id(), "id at iter %d", i)
		opt := cmpopts.IgnoreUnexported(dynamic.Message{})
		assert.DeepEqual(
			t,
			msg.NewEmpty().GrpcMsg(),
			state.Msg().GrpcMsg(),
			opt,
		)
	}

	close(term)
	<-done
	close(errs)
	assert.Assert(t, len(errs) == 0)

	st.Close()
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
		NewState(3, testMergeMessage(t, 3)),
		NewState(6, testMergeMessage(t, 6)),
	}

	go func() {
		input1 <- NewState(1, testMergeInner1Message(t, 1))
		input1 <- NewState(2, testMergeInner1Message(t, 2))
		input1 <- NewState(3, testMergeInner1Message(t, 3))
		input1 <- NewState(6, testMergeInner1Message(t, 6))
	}()

	go func() {
		input2 <- NewState(2, testMergeInner2Message(t, 2))
		input2 <- NewState(3, testMergeInner2Message(t, 3))
		input2 <- NewState(5, testMergeInner2Message(t, 5))
		input2 <- NewState(6, testMergeInner2Message(t, 6))
	}()

	go func() {
		input3 <- NewState(1, testMergeInner3Message(t, 2))
		input3 <- NewState(3, testMergeInner3Message(t, 3))
		input3 <- NewState(5, testMergeInner3Message(t, 5))
		input3 <- NewState(6, testMergeInner3Message(t, 6))
	}()

	go s.Run(
		&RunCfg{
			pubSub: &events.MockPubSub{},
			term:   term,
			done:   done,
			errs:   errs,
		},
	)

	for i, exp := range expected {
		out := <-output
		assert.NilError(t, out.Err(), "out err at iter %d", i)
		assert.Equal(t, exp.id, out.Id(), "id at iter %d", i)
		expDyn, ok := exp.Msg().GrpcMsg().(*dynamic.Message)
		expMsg := &pb.MergeMessage{}
		err = expDyn.ConvertTo(expMsg)
		assert.NilError(t, err, "convert dyn to exp")
		assert.Assert(t, ok, "cast for exp at iter %d", i)
		dynMsg, ok := out.Msg().GrpcMsg().(*dynamic.Message)
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

	close(input1)
	close(input2)
	close(input3)
	s.Close()
}

func testMergeMessage(t *testing.T, val int32) rpc.DynMessage {
	protoMsg := &pb.MergeMessage{
		In1: &pb.MergeInner1{Val: val},
		In2: &pb.MergeInner2{Val: val},
		In3: &pb.MergeInner3{Val: val},
	}
	msg, err := rpc.DynMessageFromProto(protoMsg)
	assert.NilError(t, err, "create merge message")
	return msg
}

func testMergeInner1Message(t *testing.T, val int32) rpc.DynMessage {
	protoMsg := &pb.MergeInner1{Val: val}
	msg, err := rpc.DynMessageFromProto(protoMsg)
	assert.NilError(t, err, "create merge inner 1message")
	return msg
}

func testMergeInner2Message(t *testing.T, val int32) rpc.DynMessage {
	protoMsg := &pb.MergeInner2{Val: val}
	msg, err := rpc.DynMessageFromProto(protoMsg)
	assert.NilError(t, err, "create merge inner 2 message")
	return msg
}

func testMergeInner3Message(t *testing.T, val int32) rpc.DynMessage {
	protoMsg := &pb.MergeInner3{Val: val}
	msg, err := rpc.DynMessageFromProto(protoMsg)
	assert.NilError(t, err, "create merge inner 3 message")
	return msg
}

func TestSplitStage_Run(t *testing.T) {
	var err error
	fields := []string{"out1", "", "out2"}

	input := make(chan *State)

	output1 := make(chan *State)
	output2 := make(chan *State)
	output3 := make(chan *State)

	outputs := []chan<- *State{output1, output2, output3}

	s := NewSplitStage(fields, input, outputs)

	term := make(chan struct{})
	done := make(chan struct{})
	errs := make(chan error)

	expected1 := []*State{
		NewState(Id(1), testSplitInner1Message(t, 1)),
		NewState(Id(3), testSplitInner1Message(t, 3)),
		NewState(Id(5), testSplitInner1Message(t, 5)),
	}
	expected2 := []*State{
		NewState(Id(1), testSplitMessage(t, 1)),
		NewState(Id(3), testSplitMessage(t, 3)),
		NewState(Id(5), testSplitMessage(t, 5)),
	}
	expected3 := []*State{
		NewState(Id(1), testSplitInner2Message(t, 1)),
		NewState(Id(3), testSplitInner2Message(t, 3)),
		NewState(Id(5), testSplitInner2Message(t, 5)),
	}

	go func() {
		input <- NewState(Id(1), testSplitMessage(t, 1))
		input <- NewState(Id(3), testSplitMessage(t, 3))
		input <- NewState(Id(5), testSplitMessage(t, 5))
	}()

	go s.Run(
		&RunCfg{
			pubSub: &events.MockPubSub{},
			term:   term,
			done:   done,
			errs:   errs,
		},
	)

	for i := 0; i < len(expected1); i++ {
		exp1 := expected1[i]
		out1 := <-output1
		assert.NilError(t, out1.Err(), "err 1 at iter %d", i)
		assert.Equal(t, exp1.Id(), out1.Id(), "id 1 at iter %d", i)
		expDyn1, ok := exp1.Msg().GrpcMsg().(*dynamic.Message)
		assert.Assert(t, ok, "cast for exp 1 at iter %d", i)
		expMsg1 := &pb.SplitInner1{}
		err = expDyn1.ConvertTo(expMsg1)
		assert.NilError(t, err, "convert dyn 1 to exp 1")
		outDyn1, ok := out1.Msg().GrpcMsg().(*dynamic.Message)
		assert.Assert(t, ok, "cast for out 1 at iter %d", i)
		outMsg1 := &pb.SplitInner1{}
		err = outDyn1.ConvertTo(outMsg1)
		assert.Equal(t, expMsg1.Val, outMsg1.Val)

		exp2 := expected2[i]
		out2 := <-output2
		assert.NilError(t, out2.Err(), "err 2 at iter %d", i)
		assert.Equal(t, exp2.Id(), out2.Id(), "id 2 at iter %d", i)
		expDyn2, ok := exp2.Msg().GrpcMsg().(*dynamic.Message)
		assert.Assert(t, ok, "cast for exp 2 at iter %d", i)
		expMsg2 := &pb.SplitMessage{}
		err = expDyn2.ConvertTo(expMsg2)
		assert.NilError(t, err, "convert dyn 2 to exp 2")
		dynDyn2, ok := out2.Msg().GrpcMsg().(*dynamic.Message)
		assert.Assert(t, ok, "cast for out 2 at iter %d", i)
		outMsg2 := &pb.SplitMessage{}
		err = dynDyn2.ConvertTo(outMsg2)
		assert.NilError(t, err, "convert dyn 2 to out 2")
		assert.Equal(t, expMsg2.Out1.Val, outMsg2.Out1.Val)
		assert.Equal(t, expMsg2.Val, outMsg2.Val)
		assert.Equal(t, expMsg2.Out2.Val, outMsg2.Out2.Val)

		exp3 := expected3[i]
		out3 := <-output3
		assert.NilError(t, out3.Err(), "err 3 at iter %d", i)
		assert.Equal(t, exp3.Id(), out3.Id(), "id 3 at iter %d", i)
		expDyn3, ok := exp3.Msg().GrpcMsg().(*dynamic.Message)
		assert.Assert(t, ok, "cast for exp 3 at iter %d", i)
		expMsg3 := &pb.SplitInner2{}
		err = expDyn3.ConvertTo(expMsg3)
		assert.NilError(t, err, "convert dyn 3 to exp 3")
		outDyn3, ok := out3.Msg().GrpcMsg().(*dynamic.Message)
		assert.Assert(t, ok, "cast for out 3 at iter %d", i)
		outMsg3 := &pb.SplitInner2{}
		err = outDyn3.ConvertTo(outMsg3)
		assert.Equal(t, expMsg3.Val, outMsg3.Val)
	}

	close(term)
	<-done
	close(errs)
	assert.Assert(t, len(errs) == 0)
	close(input)
	s.Close()
}

func testSplitMessage(t *testing.T, val int32) rpc.DynMessage {
	protoMsg := &pb.SplitMessage{
		Out1: &pb.SplitInner1{Val: val},
		Val:  val,
		Out2: &pb.SplitInner2{Val: val},
	}
	msg, err := rpc.DynMessageFromProto(protoMsg)
	assert.NilError(t, err, "create split 1message")
	return msg
}

func testSplitInner1Message(t *testing.T, val int32) rpc.DynMessage {
	protoMsg := &pb.SplitInner1{Val: val}
	msg, err := rpc.DynMessageFromProto(protoMsg)
	assert.NilError(t, err, "create split inner 1message")
	return msg
}

func testSplitInner2Message(t *testing.T, val int32) rpc.DynMessage {
	protoMsg := &pb.SplitInner2{Val: val}
	msg, err := rpc.DynMessageFromProto(protoMsg)
	assert.NilError(t, err, "create split inner 2 message")
	return msg
}
