package create

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/cli/client"
	"github.com/spf13/cobra"
)

const (
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
		Args:  cobra.MaximumNArgs(1),
		Run:   runCreateAsset,
	}

	cmd.Flags().StringVar(&createAssetOpts.image, imageFlag, "", imageUsage)

	return cmd
}

func runCreateAsset(cmd *cobra.Command, args []string) {
	var asset *pb.Asset
	// Create from files
	if cmd.Flag(fileFull).Changed {
		if len(args) > 0 {
			fmt.Printf(
				"warning: creating from files, positional arguments ignored\n")
		}
		if cmd.Flag(imageFlag).Changed {
			fmt.Printf(
				"warning: creating from files, %v flag ignored\n",
				imageFlag)
		}
		err := createFromFiles(createOpts.files, createOpts.addr, assetKind)
		if err != nil {
			fmt.Printf("unable to create resources: %v\n", err)
		}
	} else {
		asset = &pb.Asset{
			Name:  args[0],
			Image: createAssetOpts.image,
		}
		err := client.CreateAsset(asset, createOpts.addr)

		if err != nil {
			fmt.Printf("Unable to create asset: %v\n", err)
			return
		}
		fmt.Printf("created asset with name '%v'\n", args[0])
	}
}
