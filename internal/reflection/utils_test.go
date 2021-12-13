package reflection

import (
	"github.com/DuarteMRAlves/maestro/tests/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gotest.tools/v3/assert"
	"net"
	"testing"
)

type service struct {
	pb.UnimplementedTestServiceServer
	pb.UnimplementedExtraServiceServer
}

func startServer(
	t *testing.T,
	lis net.Listener,
	reflectionFlag bool,
) *grpc.Server {
	testServer := grpc.NewServer()
	pb.RegisterTestServiceServer(testServer, &service{})
	pb.RegisterExtraServiceServer(testServer, &service{})

	if reflectionFlag {
		reflection.Register(testServer)
	}

	go func() {
		err := testServer.Serve(lis)
		assert.NilError(t, err, "test server error")
	}()
	return testServer
}
