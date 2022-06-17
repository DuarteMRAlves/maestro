package grpc

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/DuarteMRAlves/maestro/test/protobuf/unit"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
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
	if err != nil {
		t.Fatalf("failed to listen: %s", err)
	}
	addr := lis.Addr().String()
	testServer := testMethodStartServer(t, lis)
	defer testServer.Stop()

	inDesc, err := newMessageDescriptor(&unit.TestMethodRequest{})
	if err != nil {
		t.Fatalf("create input message descriptor: %s", err)
	}
	outDesc, err := newMessageDescriptor(&unit.TestMethodReply{})
	if err != nil {
		t.Fatalf("create output message descriptor: %s", err)
	}

	methodName := "unit.TestMethodService/CorrectMethod"
	method := newUnaryMethod(addr, methodName, inDesc, outDesc)

	conn, err := method.Dial()
	if err != nil {
		t.Fatalf("build conn: %s", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			t.Fatalf("close conn: %s", err)
		}
	}()

	req, err := newMessage(correctRequest)
	if err != nil {
		t.Fatalf("create request: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	reply, err := conn.Call(ctx, req)
	if err != nil {
		t.Fatalf("call method: %s", err)
	}

	grpcMsg, ok := reply.(*message)
	if !ok {
		t.Fatalf("cast reply to grpcMsg")
	}

	pbReply := &unit.TestMethodReply{}
	err = grpcMsg.dynMsg.ConvertTo(pbReply)
	if err != nil {
		t.Fatalf("convert grpcMsg to pbReply: %s", err)
	}

	cmpOpts := cmpopts.IgnoreUnexported(
		unit.TestMethodReply{},
		unit.TestMethodInnerMessage{},
	)
	if diff := cmp.Diff(expectedReply, pbReply, cmpOpts); diff != "" {
		t.Fatalf("reply mismatch:\n%s", diff)
	}
}

func TestUnaryClient_Invoke_ErrorReturned(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen: %s", err)
	}
	addr := lis.Addr().String()
	testServer := testMethodStartServer(t, lis)
	defer testServer.Stop()

	inDesc, err := newMessageDescriptor(&unit.TestMethodRequest{})
	if err != nil {
		t.Fatalf("create input message descriptor: %s", err)
	}
	outDesc, err := newMessageDescriptor(&unit.TestMethodReply{})
	if err != nil {
		t.Fatalf("create input output descriptor: %s", err)
	}

	methodName := "unit.TestMethodService/CorrectMethod"
	method := newUnaryMethod(addr, methodName, inDesc, outDesc)

	conn, err := method.Dial()
	if err != nil {
		t.Fatalf("build conn: %s", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			t.Fatalf("close conn: %s", err)
		}
	}()

	req, err := newMessage(errorRequest)
	if err != nil {
		t.Fatalf("create request: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	reply, err := conn.Call(ctx, req)
	if reply != nil {
		t.Fatalf("replay is not nil")
	}
	if err == nil {
		t.Fatalf("expected error but received nil")
	}
	cause, ok := errors.Unwrap(err).(interface {
		GRPCStatus() *status.Status
	})
	if !ok {
		t.Fatalf("error does not implement grpc interface")
	}
	st := cause.GRPCStatus()
	if diff := cmp.Diff(codes.Unknown, st.Code()); diff != "" {
		t.Fatalf("code mismatch:\n%s", diff)
	}
}

func TestUnaryClient_Invoke_MethodUnimplemented(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen: %s", err)
	}
	addr := lis.Addr().String()
	testServer := startServer(t, lis, false)
	defer testServer.Stop()

	inDesc, err := newMessageDescriptor(&unit.TestMethodRequest{})
	if err != nil {
		t.Fatalf("create input input descriptor: %s", err)
	}
	outDesc, err := newMessageDescriptor(&unit.TestMethodReply{})
	if err != nil {
		t.Fatalf("create input output descriptor: %s", err)
	}

	methodName := "unit.TestMethodService/UnimplementedMethod"
	method := newUnaryMethod(addr, methodName, inDesc, outDesc)

	conn, err := method.Dial()
	if err != nil {
		t.Fatalf("build conn: %s", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			t.Fatalf("close conn: %s", err)
		}
	}()

	req, err := newMessage(errorRequest)
	if err != nil {
		t.Fatalf("create request: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	reply, err := conn.Call(ctx, req)
	if reply != nil {
		t.Fatalf("replay is not nil")
	}
	if err == nil {
		t.Fatalf("expected error but received nil")
	}
	cause, ok := errors.Unwrap(err).(interface {
		GRPCStatus() *status.Status
	})
	if !ok {
		t.Fatalf("error does not implement grpc interface")
	}
	st := cause.GRPCStatus()
	if diff := cmp.Diff(codes.Unimplemented, st.Code()); diff != "" {
		t.Fatalf("code mismatch:\n%s", diff)
	}
}

var errDummy = errors.New("dummy error")

type testMethodService struct {
	unit.UnimplementedTestMethodServiceServer
}

func (s *testMethodService) CorrectMethod(
	_ context.Context,
	request *unit.TestMethodRequest,
) (*unit.TestMethodReply, error) {
	if request.StringField == "error" {
		return nil, errDummy
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
		if err := testServer.Serve(lis); err != nil {
			t.Errorf("test server: %s", err)
		}
	}()
	return testServer
}
