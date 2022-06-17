package grpc

import (
	"context"
	"errors"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/DuarteMRAlves/maestro/test/protobuf/unit"
	protocdesc "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/google/go-cmp/cmp"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func TestReflectionClient_ListServices(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen: %s", err)
	}
	addr := lis.Addr().String()
	testServer := startServer(t, lis, true)
	defer testServer.GracefulStop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("dial error: %s", err)
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			t.Errorf("close connection: %s", err)
		}
	}(conn)

	m := ReflectionMethodLoader{timeout: 5 * time.Second}
	services, err := m.listServices(ctx, conn)
	if err != nil {
		t.Fatalf("list services: %s", err)
	}

	if diff := cmp.Diff(1, len(services)); diff != "" {
		t.Fatalf("mismatch on number of services:\n%s", diff)
	}
	counts := map[string]int{"unit.MethodLoaderTestService": 0}
	for _, s := range services {
		_, serviceExists := counts[string(s)]
		if !serviceExists {
			t.Fatalf("unexpected service %s", s)
		}
		counts[string(s)]++
	}
	for service, count := range counts {
		if diff := cmp.Diff(1, count); diff != "" {
			t.Fatalf("mismatch service %s occurences:\n%s", service, diff)
		}
	}
}

func TestReflectionClient_ListServicesNoReflection(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen: %s", err)
	}
	addr := lis.Addr().String()
	testServer := startServer(t, lis, false)
	defer testServer.GracefulStop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("dial error: %s", err)
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			t.Errorf("close connection: %s", err)
		}
	}(conn)

	m := ReflectionMethodLoader{timeout: 5 * time.Second}
	services, err := m.listServices(ctx, conn)
	if err == nil {
		t.Fatalf("expected non nil error at listServices")
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
	if services != nil {
		t.Fatalf("services are not nil")
	}
}

func TestReflectionClient_ResolveService_TestService(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen: %s", err)
	}
	addr := lis.Addr().String()
	testServer := startServer(t, lis, true)
	defer testServer.GracefulStop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("dial error: %s", err)
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			t.Errorf("close connection: %s", err)
		}
	}(conn)

	serviceName := Service("unit.MethodLoaderTestService")
	m := ReflectionMethodLoader{timeout: 5 * time.Second}
	serv, err := m.resolveService(ctx, conn, serviceName)
	if err != nil {
		t.Fatalf("resolve service: %s", err)
	}
	assertTestService(t, serv)
}

func assertTestService(t *testing.T, descriptor *desc.ServiceDescriptor) {
	methods := descriptor.GetMethods()
	if diff := cmp.Diff(4, len(methods)); diff != "" {
		t.Fatalf("number of methods mismatch:\n%s", diff)
	}

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
	if diff := cmp.Diff(int32(1), stringField.GetNumber()); diff != "" {
		t.Fatalf("stringField number mismatch:\n%s", diff)
	}
	if diff := cmp.Diff(
		protocdesc.FieldDescriptorProto_TYPE_STRING, stringField.GetType(),
	); diff != "" {
		t.Fatalf("stringField type mismatch:\n%s", diff)
	}

	repeatedField := descriptor.FindFieldByName("repeatedField")
	if diff := cmp.Diff(int32(2), repeatedField.GetNumber()); diff != "" {
		t.Fatalf("repeatedField number mismatch:\n%s", diff)
	}
	if diff := cmp.Diff(
		protocdesc.FieldDescriptorProto_TYPE_INT64, repeatedField.GetType(),
	); diff != "" {
		t.Fatalf("repeatedField type mismatch:\n%s", diff)
	}
	if !repeatedField.IsRepeated() {
		t.Fatalf("repeatedField is not repeated")
	}

	repeatedInnerMsg := descriptor.FindFieldByName("repeatedInnerMsg")
	if diff := cmp.Diff(int32(3), repeatedInnerMsg.GetNumber()); diff != "" {
		t.Fatalf("repeatedInnerMsg number mismatch:\n%s", diff)
	}
	if diff := cmp.Diff(
		protocdesc.FieldDescriptorProto_TYPE_MESSAGE, repeatedInnerMsg.GetType(),
	); diff != "" {
		t.Fatalf("repeatedInnerMsg type mismatch:\n%s", diff)
	}
	if !repeatedInnerMsg.IsRepeated() {
		t.Fatalf("repeatedInnerMsg is not repeated")
	}

	innerType := repeatedInnerMsg.GetMessageType()
	if innerType == nil {
		t.Fatalf("inner type is nil")
	}
	assertInnerMessageType(t, innerType)
}

func assertReplyType(t *testing.T, descriptor *desc.MessageDescriptor) {
	doubleField := descriptor.FindFieldByName("doubleField")
	if diff := cmp.Diff(int32(1), doubleField.GetNumber()); diff != "" {
		t.Fatalf("doubleField number mismatch:\n%s", diff)
	}
	if diff := cmp.Diff(
		protocdesc.FieldDescriptorProto_TYPE_DOUBLE, doubleField.GetType(),
	); diff != "" {
		t.Fatalf("doubleField type mismatch:\n%s", diff)
	}

	innerMsg := descriptor.FindFieldByName("innerMsg")
	if diff := cmp.Diff(int32(2), innerMsg.GetNumber()); diff != "" {
		t.Fatalf("innerMsg number mismatch:\n%s", diff)
	}
	if diff := cmp.Diff(
		protocdesc.FieldDescriptorProto_TYPE_MESSAGE, innerMsg.GetType(),
	); diff != "" {
		t.Fatalf("innerMsg type mismatch:\n%s", diff)
	}
	innerType := innerMsg.GetMessageType()
	if innerType == nil {
		t.Fatalf("inner type is nil")
	}
	assertInnerMessageType(t, innerType)
}

func assertInnerMessageType(t *testing.T, descriptor *desc.MessageDescriptor) {
	repeatedString := descriptor.FindFieldByName("repeatedString")
	if diff := cmp.Diff(int32(1), repeatedString.GetNumber()); diff != "" {
		t.Fatalf("repeatedString number mismatch:\n%s", diff)
	}
	if diff := cmp.Diff(
		protocdesc.FieldDescriptorProto_TYPE_STRING, repeatedString.GetType(),
	); diff != "" {
		t.Fatalf("repeatedString type mismatch:\n%s", diff)
	}
	if !repeatedString.IsRepeated() {
		t.Fatalf("repeatedString is not repeated")
	}
}

func TestReflectionClient_ResolveServiceNoReflection(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen: %s", err)
	}
	addr := lis.Addr().String()
	testServer := startServer(t, lis, false)
	defer testServer.GracefulStop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("dial: %s", err)
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			t.Errorf("close connection: %s", err)
		}
	}(conn)

	serviceName := Service("pb.TestService")
	m := ReflectionMethodLoader{timeout: 5 * time.Second}
	serv, err := m.resolveService(ctx, conn, serviceName)
	if err == nil {
		t.Fatalf("expected non nil error at resolveService")
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
	if serv != nil {
		t.Fatalf("serv is not nil")
	}
}

func TestReflectionClient_ResolveServiceUnknownService(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen: %s", err)
	}
	addr := lis.Addr().String()
	testServer := startServer(t, lis, true)
	defer testServer.GracefulStop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("dial: %s", err)
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			t.Errorf("close connection: %s", err)
		}
	}(conn)

	serviceName := Service("pb.UnknownService")
	m := ReflectionMethodLoader{timeout: 5 * time.Second}
	serv, err := m.resolveService(ctx, conn, serviceName)
	if err == nil {
		t.Fatalf("expected non nil error at listServices")
	}

	var notFoundErr interface{ NotFound() }
	if !errors.As(err, &notFoundErr) {
		t.Fatalf("error does not implement not found")
	}
	var nf *serviceNotFound
	if !errors.As(err, &nf) {
		format := "Wrong error type: expected %s, got %s"
		t.Fatalf(format, reflect.TypeOf(nf), reflect.TypeOf(err))
	}
	expError := &serviceNotFound{srv: string(serviceName)}
	cmpOpts := cmp.AllowUnexported(serviceNotFound{})
	if diff := cmp.Diff(expError, nf, cmpOpts); diff != "" {
		t.Fatalf("error mismatch:\n%s", diff)
	}
	if serv != nil {
		t.Fatalf("serv is not nil")
	}
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
		if err := testServer.Serve(lis); err != nil {
			t.Errorf("test server: %s", err)
		}
	}()
	return testServer
}
