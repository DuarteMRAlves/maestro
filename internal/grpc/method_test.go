package grpc

import (
	"context"
	"errors"
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/test/protobuf/unit"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"gotest.tools/v3/assert"
	"net"
	"testing"
	"time"
)

var (
	correctRequest = &unit.TestMethodRequest{
		StringField:   "some-string",
		RepeatedField: []int64{1, 2, 3, 4},
		RepeatedInnerMsg: []*unit.TestMethodInnerMessage{
			{
				RepeatedString: []string{
					"hello",
					"world",
				},
			},
			{
				RepeatedString: []string{
					"other",
					"message",
				},
			},
		},
	}

	errorRequest = &unit.TestMethodRequest{StringField: "error"}

	expectedReply = &unit.TestMethodReply{
		// Value equal to len(StringField) + sum(RepeatedField)
		DoubleField: 21,
		InnerMsg: &unit.TestMethodInnerMessage{
			RepeatedString: []string{
				"helloworld",
				"othermessage",
			},
		},
	}
)

func TestUnaryClient_Invoke(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:0")
	assert.NilError(t, err, "failed to listen")
	addr := lis.Addr().String()
	testServer := testMethodStartServer(t, lis)
	defer testServer.Stop()

	inDesc, err := newMessageDescriptor(&unit.TestMethodRequest{})
	assert.NilError(t, err, "create input message descriptor")
	outDesc, err := newMessageDescriptor(&unit.TestMethodReply{})
	assert.NilError(t, err, "create output message descriptor")

	methodName := "unit.TestMethodService/CorrectMethod"
	method := newUnaryMethod(methodName, inDesc, outDesc)

	clientBuilder := method.ClientBuilder()
	client, err := clientBuilder(internal.NewAddress(addr))
	assert.NilError(t, err, "build client")
	defer func() {
		assert.NilError(t, client.Close())
	}()

	req, err := newMessage(correctRequest)
	assert.NilError(t, err, "create request")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	reply, err := client.Call(ctx, req)
	assert.NilError(t, err, "invoke error not nil")

	grpcMsg, ok := reply.(*message)
	assert.Assert(t, ok, "cast reply to grpc message")

	pbReply := &unit.TestMethodReply{}
	err = grpcMsg.dynMsg.ConvertTo(pbReply)
	assert.NilError(t, err, "convert dynamic replay to message")

	cmpOpts := cmpopts.IgnoreUnexported(
		unit.TestMethodReply{},
		unit.TestMethodInnerMessage{},
	)
	assert.DeepEqual(t, expectedReply, pbReply, cmpOpts)
}

func TestUnaryClient_Invoke_ErrorReturned(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:0")
	assert.NilError(t, err, "failed to listen")
	addr := lis.Addr().String()
	testServer := testMethodStartServer(t, lis)
	defer testServer.Stop()

	inDesc, err := newMessageDescriptor(&unit.TestMethodRequest{})
	assert.NilError(t, err, "create input message descriptor")
	outDesc, err := newMessageDescriptor(&unit.TestMethodReply{})
	assert.NilError(t, err, "create output message descriptor")

	methodName := "unit.TestMethodService/CorrectMethod"
	method := newUnaryMethod(methodName, inDesc, outDesc)

	clientBuilder := method.ClientBuilder()
	client, err := clientBuilder(internal.NewAddress(addr))
	assert.NilError(t, err, "build client")
	defer func() {
		assert.NilError(t, client.Close())
	}()

	req, err := newMessage(errorRequest)
	assert.NilError(t, err, "create request")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	reply, err := client.Call(ctx, req)
	assert.Assert(t, reply == nil)
	assert.Assert(t, err != nil)
	cause, ok := errors.Unwrap(err).(interface {
		GRPCStatus() *status.Status
	})
	st := cause.GRPCStatus()
	assert.Assert(t, ok, "correct type")
	assert.Equal(t, codes.Unknown, st.Code())
}

func TestUnaryClient_Invoke_MethodUnimplemented(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:0")
	assert.NilError(t, err, "failed to listen")
	addr := lis.Addr().String()
	testServer := startServer(t, lis, false)
	defer testServer.Stop()

	inDesc, err := newMessageDescriptor(&unit.TestMethodRequest{})
	assert.NilError(t, err, "create input message descriptor")
	outDesc, err := newMessageDescriptor(&unit.TestMethodReply{})
	assert.NilError(t, err, "create output message descriptor")

	methodName := "unit.TestMethodService/UnimplementedMethod"
	method := newUnaryMethod(methodName, inDesc, outDesc)

	clientBuilder := method.ClientBuilder()
	client, err := clientBuilder(internal.NewAddress(addr))
	assert.NilError(t, err, "build client")
	defer func() {
		assert.NilError(t, client.Close())
	}()

	req, err := newMessage(errorRequest)
	assert.NilError(t, err, "create request")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	reply, err := client.Call(ctx, req)
	assert.Assert(t, reply == nil)
	assert.Assert(t, err != nil)
	cause, ok := errors.Unwrap(err).(interface {
		GRPCStatus() *status.Status
	})
	st := cause.GRPCStatus()
	assert.Assert(t, ok, "correct type")
	assert.Equal(t, codes.Unimplemented, st.Code())
}

var dummyErr = errors.New("dummy error")

type testMethodService struct {
	unit.UnimplementedTestMethodServiceServer
}

func (s *testMethodService) CorrectMethod(
	_ context.Context,
	request *unit.TestMethodRequest,
) (*unit.TestMethodReply, error) {
	if request.StringField == "error" {
		return nil, dummyErr
	} else {
		return testReplyFromRequest(request), nil
	}
}

func testReplyFromRequest(req *unit.TestMethodRequest) *unit.TestMethodReply {
	doubleField := float64(len(req.StringField))
	for _, val := range req.RepeatedField {
		doubleField += float64(val)
	}

	innerMsg := &unit.TestMethodInnerMessage{RepeatedString: []string{}}
	for _, inner := range req.RepeatedInnerMsg {
		repeatedString := ""
		for _, str := range inner.RepeatedString {
			repeatedString += str
		}
		innerMsg.RepeatedString = append(
			innerMsg.RepeatedString,
			repeatedString,
		)
	}
	return &unit.TestMethodReply{
		DoubleField: doubleField,
		InnerMsg:    innerMsg,
	}
}

func testMethodStartServer(t *testing.T, lis net.Listener) *grpc.Server {
	testServer := grpc.NewServer()
	unit.RegisterTestMethodServiceServer(testServer, &testMethodService{})
	reflection.Register(testServer)
	go func() {
		err := testServer.Serve(lis)
		assert.NilError(t, err, "test server error")
	}()
	return testServer
}
