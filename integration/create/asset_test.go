package create

import (
	"context"
	"github.com/DuarteMRAlves/maestro/internal/cli/client"
	"github.com/DuarteMRAlves/maestro/internal/cli/resources"
	"github.com/DuarteMRAlves/maestro/internal/server"
	"gotest.tools/v3/assert"
	"net"
	"testing"
	"time"
)

func TestCreateAssetFromCli(t *testing.T) {
	address := "localhost:50051"
	lis, err := net.Listen("tcp", address)
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}

	s := server.NewBuilder().WithGrpc().Build()

	go func() {
		if err := s.ServeGrpc(lis); err != nil {
			t.Fatalf("Failed to serve: %v", err)
		}
	}()
	defer s.GracefulStopGrpc()

	asset := &resources.AssetResource{
		Name:  "asset",
		Image: "image",
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second)
	defer cancel()

	err = client.CreateAsset(ctx, asset, address)
	assert.NilError(t, err, "error is %v", err)
}