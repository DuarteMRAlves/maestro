package client

import (
	"context"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"time"
)

func CreateAsset(a *pb.Asset, addr string) error {
	conn := NewConnection(addr)
	defer conn.Close()

	c := pb.NewAssetManagementClient(conn)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second)
	defer cancel()

	_, err := c.Create(ctx, a)

	return err
}

func CreateStage(s *pb.Stage, addr string) error {
	conn := NewConnection(addr)
	defer conn.Close()

	c := pb.NewStageManagementClient(conn)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second)
	defer cancel()

	_, err := c.Create(ctx, s)

	return err
}
