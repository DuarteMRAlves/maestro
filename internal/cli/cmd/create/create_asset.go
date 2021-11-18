package create

import (
	"context"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/cli"
	"github.com/spf13/cobra"
	"time"
)

const (
	nameFlag   = "name"
	nameUsage  = "Name to give the asset (required)"
	imageFlag  = "image"
	imageUsage = "Docker image for the asset"
)

var createAssetOpts = &struct {
	image string
}{}

func NewCmdCreateAsset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "asset NAME",
		Short: "create a new Asset",
		Args:  cobra.ExactArgs(1),
		Run:   runCreateAsset,
	}

	cmd.Flags().StringVar(&createAssetOpts.image, imageFlag, "", imageUsage)

	return cmd
}

func runCreateAsset(_ *cobra.Command, args []string) {
	conn := cli.NewConnection(createOpts.addr)
	defer conn.Close()

	c := pb.NewAssetManagementClient(conn)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second)
	defer cancel()
	_, err := c.Create(
		ctx,
		&pb.Asset{
			Name:  args[0],
			Image: createAssetOpts.image,
		})
	if err != nil {
		fmt.Printf("Unable to create asset: %v\n", err)
		return
	}
	fmt.Printf("created asset with name '%v'\n", args[0])
}
