package invoke

import (
	"context"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/test/protobuf/unit"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/grpc"
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

	assert.Assert(t, errdefs.IsUnknown(err), "error is not unknown")
	assert.ErrorContains(t, err, "unary invoke: ")
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

	invoke := NewUnaryInvoke("unit.ExtraService/ExtraMethod", conn)

	reply := &unit.Reply{}
	err = invoke(ctx, errorRequest, reply)

	assert.Assert(
		t,
		errdefs.IsFailedPrecondition(err),
		"error is not FailedPrecondition",
	)
	assert.ErrorContains(t, err, "unary invoke: ")
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

	invoke := NewUnaryInvoke("unit.TestService/NoSuchMethod", conn)

	reply := &unit.Reply{}
	err = invoke(ctx, errorRequest, reply)

	assert.Assert(
		t,
		errdefs.IsFailedPrecondition(err),
		"error is not FailedPrecondition",
	)
	assert.ErrorContains(t, err, "unary invoke: ")
}
