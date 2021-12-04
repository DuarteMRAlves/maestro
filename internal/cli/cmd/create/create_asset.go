package create

import (
	"context"
	"github.com/DuarteMRAlves/maestro/internal/cli/client"
	"github.com/DuarteMRAlves/maestro/internal/cli/resources"
	"github.com/DuarteMRAlves/maestro/internal/cli/util"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/spf13/cobra"
	"time"
)

const (
	imageFlag  = "image"
	imageUsage = "Docker image for the asset"
)

// CreateAssetOptions executes a create asset command
type CreateAssetOptions struct {
	addr string

	name  string
	image string
}

// NewCmdCreateAsset returns a new command that creates an asset from command
// line arguments
func NewCmdCreateAsset() *cobra.Command {
	o := &CreateAssetOptions{}

	cmd := &cobra.Command{
		Use:   "asset NAME",
		Short: "create a new Asset",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := o.complete(args)
			if err != nil {
				util.WriteOut(cmd, util.DisplayMsgFromError(err))
				return
			}
			err = o.validate()
			if err != nil {
				util.WriteOut(cmd, util.DisplayMsgFromError(err))
				return
			}
			err = o.run()
			if err != nil {
				util.WriteOut(cmd, util.DisplayMsgFromError(err))
				return
			}
		},
	}

	o.addFlags(cmd)

	return cmd
}

// addFlags adds the necessary flags to the cobra.Command instance that will
// execute
func (o *CreateAssetOptions) addFlags(cmd *cobra.Command) {
	addAddrFlag(cmd, &o.addr, addrHelp)

	cmd.Flags().StringVar(&o.image, imageFlag, "", imageUsage)
}

// complete fills any remaining information for the runner that is not specified
// by the flags.
func (o *CreateAssetOptions) complete(args []string) error {
	if len(args) == 1 {
		o.name = args[0]
	}
	return nil
}

// validate checks if the user options are compatible and the command can
// be executed
func (o *CreateAssetOptions) validate() error {
	if o.name == "" {
		return errdefs.InvalidArgumentWithMsg("please specify the asset name")
	}
	return nil
}

func (o *CreateAssetOptions) run() error {
	asset := &resources.AssetResource{
		Name:  o.name,
		Image: o.image,
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second)
	defer cancel()

	return client.CreateAsset(ctx, asset, o.addr)
}
