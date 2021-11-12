package cmd

import (
	"context"
	"fmt"
	pb "github.com/DuarteMRAlves/maestro/api/pb"
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
)

const (
	colPad     = 2
	minColSize = 15
)

var listAssetsCmd = &cobra.Command{
	Use:   "list",
	Short: "maestro-cli asset list allows you to list created assets",
	Run: func(cmd *cobra.Command, args []string) {
		conn := cli.NewConnection(addr)
		defer conn.Close()

		c := pb.NewAssetManagementClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		stream, err := c.List(ctx, &pb.SearchQuery{})
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
	},
}

func displayAssets(assets []*pb.Asset) error {
	numAssets := len(assets)
	ids := make([]string, 0, numAssets)
	names := make([]string, 0, numAssets)
	images := make([]string, 0, numAssets)
	for _, a := range assets {
		ids = append(ids, a.Id.Val)
		names = append(names, a.Name)
		images = append(images, a.Image)
	}

	t := table.NewBuilder().
		WithPadding(colPad).
		WithMinColSize(minColSize).
		Build()

	if err := t.AddColumn(IdText, ids); err != nil {
		return err
	}
	if err := t.AddColumn(NameText, names); err != nil {
		return err
	}
	if err := t.AddColumn(ImageText, images); err != nil {
		return err
	}
	fmt.Print(t)
	return nil
}

func init() {
	assetCmd.AddCommand(listAssetsCmd)
}
