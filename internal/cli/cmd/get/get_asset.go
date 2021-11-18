package get

import (
	"context"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/cli"
	"github.com/DuarteMRAlves/maestro/internal/cli/display/table"
	"github.com/spf13/cobra"
	"io"
	"log"
	"time"
)

const (
	IdText    = "ID"
	NameText  = "NAME"
	ImageText = "IMAGE"

	colPad     = 2
	minColSize = 15
)

func NewCmdGetAsset() *cobra.Command {
	return &cobra.Command{
		Use:   "asset",
		Short: "list one or more Assets",
		Run:   runGetAsset,
	}
}

func runGetAsset(_ *cobra.Command, _ []string) {
	conn := cli.NewConnection(createOpts.addr)
	defer conn.Close()

	c := pb.NewAssetManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := c.Get(ctx, &pb.Asset{})
	assets := make([]*pb.Asset, 0)
	if err != nil {
		log.Fatalf("Unable to create asset: %v", err)
	}
	for {
		a, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("unable to list assets: %v", err)
		}
		assets = append(assets, a)
	}
	if err := displayAssets(assets); err != nil {
		log.Fatalf("display assets: %v", err)
	}
}

func displayAssets(assets []*pb.Asset) error {
	numAssets := len(assets)
	names := make([]string, 0, numAssets)
	images := make([]string, 0, numAssets)
	for _, a := range assets {
		names = append(names, a.Name)
		images = append(images, a.Image)
	}

	t := table.NewBuilder().
		WithPadding(colPad).
		WithMinColSize(minColSize).
		Build()

	if err := t.AddColumn(NameText, names); err != nil {
		return err
	}
	if err := t.AddColumn(ImageText, images); err != nil {
		return err
	}
	fmt.Print(t)
	return nil
}
