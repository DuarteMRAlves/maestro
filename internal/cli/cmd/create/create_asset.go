package create

import (
	"context"
	"github.com/DuarteMRAlves/maestro/internal/cli/client"
	"github.com/DuarteMRAlves/maestro/internal/cli/resources"
	"github.com/DuarteMRAlves/maestro/internal/cli/util"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"time"
)

const (
	imageFlag  = "image"
	imageUsage = "Docker image for the asset"
)

// AssetOpts executes a create asset command
type AssetOpts struct {
	// address for the maestro server
	maestro string

	name  string
	image string
}

// NewCmdCreateAsset returns a new command that creates an asset from command
// line arguments
func NewCmdCreateAsset() *cobra.Command {
	o := &AssetOpts{}

	cmd := &cobra.Command{
		Use:   "asset NAME",
		Short: "create a new Asset",
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
func (o *AssetOpts) addFlags(cmd *cobra.Command) {
	util.AddMaestroFlag(cmd, &o.maestro)

	cmd.Flags().StringVar(&o.image, imageFlag, "", imageUsage)
}

// complete fills any remaining information for the runner that is not specified
// by the flags.
func (o *AssetOpts) complete(args []string) error {
	if len(args) == 1 {
		o.name = args[0]
	}
	return nil
}

// validate checks if the user options are compatible and the command can
// be executed
func (o *AssetOpts) validate() error {
	if o.name == "" {
		return errdefs.InvalidArgumentWithMsg("please specify the asset name")
	}
	return nil
}

func (o *AssetOpts) run() error {
	asset := &resources.AssetSpec{
		Name:  o.name,
		Image: o.image,
	}

	conn, err := grpc.Dial(o.maestro, grpc.WithInsecure())
	if err != nil {
		return errdefs.UnavailableWithMsg("create connection: %v", err)
	}
	defer conn.Close()

	c := client.New(conn)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second)
	defer cancel()

	return c.CreateAsset(ctx, asset)
}
