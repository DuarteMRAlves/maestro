package grpc

import (
	"context"
	"errors"
	"github.com/DuarteMRAlves/maestro/internal"
	protocdesc "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gotest.tools/v3/assert"
	"net"
	"testing"
	"time"
)

func TestReflectionClient_ListServices(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:0")
	assert.NilError(t, err, "failed to listen")
	addr := lis.Addr().String()
	testServer := startServer(t, lis, true)
	defer testServer.GracefulStop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	assert.NilError(t, err, "dial error")
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		assert.NilError(t, err, "close connection")
	}(conn)

	services, err := listServices(ctx, conn)
	assert.NilError(t, err, "list services error")

	assert.Equal(t, 2, len(services), "number of services")
	counts := map[string]int{
		"unit.TestService":  0,
		"unit.ExtraService": 0,
	}
	for _, s := range services {
		_, serviceExists := counts[s.Unwrap()]
		assert.Assert(t, serviceExists, "unexpected service %v", s)
		counts[s.Unwrap()]++
	}
	for service, count := range counts {
		assert.Equal(
			t,
			1,
			count,
			"service %v did not appear only once",
			service,
		)
	}
}

func TestReflectionClient_ListServicesNoReflection(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:0")
	assert.NilError(t, err, "failed to listen")
	addr := lis.Addr().String()
	testServer := startServer(t, lis, false)
	defer testServer.GracefulStop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	assert.NilError(t, err, "dial error")
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		assert.NilError(t, err, "close connection")
	}(conn)

	services, err := listServices(ctx, conn)
	assert.Assert(t, err != nil)
	cause, ok := errors.Unwrap(err).(interface {
		GRPCStatus() *status.Status
	})
	st := cause.GRPCStatus()
	assert.Assert(t, ok, "correct type")
	assert.Equal(t, codes.Unimplemented, st.Code())
	assert.Assert(t, services == nil, "services is not nil")
}

func TestReflectionClient_ResolveService_TestService(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:0")
	assert.NilError(t, err, "failed to listen")
	addr := lis.Addr().String()
	testServer := startServer(t, lis, true)
	defer testServer.GracefulStop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	assert.NilError(t, err, "dial error")
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		assert.NilError(t, err, "close connection")
	}(conn)

	serviceName := internal.NewService("unit.TestService")
	serv, err := resolveService(ctx, conn, serviceName)
	assert.NilError(t, err, "resolve service error")
	assertTestService(t, serv)
}

func assertTestService(t *testing.T, descriptor *desc.ServiceDescriptor) {
	methods := descriptor.GetMethods()
	assert.Equal(t, 4, len(methods), "number of methods")

	names := []string{
		"unit.TestService.Unary",
		"unit.TestService.ClientStream",
		"unit.TestService.ServerStream",
		"unit.TestService.BidiStream",
	}
	for _, m := range methods {
		foundName := false
		for _, n := range names {
			if n == m.GetFullyQualifiedName() {
				foundName = true
			}
		}
		if !foundName {
			t.Fatalf("unknown method name '%v'", m.GetFullyQualifiedName())
		}
		assertRequestType(t, m.GetInputType())
		assertReplyType(t, m.GetOutputType())
	}
}

func assertRequestType(t *testing.T, descriptor *desc.MessageDescriptor) {
	stringField := descriptor.FindFieldByName("stringField")
	assert.Equal(t, int32(1), stringField.GetNumber())
	assert.Equal(
		t,
		protocdesc.FieldDescriptorProto_TYPE_STRING,
		stringField.GetType(),
	)

	repeatedField := descriptor.FindFieldByName("repeatedField")
	assert.Equal(t, int32(2), repeatedField.GetNumber())
	assert.Equal(
		t,
		protocdesc.FieldDescriptorProto_TYPE_INT64,
		repeatedField.GetType(),
	)
	assert.Assert(t, repeatedField.IsRepeated())

	repeatedInnerMsg := descriptor.FindFieldByName("repeatedInnerMsg")
	assert.Equal(t, int32(3), repeatedInnerMsg.GetNumber())
	assert.Equal(
		t,
		protocdesc.FieldDescriptorProto_TYPE_MESSAGE,
		repeatedInnerMsg.GetType(),
	)
	assert.Assert(t, repeatedInnerMsg.IsRepeated())

	innerType := repeatedInnerMsg.GetMessageType()
	assert.Assert(t, innerType != nil)
	assertInnerMessageType(t, innerType)
}

func assertReplyType(t *testing.T, descriptor *desc.MessageDescriptor) {
	doubleField := descriptor.FindFieldByName("doubleField")
	assert.Equal(t, int32(1), doubleField.GetNumber())
	assert.Equal(
		t,
		protocdesc.FieldDescriptorProto_TYPE_DOUBLE,
		doubleField.GetType(),
	)

	innerMsg := descriptor.FindFieldByName("innerMsg")
	assert.Equal(t, int32(2), innerMsg.GetNumber())
	assert.Equal(
		t,
		protocdesc.FieldDescriptorProto_TYPE_MESSAGE,
		innerMsg.GetType(),
	)

	innerType := innerMsg.GetMessageType()
	assert.Assert(t, innerType != nil)
	assertInnerMessageType(t, innerType)
}

func assertInnerMessageType(t *testing.T, descriptor *desc.MessageDescriptor) {
	repeatedString := descriptor.FindFieldByName("repeatedString")
	assert.Equal(t, int32(1), repeatedString.GetNumber())
	assert.Equal(
		t,
		protocdesc.FieldDescriptorProto_TYPE_STRING,
		repeatedString.GetType(),
	)
	assert.Assert(t, repeatedString.IsRepeated())
}

func TestReflectionClient_ResolveService_ExtraService(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:0")
	assert.NilError(t, err, "failed to listen")
	addr := lis.Addr().String()
	testServer := startServer(t, lis, true)
	defer testServer.GracefulStop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	assert.NilError(t, err, "dial error")
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		assert.NilError(t, err, "close connection")
	}(conn)

	serviceName := internal.NewService("unit.ExtraService")
	serv, err := resolveService(ctx, conn, serviceName)
	assert.NilError(t, err, "resolve service error")
	assertExtraService(t, serv)
}

func assertExtraService(t *testing.T, descriptor *desc.ServiceDescriptor) {
	methods := descriptor.GetMethods()
	assert.Equal(t, 1, len(methods), "number of methods")

	m := methods[0]
	assertExtraRequestType(t, m.GetInputType())
	assertExtraReplyType(t, m.GetOutputType())
}

func assertExtraRequestType(t *testing.T, descriptor *desc.MessageDescriptor) {
	repeatedStringField := descriptor.FindFieldByName("repeatedStringField")
	assert.Equal(t, int32(1), repeatedStringField.GetNumber())
	assert.Equal(
		t,
		protocdesc.FieldDescriptorProto_TYPE_STRING,
		repeatedStringField.GetType(),
	)
	assert.Assert(t, repeatedStringField.IsRepeated())

	innerMsg := descriptor.FindFieldByName("innerMsg")
	assert.Equal(t, int32(2), innerMsg.GetNumber())
	assert.Equal(
		t,
		protocdesc.FieldDescriptorProto_TYPE_MESSAGE,
		innerMsg.GetType(),
	)

	innerType := innerMsg.GetMessageType()
	assert.Assert(t, innerType != nil)
	assertExtraInnerMessageType(t, innerType)
}

func assertExtraReplyType(t *testing.T, descriptor *desc.MessageDescriptor) {
	oneOfs := descriptor.GetOneOfs()
	assert.Equal(t, 1, len(oneOfs))

	oneOf := oneOfs[0]
	assert.Equal(t, "oneOfField", oneOf.GetName())
	oneOfChoices := oneOf.GetChoices()
	assert.Equal(t, 2, len(oneOfChoices))

	oneOfChoice1 := oneOfChoices[0]
	assert.Equal(t, "intOpt", oneOfChoice1.GetName())
	assert.Equal(t, int32(1), oneOfChoice1.GetNumber())
	assert.Equal(
		t,
		protocdesc.FieldDescriptorProto_TYPE_INT64,
		oneOfChoice1.GetType(),
	)

	oneOfChoice2 := oneOfChoices[1]
	assert.Equal(t, "extraInnerMsg", oneOfChoice2.GetName())
	assert.Equal(t, int32(2), oneOfChoice2.GetNumber())
	assert.Equal(
		t,
		protocdesc.FieldDescriptorProto_TYPE_MESSAGE,
		oneOfChoice2.GetType(),
	)

	extraInnerMsg := oneOfChoice2.GetMessageType()
	assert.Assert(t, extraInnerMsg != nil)
	assertExtraInnerMessageType(t, extraInnerMsg)

	repeatedDoubleField := descriptor.FindFieldByName("repeatedDoubleField")
	assert.Equal(t, int32(3), repeatedDoubleField.GetNumber())
	assert.Equal(
		t,
		protocdesc.FieldDescriptorProto_TYPE_DOUBLE,
		repeatedDoubleField.GetType(),
	)
	assert.Assert(t, repeatedDoubleField.IsRepeated())
}

func assertExtraInnerMessageType(
	t *testing.T,
	descriptor *desc.MessageDescriptor,
) {
	repeatedString := descriptor.FindFieldByName("repeatedString")
	assert.Equal(t, int32(1), repeatedString.GetNumber())
	assert.Equal(
		t,
		protocdesc.FieldDescriptorProto_TYPE_STRING,
		repeatedString.GetType(),
	)
	assert.Assert(t, repeatedString.IsRepeated())
}

func TestReflectionClient_ResolveServiceNoReflection(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:0")
	assert.NilError(t, err, "failed to listen")
	addr := lis.Addr().String()
	testServer := startServer(t, lis, false)
	defer testServer.GracefulStop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	assert.NilError(t, err, "dial error")
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		assert.NilError(t, err, "close connection")
	}(conn)

	serviceName := internal.NewService("pb.TestService")
	serv, err := resolveService(ctx, conn, serviceName)

	assert.Assert(t, err != nil)
	cause, ok := errors.Unwrap(err).(interface {
		GRPCStatus() *status.Status
	})
	st := cause.GRPCStatus()
	assert.Assert(t, ok, "correct type")
	assert.Equal(t, codes.Unimplemented, st.Code())
	assert.Assert(t, serv == nil, "service is not nil")
}

func TestReflectionClient_ResolveServiceUnknownService(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:0")
	assert.NilError(t, err, "failed to listen")
	addr := lis.Addr().String()
	testServer := startServer(t, lis, true)
	defer testServer.GracefulStop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	assert.NilError(t, err, "dial error")
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		assert.NilError(t, err, "close connection")
	}(conn)

	serviceName := internal.NewService("pb.UnknownService")
	serv, err := resolveService(ctx, conn, serviceName)

	var notFound *internal.NotFound
	assert.Assert(t, errors.As(err, &notFound), "resolve service error")
	assert.Equal(t, "service", notFound.Type)
	assert.Equal(t, serviceName.Unwrap(), notFound.Ident)
	assert.Assert(t, serv == nil, "service is not nil")
}
