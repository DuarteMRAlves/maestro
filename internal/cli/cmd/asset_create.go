package cmd

import (
	"context"
	"fmt"
	pb "github.com/DuarteMRAlves/maestro/api/pb"
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

var createAssetArgs = &struct {
	name  string
	image string
}{}

var createAssetCmd = &cobra.Command{
	Use:   "create",
	Short: "maestro-cli asset create allows you to create a new Asset",
	Run: func(cmd *cobra.Command, args []string) {
		conn := cli.NewConnection(addr)
		defer conn.Close()

		c := pb.NewAssetManagementClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		asset, err := c.Create(
			ctx,
			&pb.Asset{Name: createAssetArgs.name, Image: createAssetArgs.image})

		if err != nil {
			log.Fatalf("Unable to create asset: %v", err)
		}
		fmt.Printf("created asset with identifier '%v'\n", asset.Val)
	},
}

func init() {
	assetCmd.AddCommand(createAssetCmd)

	createAssetCmd.Flags().StringVar(
		&createAssetArgs.name,
		nameFlag,
		"",
		nameUsage)
	createAssetCmd.MarkFlagRequired(nameFlag)

	createAssetCmd.Flags().StringVar(
		&createAssetArgs.image,
		imageFlag,
		"",
		imageUsage)
}
