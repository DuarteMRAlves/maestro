package get

import (
	"context"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/cli/client"
	"github.com/DuarteMRAlves/maestro/internal/cli/util"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"io"
	"time"
)

// GetAssetOptions executes a get asset command
type GetAssetOptions struct {
	addr string
}

func NewCmdGetAsset() *cobra.Command {
	o := &GetAssetOptions{}

	cmd := &cobra.Command{
		Use:     "asset",
		Short:   "list one or more Assets",
		Aliases: []string{"assets"},
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
func (o *GetAssetOptions) addFlags(cmd *cobra.Command) {
	addAddrFlag(cmd, &o.addr)
}

// complete fills any remaining information for the runner that is not specified
// by the flags.
func (o *GetAssetOptions) complete(_ []string) error {
	return nil
}

// validate checks if the user options are compatible and the command can
// be executed
func (o *GetAssetOptions) validate() error {
	return nil
}

// run executes the get asset command
func (o *GetAssetOptions) run() error {
	conn := client.NewConnection(o.addr)
	defer conn.Close()

	c := pb.NewAssetManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := c.Get(ctx, &pb.Asset{})
	assets := make([]*pb.Asset, 0)
	if err != nil {
		return err
	}
	for {
		a, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		assets = append(assets, a)
	}
	return displayAssets(assets)
}

func displayAssets(assets []*pb.Asset) error {
	numAssets := len(assets)
	// Add space for all assets plus the header
	data := make([][]string, 0, numAssets+1)
	data = append(data, []string{NameText, ImageText})
	for _, a := range assets {
		data = append(data, []string{a.Name, a.Image})
	}
	err := pterm.DefaultTable.WithHasHeader().WithData(data).Render()
	if err != nil {
		return errdefs.UnknownWithMsg("display assets: %v", err)
	}
	return nil
}
