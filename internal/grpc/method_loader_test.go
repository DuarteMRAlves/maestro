package grpc

import (
	"context"
	"errors"
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/test/protobuf/unit"
	protocdesc "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
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

	assert.Equal(t, 1, len(services), "number of services")
	counts := map[string]int{"unit.MethodLoaderTestService": 0}
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

	serviceName := internal.NewService("unit.MethodLoaderTestService")
	serv, err := resolveService(ctx, conn, serviceName)
	assert.NilError(t, err, "resolve service error")
	assertTestService(t, serv)
}

func assertTestService(t *testing.T, descriptor *desc.ServiceDescriptor) {
	methods := descriptor.GetMethods()
	assert.Equal(t, 4, len(methods), "number of methods")

	names := []string{
		"unit.MethodLoaderTestService.Unary",
		"unit.MethodLoaderTestService.ClientStream",
		"unit.MethodLoaderTestService.ServerStream",
		"unit.MethodLoaderTestService.BidiStream",
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

type testService struct {
	unit.UnimplementedMethodLoaderTestServiceServer
}

func (s *testService) Unary(
	_ context.Context,
	_ *unit.MethodLoaderRequest,
) (*unit.MethodLoaderReply, error) {
	panic("Not implemented should not be called.")
}

func startServer(
	t *testing.T,
	lis net.Listener,
	reflectionFlag bool,
) *grpc.Server {
	testServer := grpc.NewServer()
	unit.RegisterMethodLoaderTestServiceServer(testServer, &testService{})

	if reflectionFlag {
		reflection.Register(testServer)
	}

	go func() {
		err := testServer.Serve(lis)
		assert.NilError(t, err, "test server error")
	}()
	return testServer
}
