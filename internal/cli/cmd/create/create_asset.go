package create

import (
	"context"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/cli"
	"github.com/spf13/cobra"
	"log"
	"time"
)

const (
	nameFlag   = "name"
	nameUsage  = "Name to give the asset (required)"
	imageFlag  = "image"
	imageUsage = "Docker image for the asset"
)

var createAssetOpts = &struct {
	name  string
	image string
}{}

func NewCmdCreateAsset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "asset",
		Short: "create a new Asset",
		Run:   runCreateAsset,
	}

	cmd.Flags().StringVar(&createAssetOpts.name, nameFlag, "", nameUsage)
	cmd.MarkFlagRequired(nameFlag)

	cmd.Flags().StringVar(&createAssetOpts.image, imageFlag, "", imageUsage)

	return cmd
}

func runCreateAsset(_ *cobra.Command, _ []string) {
	conn := cli.NewConnection(createOpts.addr)
	defer conn.Close()

	c := pb.NewAssetManagementClient(conn)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second)
	defer cancel()
	asset, err := c.Create(
		ctx,
		&pb.Asset{
			Name:  createAssetOpts.name,
			Image: createAssetOpts.image,
		})

	if err != nil {
		log.Fatalf("Unable to create asset: %v", err)
	}
	fmt.Printf("created asset with identifier '%v'\n", asset.Val)
}
