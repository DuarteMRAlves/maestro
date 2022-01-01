package testutil

import (
	"github.com/DuarteMRAlves/maestro/tests/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gotest.tools/v3/assert"
	"net"
	"testing"
)

type testService struct {
	pb.UnimplementedTestServiceServer
	pb.UnimplementedExtraServiceServer
}

func StartTestServer(
	t *testing.T,
	lis net.Listener,
	registerTest bool,
	registerExtra bool,
) *grpc.Server {
	testServer := grpc.NewServer()
	if registerTest {
		pb.RegisterTestServiceServer(testServer, &testService{})
	}
	if registerExtra {
		pb.RegisterExtraServiceServer(testServer, &testService{})
	}

	reflection.Register(testServer)

	go func() {
		err := testServer.Serve(lis)
		assert.NilError(t, err, "test server error")
	}()
	return testServer
}
