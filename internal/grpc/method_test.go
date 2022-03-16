package grpc

import (
	"context"
	"errors"
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/test/protobuf/unit"
	"github.com/google/go-cmp/cmp/cmpopts"
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

	inDesc, err := newMessageDescriptor(&unit.Request{})
	assert.NilError(t, err, "create input message descriptor")
	outDesc, err := newMessageDescriptor(&unit.Reply{})
	assert.NilError(t, err, "create output message descriptor")

	method := newUnaryMethod("unit.TestService/Unary", inDesc, outDesc)

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

	pbReply := &unit.Reply{}
	err = grpcMsg.dynMsg.ConvertTo(pbReply)
	assert.NilError(t, err, "convert dynamic replay to message")

	cmpOpts := cmpopts.IgnoreUnexported(unit.Reply{}, unit.InnerMessage{})
	assert.DeepEqual(t, expectedReply, pbReply, cmpOpts)
}

func TestUnaryClient_Invoke_ErrorReturned(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:0")
	assert.NilError(t, err, "failed to listen")
	addr := lis.Addr().String()
	testServer := startServer(t, lis, false)
	defer testServer.Stop()

	inDesc, err := newMessageDescriptor(&unit.Request{})
	assert.NilError(t, err, "create input message descriptor")
	outDesc, err := newMessageDescriptor(&unit.Reply{})
	assert.NilError(t, err, "create output message descriptor")

	method := newUnaryMethod("unit.TestService/Unary", inDesc, outDesc)

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

	inDesc, err := newMessageDescriptor(&unit.ExtraRequest{})
	assert.NilError(t, err, "create input message descriptor")
	outDesc, err := newMessageDescriptor(&unit.ExtraReply{})
	assert.NilError(t, err, "create output message descriptor")

	method := newUnaryMethod("unit.ExtraService/ExtraMethod", inDesc, outDesc)

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
