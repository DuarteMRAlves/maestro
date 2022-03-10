package execute

import (
	"context"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/invoke"
	"github.com/DuarteMRAlves/maestro/test/protobuf/unit"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jhump/protoreflect/dynamic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gotest.tools/v3/assert"
	"net"
	"testing"
)

func TestUnaryStage_Run(t *testing.T) {
	var collect []state

	stageDone := make(chan struct{})
	serverDone := make(chan struct{})

	requests := testRequests(t)
	states := []state{
		newState(1, requests[0]),
		newState(3, requests[1]),
		newState(5, requests[2]),
	}

	lis, err := net.Listen("tcp", "localhost:0")
	assert.NilError(t, err, "failed to listen")
	addr := lis.Addr().String()
	server := startTestServer(t, lis)
	defer server.Stop()

	input := make(chan state, len(requests))
	output := make(chan state, len(requests))

	outDesc, err := invoke.NewMessageDescriptor(&unit.Reply{})
	assert.NilError(t, err, "create input descriptor")

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	assert.NilError(t, err, "unable to connect to address: %s", addr)

	invokeFn := invoke.NewUnaryInvoke("unit.TestService/Unary", conn)

	stage := newUnaryStage(input, output, outDesc.MessageGenerator(), invokeFn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = stage.Run(ctx)
		close(stageDone)
	}()

	go func() {
		for i := 0; i < len(states); i++ {
			c := <-output
			collect = append(collect, c)
		}
		close(serverDone)
	}()

	input <- states[0]
	input <- states[1]
	input <- states[2]
	<-serverDone
	cancel()
	<-stageDone
	close(input)

	assert.Assert(t, len(collect) == len(states))
	for i, c := range collect {
		in := states[i]
		assert.Equal(t, in.id, c.id, "correct received id")

		dynReq, err := dynamic.AsDynamicMessage(in.msg.GrpcMessage())
		assert.Assert(t, err, "request dynamic message")
		req := &unit.Request{}
		err = dynReq.ConvertTo(req)
		assert.NilError(t, err, "convert dynamic to Request")

		dynRep, err := dynamic.AsDynamicMessage(c.msg.GrpcMessage())
		assert.NilError(t, err, "reply dynamic message")
		rep := &unit.Reply{}
		err = dynRep.ConvertTo(rep)
		assert.NilError(t, err, "convert dynamic to Reply")

		assertUnaryRequest(t, req, rep)
	}
}

func testRequests(t *testing.T) []invoke.DynamicMessage {
	msg1, err := invoke.NewDynamicMessage(
		&unit.Request{
			StringField:   "string-1",
			RepeatedField: []int64{1, 2, 3, 4},
			RepeatedInnerMsg: []*unit.InnerMessage{
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
	msg2, err := invoke.NewDynamicMessage(
		&unit.Request{
			StringField:   "string-2",
			RepeatedField: []int64{1, 2, 3, 4},
			RepeatedInnerMsg: []*unit.InnerMessage{
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
	msg3, err := invoke.NewDynamicMessage(
		&unit.Request{
			StringField:   "string-3",
			RepeatedField: []int64{1, 2, 3, 4},
			RepeatedInnerMsg: []*unit.InnerMessage{
				{
					RepeatedString: []string{"hello", "world", "3"},
				},
				{
					RepeatedString: []string{"other", "message", "3"},
				},
			},
		},
	)
	assert.NilError(t, err, "create message 3")

	return []invoke.DynamicMessage{msg1, msg2, msg3}
}

func assertUnaryRequest(t *testing.T, req *unit.Request, rep *unit.Reply) {
	expected := replyFromRequest(req)
	opts := cmpopts.IgnoreUnexported(unit.Reply{}, unit.InnerMessage{})
	assert.DeepEqual(t, expected, rep, opts)
}

func startTestServer(
	t *testing.T,
	lis net.Listener,
) *grpc.Server {
	testServer := grpc.NewServer()
	service := &testService{}
	unit.RegisterTestServiceServer(testServer, service)

	reflection.Register(testServer)

	go func() {
		err := testServer.Serve(lis)
		assert.NilError(t, err, "test server error")
	}()
	return testServer
}

type testService struct {
	unit.UnimplementedTestServiceServer
	unit.UnimplementedExtraServiceServer
}

func (s *testService) Unary(
	_ context.Context,
	request *unit.Request,
) (*unit.Reply, error) {
	if request.StringField == "error" {
		return nil, fmt.Errorf("dummy error")
	}
	return replyFromRequest(request), nil
}

func replyFromRequest(request *unit.Request) *unit.Reply {
	doubleField := float64(len(request.StringField))
	for _, val := range request.RepeatedField {
		doubleField += float64(val)
	}

	innerMsg := &unit.InnerMessage{RepeatedString: []string{}}
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
	return &unit.Reply{
		DoubleField: doubleField,
		InnerMsg:    innerMsg,
	}
}

func TestSourceStage_Run(t *testing.T) {
	start := int32(1)
	numRequest := 10

	msgDesc, err := invoke.NewMessageDescriptor(&unit.Request{})
	assert.NilError(t, err, "request message descriptor")

	output := make(chan state)
	s := newSourceStage(start, msgDesc.MessageGenerator(), output)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})

	go func() {
		_ = s.Run(ctx)
		close(done)
	}()

	generated := make([]state, 0, numRequest)
	for i := 0; i < numRequest; i++ {
		generated = append(generated, <-output)
	}
	cancel()
	<-done

	for i, g := range generated {
		assert.Equal(t, int32(g.id), int32(i+1))
	}
}

func TestMergeStage_Run(t *testing.T) {
	f1, err := domain.NewMessageField("in1")
	assert.NilError(t, err, "create field 1")

	f2, err := domain.NewMessageField("in2")
	assert.NilError(t, err, "create field 2")

	f3, err := domain.NewMessageField("in3")
	assert.NilError(t, err, "create field 3")

	fields := []domain.MessageField{f1, f2, f3}

	input1 := make(chan state)
	defer close(input1)
	input2 := make(chan state)
	defer close(input2)
	input3 := make(chan state)
	defer close(input3)
	inputs := []<-chan state{input1, input2, input3}

	output := make(chan state)

	outDesc, err := invoke.NewMessageDescriptor(&unit.MergeMessage{})
	assert.NilError(t, err, "create merge message descriptor")

	s := newMergeStage(fields, inputs, output, outDesc.MessageGenerator())

	expected := []state{
		newState(3, testMergeMessage(t, 3)),
		newState(6, testMergeMessage(t, 6)),
	}

	go func() {
		input1 <- newState(1, testMergeInner1Message(t, 1))
		input1 <- newState(2, testMergeInner1Message(t, 2))
		input1 <- newState(3, testMergeInner1Message(t, 3))
		input1 <- newState(6, testMergeInner1Message(t, 6))
	}()

	go func() {
		input2 <- newState(2, testMergeInner2Message(t, 2))
		input2 <- newState(3, testMergeInner2Message(t, 3))
		input2 <- newState(5, testMergeInner2Message(t, 5))
		input2 <- newState(6, testMergeInner2Message(t, 6))
	}()

	go func() {
		input3 <- newState(1, testMergeInner3Message(t, 2))
		input3 <- newState(3, testMergeInner3Message(t, 3))
		input3 <- newState(5, testMergeInner3Message(t, 5))
		input3 <- newState(6, testMergeInner3Message(t, 6))
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})

	go func() {
		err := s.Run(ctx)
		assert.NilError(t, err, "run error")
		close(done)
	}()

	for i, exp := range expected {
		out := <-output
		assert.Equal(t, exp.id, out.id, "id at iter %d", i)
		expDyn, ok := exp.msg.GrpcMessage().(*dynamic.Message)
		expMsg := &unit.MergeMessage{}
		err = expDyn.ConvertTo(expMsg)
		assert.NilError(t, err, "convert dyn to exp")
		assert.Assert(t, ok, "cast for exp at iter %d", i)
		dynMsg, ok := out.msg.GrpcMessage().(*dynamic.Message)
		assert.Assert(t, ok, "cast for out at iter %d", i)
		outMsg := &unit.MergeMessage{}
		err = dynMsg.ConvertTo(outMsg)
		assert.NilError(t, err, "convert dyn to out")
		assert.Equal(t, expMsg.In1.Val, outMsg.In1.Val)
		assert.Equal(t, expMsg.In2.Val, outMsg.In2.Val)
		assert.Equal(t, expMsg.In3.Val, outMsg.In3.Val)
	}
	cancel()
	<-done
}

func testMergeMessage(t *testing.T, val int32) invoke.DynamicMessage {
	protoMsg := &unit.MergeMessage{
		In1: &unit.MergeInner1{Val: val},
		In2: &unit.MergeInner2{Val: val},
		In3: &unit.MergeInner3{Val: val},
	}
	msg, err := invoke.NewDynamicMessage(protoMsg)
	assert.NilError(t, err, "create merge message")
	return msg
}

func testMergeInner1Message(t *testing.T, val int32) invoke.DynamicMessage {
	protoMsg := &unit.MergeInner1{Val: val}
	msg, err := invoke.NewDynamicMessage(protoMsg)
	assert.NilError(t, err, "create merge inner 1message")
	return msg
}

func testMergeInner2Message(t *testing.T, val int32) invoke.DynamicMessage {
	protoMsg := &unit.MergeInner2{Val: val}
	msg, err := invoke.NewDynamicMessage(protoMsg)
	assert.NilError(t, err, "create merge inner 2 message")
	return msg
}

func testMergeInner3Message(t *testing.T, val int32) invoke.DynamicMessage {
	protoMsg := &unit.MergeInner3{Val: val}
	msg, err := invoke.NewDynamicMessage(protoMsg)
	assert.NilError(t, err, "create merge inner 3 message")
	return msg
}

func TestSplitStage_Run(t *testing.T) {
	f1, err := domain.NewMessageField("out1")
	assert.NilError(t, err, "create message field 1")

	f3, err := domain.NewMessageField("out2")
	assert.NilError(t, err, "create message field 3")

	fields := []domain.OptionalMessageField{
		domain.NewPresentMessageField(f1),
		domain.NewEmptyMessageField(),
		domain.NewPresentMessageField(f3),
	}

	input := make(chan state)

	output1 := make(chan state)
	output2 := make(chan state)
	output3 := make(chan state)

	outputs := []chan<- state{output1, output2, output3}

	s := newSplitStage(fields, input, outputs)

	expected1 := []state{
		newState(id(1), testSplitInner1Message(t, 1)),
		newState(id(3), testSplitInner1Message(t, 3)),
		newState(id(5), testSplitInner1Message(t, 5)),
	}
	expected2 := []state{
		newState(id(1), testSplitMessage(t, 1)),
		newState(id(3), testSplitMessage(t, 3)),
		newState(id(5), testSplitMessage(t, 5)),
	}
	expected3 := []state{
		newState(id(1), testSplitInner2Message(t, 1)),
		newState(id(3), testSplitInner2Message(t, 3)),
		newState(id(5), testSplitInner2Message(t, 5)),
	}

	go func() {
		input <- newState(id(1), testSplitMessage(t, 1))
		input <- newState(id(3), testSplitMessage(t, 3))
		input <- newState(id(5), testSplitMessage(t, 5))
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})

	go func() {
		err := s.Run(ctx)
		assert.NilError(t, err, "run error")
		close(done)
	}()

	for i := 0; i < len(expected1); i++ {
		exp1 := expected1[i]
		out1 := <-output1
		assert.Equal(t, exp1.id, out1.id, "id 1 at iter %d", i)
		expDyn1, ok := exp1.msg.GrpcMessage().(*dynamic.Message)
		assert.Assert(t, ok, "cast for exp 1 at iter %d", i)
		expMsg1 := &unit.SplitInner1{}
		err = expDyn1.ConvertTo(expMsg1)
		assert.NilError(t, err, "convert dyn 1 to exp 1")
		outDyn1, ok := out1.msg.GrpcMessage().(*dynamic.Message)
		assert.Assert(t, ok, "cast for out 1 at iter %d", i)
		outMsg1 := &unit.SplitInner1{}
		err = outDyn1.ConvertTo(outMsg1)
		assert.Equal(t, expMsg1.Val, outMsg1.Val)

		exp2 := expected2[i]
		out2 := <-output2
		assert.Equal(t, exp2.id, out2.id, "id 2 at iter %d", i)
		expDyn2, ok := exp2.msg.GrpcMessage().(*dynamic.Message)
		assert.Assert(t, ok, "cast for exp 2 at iter %d", i)
		expMsg2 := &unit.SplitMessage{}
		err = expDyn2.ConvertTo(expMsg2)
		assert.NilError(t, err, "convert dyn 2 to exp 2")
		dynDyn2, ok := out2.msg.GrpcMessage().(*dynamic.Message)
		assert.Assert(t, ok, "cast for out 2 at iter %d", i)
		outMsg2 := &unit.SplitMessage{}
		err = dynDyn2.ConvertTo(outMsg2)
		assert.NilError(t, err, "convert dyn 2 to out 2")
		assert.Equal(t, expMsg2.Out1.Val, outMsg2.Out1.Val)
		assert.Equal(t, expMsg2.Val, outMsg2.Val)
		assert.Equal(t, expMsg2.Out2.Val, outMsg2.Out2.Val)

		exp3 := expected3[i]
		out3 := <-output3
		assert.Equal(t, exp3.id, out3.id, "id 3 at iter %d", i)
		expDyn3, ok := exp3.msg.GrpcMessage().(*dynamic.Message)
		assert.Assert(t, ok, "cast for exp 3 at iter %d", i)
		expMsg3 := &unit.SplitInner2{}
		err = expDyn3.ConvertTo(expMsg3)
		assert.NilError(t, err, "convert dyn 3 to exp 3")
		outDyn3, ok := out3.msg.GrpcMessage().(*dynamic.Message)
		assert.Assert(t, ok, "cast for out 3 at iter %d", i)
		outMsg3 := &unit.SplitInner2{}
		err = outDyn3.ConvertTo(outMsg3)
		assert.Equal(t, expMsg3.Val, outMsg3.Val)
	}
	cancel()
	<-done
}

func testSplitMessage(t *testing.T, val int32) invoke.DynamicMessage {
	protoMsg := &unit.SplitMessage{
		Out1: &unit.SplitInner1{Val: val},
		Val:  val,
		Out2: &unit.SplitInner2{Val: val},
	}
	msg, err := invoke.NewDynamicMessage(protoMsg)
	assert.NilError(t, err, "create split 1message")
	return msg
}

func testSplitInner1Message(t *testing.T, val int32) invoke.DynamicMessage {
	protoMsg := &unit.SplitInner1{Val: val}
	msg, err := invoke.NewDynamicMessage(protoMsg)
	assert.NilError(t, err, "create split inner 1message")
	return msg
}

func testSplitInner2Message(t *testing.T, val int32) invoke.DynamicMessage {
	protoMsg := &unit.SplitInner2{Val: val}
	msg, err := invoke.NewDynamicMessage(protoMsg)
	assert.NilError(t, err, "create split inner 2 message")
	return msg
}
