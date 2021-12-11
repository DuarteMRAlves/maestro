package get

import (
	"context"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/cli/client"
	"github.com/DuarteMRAlves/maestro/internal/cli/display/table"
	"github.com/DuarteMRAlves/maestro/internal/cli/util"
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
		Use:   "asset",
		Short: "list one or more Assets",
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
	names := make([]string, 0, numAssets)
	images := make([]string, 0, numAssets)
	for _, a := range assets {
		names = append(names, a.Name)
		images = append(images, a.Image)
	}

	t := table.NewBuilder().
		WithPadding(colPad).
		WithMinColSize(minColSize).
		Build()

	if err := t.AddColumn(NameText, names); err != nil {
		return err
	}
	if err := t.AddColumn(ImageText, images); err != nil {
		return err
	}
	fmt.Print(t)
	return nil
}
