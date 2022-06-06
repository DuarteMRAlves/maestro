package grpcexecute

import (
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func createGrpcServer(
	t *testing.T, registerSrv func(grpc.ServiceRegistrar),
) (net.Addr, func(), func()) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen: %s", err)
	}
	s := grpc.NewServer()
	registerSrv(s)
	reflection.Register(s)

	start := func() {
		if err := s.Serve(lis); err != nil {
			t.Fatalf("Failed to server: %s", err)
		}
	}
	stop := func() {
		s.Stop()
	}
	return lis.Addr(), start, stop
}
