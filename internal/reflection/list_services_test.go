package reflection

import (
	"context"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/testutil"
	"google.golang.org/grpc"
	"gotest.tools/v3/assert"
	"testing"
	"time"
)

func TestClient_ListServices(t *testing.T) {
	lis := testutil.ListenAvailablePort(t)
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

func TestClient_ListServicesNoReflection(t *testing.T) {
	lis := testutil.ListenAvailablePort(t)
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

	c := NewClient(ctx, conn)
	services, err := c.ListServices()

	assert.Assert(t, errdefs.IsFailedPrecondition(err), "list services error")
	assert.ErrorContains(t, err, "list services:")
	assert.Assert(t, services == nil, "services is not nil")
}
