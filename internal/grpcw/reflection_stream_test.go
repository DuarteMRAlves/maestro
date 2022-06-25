package grpcw

import (
	"context"
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/DuarteMRAlves/maestro/test/protobuf/unit"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
)

func TestListServiceNames(t *testing.T) {
	addr, stop := startServer(true)
	defer stop()

	conn, close := dialServer(addr)
	defer close()

	s, err := newBlockingReflectionStream(context.Background(), conn)
	if err != nil {
		t.Fatalf("create reflection client: %s", err)
	}

	services, err := s.listServiceNames()
	if err != nil {
		t.Fatalf("list services: %s", err)
	}

	if diff := cmp.Diff(2, len(services)); diff != "" {
		t.Fatalf("mismatch on number of services:\n%s", diff)
	}
	counts := map[string]int{
		"unit.MethodLoaderTestService":             0,
		"grpc.reflection.v1alpha.ServerReflection": 0,
	}
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

func TestFilesForSymbol(t *testing.T) {
	addr, stop := startServer(true)
	defer stop()

	conn, close := dialServer(addr)
	defer close()

	s, err := newBlockingReflectionStream(context.Background(), conn)
	if err != nil {
		t.Fatalf("create reflection client: %s", err)
	}

	fds, err := s.filesForSymbol("unit.MethodLoaderTestService")
	if err != nil {
		t.Fatalf("files for symbol: %s", err)
	}

	if diff := cmp.Diff(1, len(fds)); diff != "" {
		t.Fatalf("number of files mismatch:\n%s", diff)
	}

	expDesc := unit.File_method_loader_proto
	expProto := protodesc.ToFileDescriptorProto(expDesc)
	expBytes, err := proto.Marshal(expProto)
	if err != nil {
		t.Fatalf("marshal expected proto: %s", err)
	}

	if diff := cmp.Diff(expBytes, fds[0]); diff != "" {
		t.Fatalf("file descriptor proto mismatch:\n%s", diff)
	}
}

func TestNoReflection(t *testing.T) {
	addr, stop := startServer(false)
	defer stop()

	conn, close := dialServer(addr)
	defer close()

	s, err := newBlockingReflectionStream(context.Background(), conn)
	if err != nil {
		t.Fatalf("create reflection client: %s", err)
	}

	services, err := s.listServiceNames()
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

type testService struct {
	unit.UnimplementedMethodLoaderTestServiceServer
}

func (s *testService) Unary(
	_ context.Context,
	_ *unit.MethodLoaderRequest,
) (*unit.MethodLoaderReply, error) {
	panic("Not implemented should not be called.")
}

func startServer(activateReflection bool) (string, func()) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(fmt.Sprintf("failed to listen: %s", err))
	}
	addr := lis.Addr().String()

	testServer := grpc.NewServer()
	unit.RegisterMethodLoaderTestServiceServer(testServer, &testService{})

	if activateReflection {
		reflection.Register(testServer)
	}

	go func() {
		if err := testServer.Serve(lis); err != nil {
			panic(fmt.Sprintf("test server: %s", err))
		}
	}()
	return addr, testServer.GracefulStop
}

func dialServer(addr string) (*grpc.ClientConn, func()) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		panic(fmt.Sprintf("dial error: %s", err))
	}
	closeFunc := func() {
		if err := conn.Close(); err != nil {
			panic(fmt.Sprintf("close connection: %s", err))
		}
	}
	return conn, closeFunc
}
