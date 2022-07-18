package grpcw

import (
	"context"
	"errors"
	"fmt"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/DuarteMRAlves/maestro/internal/retry"
	"github.com/DuarteMRAlves/maestro/test/protobuf/unit"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestReflectionClient_SlowServerStartup(t *testing.T) {
	var backoff retry.ExponentialBackoff

	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen: %s", err)
	}
	addr := lis.Addr().String()
	start, stop := startTestResolverServer(t, lis, true)
	defer stop()
	go func() {
		// Delay invocation to allow the resolver to start.
		time.Sleep(3 * time.Second)
		start()
	}()

	r, err := NewReflectionResolver(5*time.Second, backoff, testLogger{})
	if err != nil {
		t.Fatalf("create resolver error: %s", err)
	}

	methodAddr := fmt.Sprintf("%s/*/Unary", addr)
	fmt.Println("Resolving")
	_, err = r.Resolve(context.Background(), methodAddr)
	if err != nil {
		t.Fatalf("resolve error: %s", err)
	}
}

func TestReflectionClient_ListServices(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen: %s", err)
	}
	addr := lis.Addr().String()
	start, stop := startTestResolverServer(t, lis, true)
	defer stop()
	go start()

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

	m := ReflectionResolver{timeout: 5 * time.Second}
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
	start, stop := startTestResolverServer(t, lis, false)
	defer stop()
	go start()

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

	m := ReflectionResolver{timeout: 5 * time.Second}
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
	start, stop := startTestResolverServer(t, lis, true)
	defer stop()
	go start()

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
	var backoff retry.ExponentialBackoff
	m, err := NewReflectionResolver(5*time.Second, backoff, nil)
	if err != nil {
		t.Fatalf("create resolver: %v", err)
	}
	serv, err := m.resolveService(ctx, conn, serviceName)
	if err != nil {
		t.Fatalf("resolve service: %s", err)
	}
	assertTestService(t, serv)
}

func assertTestService(t *testing.T, descriptor protoreflect.ServiceDescriptor) {
	methods := descriptor.Methods()
	if diff := cmp.Diff(4, methods.Len()); diff != "" {
		t.Fatalf("number of methods mismatch:\n%s", diff)
	}

	names := []protoreflect.FullName{
		"unit.MethodLoaderTestService.Unary",
		"unit.MethodLoaderTestService.ClientStream",
		"unit.MethodLoaderTestService.ServerStream",
		"unit.MethodLoaderTestService.BidiStream",
	}
	for i := 0; i < methods.Len(); i++ {
		m := methods.Get(i)
		foundName := false
		for _, n := range names {
			if n == m.FullName() {
				foundName = true
			}
		}
		if !foundName {
			t.Fatalf("unknown method name '%v'", m.FullName())
		}
		assertRequestType(t, m.Input())
		assertReplyType(t, m.Output())
	}
}

func assertRequestType(t *testing.T, descriptor protoreflect.MessageDescriptor) {
	fields := descriptor.Fields()
	stringField := fields.ByName("stringField")
	if diff := cmp.Diff(protowire.Number(1), stringField.Number()); diff != "" {
		t.Fatalf("stringField number mismatch:\n%s", diff)
	}
	if diff := cmp.Diff(protoreflect.StringKind, stringField.Kind()); diff != "" {
		t.Fatalf("stringField type mismatch:\n%s", diff)
	}

	repeatedField := fields.ByName("repeatedField")
	if diff := cmp.Diff(protowire.Number(2), repeatedField.Number()); diff != "" {
		t.Fatalf("repeatedField number mismatch:\n%s", diff)
	}
	if diff := cmp.Diff(protoreflect.Int64Kind, repeatedField.Kind()); diff != "" {
		t.Fatalf("repeatedField type mismatch:\n%s", diff)
	}
	if diff := cmp.Diff(protoreflect.Repeated, repeatedField.Cardinality()); diff != "" {
		t.Fatalf("repeatedField cardinality mismatch:\n%s", diff)
	}

	repeatedInnerMsg := fields.ByName("repeatedInnerMsg")
	if diff := cmp.Diff(protowire.Number(3), repeatedInnerMsg.Number()); diff != "" {
		t.Fatalf("repeatedInnerMsg number mismatch:\n%s", diff)
	}
	if diff := cmp.Diff(protoreflect.MessageKind, repeatedInnerMsg.Kind()); diff != "" {
		t.Fatalf("repeatedInnerMsg type mismatch:\n%s", diff)
	}
	if diff := cmp.Diff(protoreflect.Repeated, repeatedInnerMsg.Cardinality()); diff != "" {
		t.Fatalf("repeatedInnerMsg cardinality mismatch:\n%s", diff)
	}

	innerType := repeatedInnerMsg.Message()
	if innerType == nil {
		t.Fatalf("inner type is nil")
	}
	assertInnerMessageType(t, innerType)
}

func assertReplyType(t *testing.T, descriptor protoreflect.MessageDescriptor) {
	fields := descriptor.Fields()
	doubleField := fields.ByName("doubleField")
	if diff := cmp.Diff(protowire.Number(1), doubleField.Number()); diff != "" {
		t.Fatalf("doubleField number mismatch:\n%s", diff)
	}
	if diff := cmp.Diff(protoreflect.DoubleKind, doubleField.Kind()); diff != "" {
		t.Fatalf("doubleField type mismatch:\n%s", diff)
	}

	innerMsg := fields.ByName("innerMsg")
	if diff := cmp.Diff(protowire.Number(2), innerMsg.Number()); diff != "" {
		t.Fatalf("innerMsg number mismatch:\n%s", diff)
	}
	if diff := cmp.Diff(protoreflect.MessageKind, innerMsg.Kind()); diff != "" {
		t.Fatalf("innerMsg type mismatch:\n%s", diff)
	}
	innerType := innerMsg.Message()
	if innerType == nil {
		t.Fatalf("inner type is nil")
	}
	assertInnerMessageType(t, innerType)
}

func assertInnerMessageType(t *testing.T, descriptor protoreflect.MessageDescriptor) {
	fields := descriptor.Fields()
	repeatedString := fields.ByName("repeatedString")
	if diff := cmp.Diff(protowire.Number(1), repeatedString.Number()); diff != "" {
		t.Fatalf("repeatedString number mismatch:\n%s", diff)
	}
	if diff := cmp.Diff(protoreflect.StringKind, repeatedString.Kind()); diff != "" {
		t.Fatalf("repeatedString type mismatch:\n%s", diff)
	}
	if diff := cmp.Diff(protoreflect.Repeated, repeatedString.Cardinality()); diff != "" {
		t.Fatalf("repeatedString cardinality mismatch:\n%s", diff)
	}
}

func TestReflectionClient_ResolveServiceNoReflection(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen: %s", err)
	}
	addr := lis.Addr().String()
	start, stop := startTestResolverServer(t, lis, false)
	defer stop()
	go start()

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
	m := ReflectionResolver{timeout: 5 * time.Second}
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
	start, stop := startTestResolverServer(t, lis, true)
	defer stop()
	go start()

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
	m := ReflectionResolver{timeout: 5 * time.Second}
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

type testResolverService struct {
	unit.UnimplementedMethodLoaderTestServiceServer
}

func (s *testResolverService) Unary(
	_ context.Context,
	_ *unit.MethodLoaderRequest,
) (*unit.MethodLoaderReply, error) {
	panic("Not implemented should not be called.")
}

func startTestResolverServer(
	t *testing.T,
	lis net.Listener,
	reflectionFlag bool,
) (func(), func()) {
	testServer := grpc.NewServer()
	unit.RegisterMethodLoaderTestServiceServer(testServer, &testResolverService{})

	if reflectionFlag {
		reflection.Register(testServer)
	}

	start := func() {
		if err := testServer.Serve(lis); err != nil {
			t.Errorf("test server: %s", err)
		}
	}
	stop := testServer.Stop
	return start, stop
}

type testLogger struct{}

func (t testLogger) Debugf(format string, args ...any) {
	fmt.Printf(format, args...)
}

func (t testLogger) Infof(format string, args ...any) {
	fmt.Printf(format, args...)
}
