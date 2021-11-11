package cmd

import (
	"context"
	"fmt"
	pb "github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/cli"
	"github.com/DuarteMRAlves/maestro/internal/cli/display/table"
	"github.com/spf13/cobra"
	"io"
	"log"
	"strings"
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
		displayAssets(assets)
	},
}

func displayAssets(assets []*pb.Asset) {
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
	
	t.AddColumn(IdText, ids)
	t.AddColumn(NameText, names)
	t.AddColumn(ImageText, images)
	fmt.Print(t)
	//sb := &strings.Builder{}
	//buildTitle(sb)
	//for _, a := range assets {
	//	buildAsset(sb, a)
	//}
	//fmt.Print(sb.String())
}

func buildTitle(sb *strings.Builder) {
	idPad := asset.IdSize - len(IdText) + colPad
	titleLen := len(IdText) + idPad + len(NameText) + 1 // Last for \n
	sb.Grow(titleLen)
	sb.WriteString(IdText)
	for i := 0; i < idPad; i++ {
		sb.WriteByte(' ')
	}
	sb.WriteString(NameText)
	sb.WriteByte('\n')
}

func buildAsset(sb *strings.Builder, a *pb.Asset) {
	id := a.Id.Val
	name := a.Name
	idPad := asset.IdSize - len(id) + colPad
	assetLen := len(id) + idPad + len(name) + 1 // Last for \n
	sb.Grow(assetLen)
	sb.WriteString(id)
	for i := 0; i < idPad; i++ {
		sb.WriteByte(' ')
	}
	sb.WriteString(name)
	sb.WriteByte('\n')
}

func init() {
	assetCmd.AddCommand(listAssetsCmd)
}
