package grpcexecute

import (
	"github.com/DuarteMRAlves/maestro/internal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"testing"
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

func createOrchName(t *testing.T, name string) internal.OrchestrationName {
	orchName, err := internal.NewOrchestrationName(name)
	if err != nil {
		t.Fatalf("create orchestration name %s: %s", name, err)
	}
	return orchName
}

func createStageName(t *testing.T, name string) internal.StageName {
	stageName, err := internal.NewStageName(name)
	if err != nil {
		t.Fatalf("create stage name %s: %s", name, err)
	}
	return stageName
}

func createLinkName(t *testing.T, name string) internal.LinkName {
	linkName, err := internal.NewLinkName(name)
	if err != nil {
		t.Fatalf("create link name %s: %s", name, err)
	}
	return linkName
}

func createMethodContext(addr internal.Address) internal.MethodContext {
	var (
		emptyService internal.Service
		emptyMethod  internal.Method
	)
	return internal.NewMethodContext(addr, emptyService, emptyMethod)
}
