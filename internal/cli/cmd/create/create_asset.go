package create

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/cli/client"
	"github.com/DuarteMRAlves/maestro/internal/cli/resources"
	"github.com/DuarteMRAlves/maestro/internal/cli/util"
	"github.com/spf13/cobra"
)

const (
	imageFlag  = "image"
	imageUsage = "Docker image for the asset"
)

var createAssetOpts = &struct {
	addr  string
	files []string
	image string
}{}

func NewCmdCreateAsset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "asset NAME",
		Short: "create a new Asset",
		Args:  cobra.MaximumNArgs(1),
		Run:   runCreateAsset,
	}

	addAddrFlag(cmd, &createAssetOpts.addr, addrHelp)
	addFilesFlag(cmd, &createAssetOpts.files, fileHelp)

	cmd.Flags().StringVar(&createAssetOpts.image, imageFlag, "", imageUsage)

	return cmd
}

func runCreateAsset(cmd *cobra.Command, args []string) {

	// Create from files
	if cmd.Flag(fileFull).Changed {
		util.WarnArgsIgnore(args, "creating from file")
		util.WarnFlagsIgnore(
			cmd,
			[]string{imageFlag},
			"creating from file")
		err := createFromFiles(
			createAssetOpts.files,
			createAssetOpts.addr,
			resources.AssetKind)
		if err != nil {
			fmt.Printf("unable to create resources: %v\n", err)
		}
	} else {
		if len(args) != 1 {
			fmt.Printf("unable to create asset: expected name argument")
			return
		}
		asset := &resources.AssetResource{
			Name:  args[0],
			Image: createAssetOpts.image,
		}
		err := client.CreateAsset(asset, createAssetOpts.addr)

		if err != nil {
			fmt.Printf("Unable to create asset: %v\n", err)
			return
		}
		fmt.Printf("created asset with name '%v'\n", args[0])
	}
}
