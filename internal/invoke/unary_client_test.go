package invoke

import (
	"context"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/util"
	"github.com/DuarteMRAlves/maestro/tests/pb"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/grpc"
	"gotest.tools/v3/assert"
	"testing"
	"time"
)

var (
	correctRequest = &pb.Request{
		StringField:   "some-string",
		RepeatedField: []int64{1, 2, 3, 4},
		RepeatedInnerMsg: []*pb.InnerMessage{
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

	errorRequest = &pb.Request{StringField: "error"}

	expectedReply = &pb.Reply{
		// Value equal to len(StringField) + sum(RepeatedField)
		DoubleField: 21,
		InnerMsg: &pb.InnerMessage{
			RepeatedString: []string{
				"helloworld",
				"othermessage",
			},
		},
	}
)

func TestUnaryClient_Invoke(t *testing.T) {
	lis := util.NewTestListener(t)
	addr := lis.Addr().String()
	testServer := startServer(t, lis)
	defer testServer.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	assert.NilError(t, err, "dial error")
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		assert.NilError(t, err, "close connection")
	}(conn)

	client := NewUnary("pb.TestService/Unary", conn)

	reply := &pb.Reply{}
	err = client.Invoke(ctx, correctRequest, reply)
	assert.NilError(t, err, "invoke error not nil")

	assert.DeepEqual(
		t,
		expectedReply,
		reply,
		cmpopts.IgnoreUnexported(pb.Reply{}, pb.InnerMessage{}),
	)
}

func TestUnaryClient_Invoke_ErrorReturned(t *testing.T) {
	lis := util.NewTestListener(t)
	addr := lis.Addr().String()
	testServer := startServer(t, lis)
	defer testServer.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	assert.NilError(t, err, "dial error")
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		assert.NilError(t, err, "close connection")
	}(conn)

	client := NewUnary("pb.TestService/Unary", conn)

	reply := &pb.Reply{}
	err = client.Invoke(ctx, errorRequest, reply)

	assert.Assert(t, errdefs.IsUnknown(err), "error is not unknown")
	assert.ErrorContains(t, err, "unary invoke: ")
}

func TestUnaryClient_Invoke_MethodUnimplemented(t *testing.T) {
	lis := util.NewTestListener(t)
	addr := lis.Addr().String()
	testServer := startServer(t, lis)
	defer testServer.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	assert.NilError(t, err, "dial error")
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		assert.NilError(t, err, "close connection")
	}(conn)

	client := NewUnary("pb.ExtraService/ExtraMethod", conn)

	reply := &pb.Reply{}
	err = client.Invoke(ctx, errorRequest, reply)

	assert.Assert(
		t,
		errdefs.IsFailedPrecondition(err),
		"error is not FailedPrecondition",
	)
	assert.ErrorContains(t, err, "unary invoke: ")
}

func TestUnaryClient_Invoke_MethodDoesNotExist(t *testing.T) {
	lis := util.NewTestListener(t)
	addr := lis.Addr().String()
	testServer := startServer(t, lis)
	defer testServer.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	assert.NilError(t, err, "dial error")
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		assert.NilError(t, err, "close connection")
	}(conn)

	client := NewUnary("pb.TestService/NoSuchMethod", conn)

	reply := &pb.Reply{}
	err = client.Invoke(ctx, errorRequest, reply)

	assert.Assert(
		t,
		errdefs.IsFailedPrecondition(err),
		"error is not FailedPrecondition",
	)
	assert.ErrorContains(t, err, "unary invoke: ")
}
