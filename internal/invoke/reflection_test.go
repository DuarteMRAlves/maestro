package invoke

import (
	"context"
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	protocdesc "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"google.golang.org/grpc"
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

	listFn := listServices(conn)
	services, err := listFn(ctx)
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

	listFn := listServices(conn)
	services, err := listFn(ctx)

	assert.Assert(t, errdefs.IsFailedPrecondition(err), "list services error")
	assert.ErrorContains(t, err, "list services:")
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

	serviceName, err := domain.NewService("unit.TestService")
	assert.NilError(t, err, "create service name")
	resolveFn := resolveService(conn)
	serv, err := resolveFn(ctx, serviceName)
	assert.NilError(t, err, "resolve service error")
	assertTestService(t, serv)
}

func assertTestService(t *testing.T, descriptor Service) {
	methods := descriptor.Methods()
	assert.Equal(t, 4, len(methods), "number of methods")

	for _, m := range methods {
		switch m.FullMethod().Unwrap() {
		case "/unit.TestService/Unary":
			assert.Assert(t, m.IsUnary())
		case "/unit.TestService/ClientStream":
			assert.Assert(t, !m.IsUnary())
		case "/unit.TestService/ServerStream":
			assert.Assert(t, !m.IsUnary())
		case "/unit.TestService/BidiStream":
			assert.Assert(t, !m.IsUnary())
		default:
			t.Fatalf("unknown method name '%v'", m.FullMethod())
		}
		assertRequestType(t, m.Input())
		assertReplyType(t, m.Output())
	}
}

func assertRequestType(t *testing.T, messageDesc MessageDescriptor) {
	empty := messageDesc.MessageGenerator()()
	dyn, err := dynamic.AsDynamicMessage(empty.GrpcMessage())
	assert.NilError(t, err, "dynamic request")
	descriptor := dyn.GetMessageDescriptor()
	assert.Equal(t, "Request", descriptor.GetName())

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

func assertReplyType(t *testing.T, messageDesc MessageDescriptor) {
	empty := messageDesc.MessageGenerator()()
	dyn, err := dynamic.AsDynamicMessage(empty.GrpcMessage())
	assert.NilError(t, err, "dynamic request")
	descriptor := dyn.GetMessageDescriptor()
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

	serviceName, err := domain.NewService("unit.ExtraService")
	assert.NilError(t, err, "create service name")
	resolveFn := resolveService(conn)
	serv, err := resolveFn(ctx, serviceName)
	assert.NilError(t, err, "resolve service error")
	assertExtraService(t, serv)
}

func assertExtraService(t *testing.T, descriptor Service) {
	methods := descriptor.Methods()
	assert.Equal(t, 1, len(methods), "number of methods")

	m := methods[0]
	assert.Assert(t, m.IsUnary())
	assertExtraRequestType(t, m.Input())
	assertExtraReplyType(t, m.Output())
}

func assertExtraRequestType(t *testing.T, messageDesc MessageDescriptor) {
	empty := messageDesc.MessageGenerator()()
	dyn, err := dynamic.AsDynamicMessage(empty.GrpcMessage())
	assert.NilError(t, err, "dynamic request")
	descriptor := dyn.GetMessageDescriptor()
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

func assertExtraReplyType(t *testing.T, messageDesc MessageDescriptor) {
	empty := messageDesc.MessageGenerator()()
	dyn, err := dynamic.AsDynamicMessage(empty.GrpcMessage())
	assert.NilError(t, err, "dynamic request")
	descriptor := dyn.GetMessageDescriptor()
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

	serviceName, err := domain.NewService("pb.TestService")
	assert.NilError(t, err, "create service name")
	resolveFn := resolveService(conn)
	serv, err := resolveFn(ctx, serviceName)

	assert.Assert(t, errdefs.IsFailedPrecondition(err), "resolve service error")
	assert.ErrorContains(t, err, "resolve service: ")
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

	serviceName, err := domain.NewService("pb.UnknownService")
	assert.NilError(t, err, "create service name")
	resolveFn := resolveService(conn)
	serv, err := resolveFn(ctx, serviceName)

	assert.Assert(t, errdefs.IsNotFound(err), "resolve service error")
	expectedMsg := "resolve service: Service not found: pb.UnknownService"
	assert.Error(t, err, expectedMsg)
	assert.Assert(t, serv == nil, "service is not nil")
}
