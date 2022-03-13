package invoke

import (
	"context"
	"errors"
	"github.com/DuarteMRAlves/maestro/test/protobuf/unit"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gotest.tools/v3/assert"
	"net"
	"testing"
	"time"
)

var (
	correctRequest = &unit.Request{
		StringField:   "some-string",
		RepeatedField: []int64{1, 2, 3, 4},
		RepeatedInnerMsg: []*unit.InnerMessage{
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

	errorRequest = &unit.Request{StringField: "error"}

	expectedReply = &unit.Reply{
		// Value equal to len(StringField) + sum(RepeatedField)
		DoubleField: 21,
		InnerMsg: &unit.InnerMessage{
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
	testServer := startServer(t, lis, false)
	defer testServer.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	assert.NilError(t, err, "dial error")
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		assert.NilError(t, err, "close connection")
	}(conn)

	invoke := NewUnaryInvoke("unit.TestService/Unary", conn)

	reply := &unit.Reply{}
	err = invoke(ctx, correctRequest, reply)
	assert.NilError(t, err, "invoke error not nil")

	assert.DeepEqual(
		t,
		expectedReply,
		reply,
		cmpopts.IgnoreUnexported(unit.Reply{}, unit.InnerMessage{}),
	)
}

func TestUnaryClient_Invoke_ErrorReturned(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:0")
	assert.NilError(t, err, "failed to listen")
	addr := lis.Addr().String()
	testServer := startServer(t, lis, false)
	defer testServer.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	assert.NilError(t, err, "dial error")
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		assert.NilError(t, err, "close connection")
	}(conn)

	invoke := NewUnaryInvoke("unit.TestService/Unary", conn)

	reply := &unit.Reply{}
	err = invoke(ctx, errorRequest, reply)

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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	assert.NilError(t, err, "dial error")
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		assert.NilError(t, err, "close connection")
	}(conn)

	method := "unit.ExtraService/ExtraMethod"
	invoke := NewUnaryInvoke(method, conn)

	reply := &unit.Reply{}
	err = invoke(ctx, errorRequest, reply)
	assert.Assert(t, err != nil)
	cause, ok := errors.Unwrap(err).(interface {
		GRPCStatus() *status.Status
	})
	st := cause.GRPCStatus()
	assert.Assert(t, ok, "correct type")
	assert.Equal(t, codes.Unimplemented, st.Code())
}

func TestUnaryClient_Invoke_MethodDoesNotExist(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:0")
	assert.NilError(t, err, "failed to listen")
	addr := lis.Addr().String()
	testServer := startServer(t, lis, false)
	defer testServer.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	assert.NilError(t, err, "dial error")
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		assert.NilError(t, err, "close connection")
	}(conn)

	method := "unit.TestService/NoSuchMethod"
	invoke := NewUnaryInvoke(method, conn)

	reply := &unit.Reply{}
	err = invoke(ctx, errorRequest, reply)
	assert.Assert(t, err != nil)
	cause, ok := errors.Unwrap(err).(interface {
		GRPCStatus() *status.Status
	})
	st := cause.GRPCStatus()
	assert.Assert(t, ok, "correct type")
	assert.Equal(t, codes.Unimplemented, st.Code())
}
