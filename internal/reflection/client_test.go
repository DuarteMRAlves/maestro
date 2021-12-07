package reflection

import (
	"context"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/tests/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gotest.tools/v3/assert"
	"net"
	"testing"
	"time"
)

type service struct {
	pb.UnimplementedTestServiceServer
	pb.UnimplementedExtraServiceServer
}

func TestClient_ListServices(t *testing.T) {
	addr := "localhost:50051"
	lis, err := net.Listen("tcp", addr)
	assert.NilError(t, err, "listen error")

	testServer := grpc.NewServer()
	pb.RegisterTestServiceServer(testServer, &service{})
	pb.RegisterExtraServiceServer(testServer, &service{})
	reflection.Register(testServer)

	go func() {
		err = testServer.Serve(lis)
		assert.NilError(t, err, "test server error")
	}()
	defer testServer.GracefulStop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	assert.NilError(t, err, "dial error")
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		assert.NilError(t, err, "close connection")
	}(conn)

	c := NewClient(ctx, conn)
	services, err := c.ListServices()
	assert.NilError(t, err, "list services error")

	assert.Equal(t, 2, len(services), "number of services")
	counts := map[string]int{
		"pb.TestService":  0,
		"pb.ExtraService": 0,
	}
	for _, s := range services {
		_, serviceExists := counts[s]
		assert.Assert(t, serviceExists, "unexpected service %v", s)
		counts[s]++
	}
	for service, count := range counts {
		assert.Equal(
			t,
			1,
			count,
			"service %v did not appear only once",
			service)
	}
}

func TestClient_ListServicesUnavailable(t *testing.T) {
	addr := "localhost:50051"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	assert.NilError(t, err, "dial error")
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		assert.NilError(t, err, "close connection")
	}(conn)

	c := NewClient(ctx, conn)
	services, err := c.ListServices()

	assert.Assert(t, errdefs.IsUnavailable(err), "list services error")
	assert.ErrorContains(t, err, "list services:")
	assert.Assert(t, services == nil, "services is not nil")
}

func TestClient_ListServicesNoReflection(t *testing.T) {
	addr := "localhost:50051"
	lis, err := net.Listen("tcp", addr)
	assert.NilError(t, err, "listen error")

	testServer := grpc.NewServer()
	pb.RegisterTestServiceServer(testServer, &service{})
	pb.RegisterExtraServiceServer(testServer, &service{})

	go func() {
		err = testServer.Serve(lis)
		assert.NilError(t, err, "test server error")
	}()
	defer testServer.GracefulStop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	assert.NilError(t, err, "dial error")
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		assert.NilError(t, err, "close connection")
	}(conn)

	c := NewClient(ctx, conn)
	services, err := c.ListServices()

	assert.Assert(t, errdefs.IsFailedPrecondition(err), "list services error")
	assert.ErrorContains(t, err, "list services:")
	assert.Assert(t, services == nil, "services is not nil")
}
