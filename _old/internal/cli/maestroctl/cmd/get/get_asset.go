package get

import (
	"context"
	"fmt"
	util2 "github.com/DuarteMRAlves/maestro/_old/internal/cli/maestroctl/cmd/util"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"io"
	"sort"
	"time"
)

// AssetOpts executes a get asset command
type AssetOpts struct {
	// address for the maestro server
	maestro string

	name  string
	image string

	// Output for the cobra.Command to be executed in order to verify outputs.
	outWriter io.Writer
}

func NewCmdGetAsset() *cobra.Command {
	o := &AssetOpts{}

	cmd := &cobra.Command{
		Use:                   "asset [NAME] [FLAGS]",
		DisableFlagsInUseLine: true,
		Short:                 "Display assets",
		Long: `Display relevant information related to assets.

The displayed assets can be filtered by specifying flags. When a flag is specified, 
only the assets with the flag value are displayed.

If a name is provided, only that asset is displayed.`,
		Args:    cobra.MaximumNArgs(1),
		Aliases: []string{"assets"},
		Run: func(cmd *cobra.Command, args []string) {
			err := o.complete(cmd, args)
			if err != nil {
				util2.WriteOut(cmd, util2.DisplayMsgFromError(err))
				return
			}
			err = o.validate()
			if err != nil {
				util2.WriteOut(cmd, util2.DisplayMsgFromError(err))
				return
			}
			err = o.run()
			if err != nil {
				util2.WriteOut(cmd, util2.DisplayMsgFromError(err))
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
	util2.AddMaestroFlag(cmd, &o.maestro)

	cmd.Flags().StringVar(&o.image, "image", "", "image name to search")
}

// complete fills any remaining information for the runner that is not specified
// by the flags.
func (o *AssetOpts) complete(cmd *cobra.Command, args []string) error {
	if len(args) == 1 {
		o.name = args[0]
	}
	o.outWriter = cmd.OutOrStdout()
	return nil
}

// validate checks if the user options are compatible and the command can
// be executed
func (o *AssetOpts) validate() error {
	return nil
}

// run executes the get asset command
func (o *AssetOpts) run() error {
	req := &pb.GetAssetRequest{
		Name:  o.name,
		Image: o.image,
	}

	conn, err := grpc.Dial(o.maestro, grpc.WithInsecure())
	if err != nil {
		return errdefs.UnavailableWithMsg("create connection: %v", err)
	}
	defer conn.Close()

	stub := pb.NewArchitectureManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := stub.GetAsset(ctx, req)
	if err != nil {
		return util2.ErrorFromGrpcError(err)
	}
	assets := make([]*pb.Asset, 0)
	for {
		a, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return util2.ErrorFromGrpcError(err)
		}
		assets = append(assets, a)
	}
	return o.displayAssets(assets)
}

func (o *AssetOpts) displayAssets(assets []*pb.Asset) error {
	sort.Slice(
		assets,
		func(i, j int) bool {
			return assets[i].Name < assets[j].Name
		},
	)
	numAssets := len(assets)
	// Add space for all assets plus the header
	data := make([][]string, 0, numAssets+1)
	data = append(data, []string{NameText, ImageText})
	for _, a := range assets {
		data = append(data, []string{a.Name, a.Image})
	}
	output, err := pterm.DefaultTable.WithHasHeader().WithData(data).Srender()
	if err != nil {
		return errdefs.UnknownWithMsg("display assets: %v", err)
	}
	_, err = fmt.Fprintln(o.outWriter, output)
	if err != nil {
		return errdefs.UnknownWithMsg("display assets: %v", err)
	}
	return nil
}