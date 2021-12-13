package get

import (
	"context"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/cli/client"
	"github.com/DuarteMRAlves/maestro/internal/cli/util"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"io"
	"sort"
	"time"
)

const (
	assetFlag  = "asset"
	assetShort = "a"
	assetUsage = "asset name to search"

	serviceFlag  = "service"
	serviceShort = "s"
	serviceUsage = "service name to search"

	methodFlag  = "method"
	methodShort = "m"
	methodUsage = "method name to search"
)

type GetStageOpts struct {
	addr string

	name    string
	asset   string
	service string
	method  string

	// Output for the cobra.Command to be executed in order to verify outputs.
	outWriter io.Writer
}

func NewCmdGetStage() *cobra.Command {
	o := &GetStageOpts{}

	cmd := &cobra.Command{
		Use:     "stage",
		Short:   "list one or more stages",
		Args:    cobra.MaximumNArgs(1),
		Aliases: []string{"stages"},
		Run: func(cmd *cobra.Command, args []string) {
			err := o.complete(cmd, args)
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
func (o *GetStageOpts) addFlags(cmd *cobra.Command) {
	addAddrFlag(cmd, &o.addr)

	cmd.Flags().StringVarP(&o.asset, assetFlag, assetShort, "", assetUsage)
	cmd.Flags().StringVarP(
		&o.service,
		serviceFlag,
		serviceShort,
		"",
		serviceUsage)
	cmd.Flags().StringVarP(&o.method, methodFlag, methodShort, "", methodUsage)
}

// complete fills any remaining information for the runner that is not specified
// by the flags.
func (o *GetStageOpts) complete(cmd *cobra.Command, args []string) error {
	if len(args) == 1 {
		o.name = args[0]
	}
	o.outWriter = cmd.OutOrStdout()
	return nil
}

// validate checks if the user options are compatible and the command can
// be executed
func (o *GetStageOpts) validate() error {
	return nil
}

// run executes the get link command
func (o *GetStageOpts) run() error {
	query := &pb.Stage{
		Name:    o.name,
		Asset:   o.asset,
		Service: o.service,
		Method:  o.method,
	}

	conn := client.NewConnection(o.addr)
	defer conn.Close()

	c := pb.NewStageManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := c.Get(ctx, query)
	if err != nil {
		return client.ErrorFromGrpcError(err)
	}
	stages := make([]*pb.Stage, 0)
	for {
		s, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return client.ErrorFromGrpcError(err)
		}
		stages = append(stages, s)
	}
	return o.displayStages(stages)
}

func (o *GetStageOpts) displayStages(stages []*pb.Stage) error {
	sort.Slice(
		stages,
		func(i, j int) bool {
			return stages[i].Name < stages[j].Name
		})
	numStages := len(stages)
	// Add space for all assets plus the header
	data := make([][]string, 0, numStages+1)

	head := []string{NameText, AssetText, ServiceText, MethodText, AddressText}
	data = append(data, head)
	for _, s := range stages {
		stageData := []string{s.Name, s.Asset, s.Service, s.Method, s.Address}
		data = append(data, stageData)
	}

	output, err := pterm.DefaultTable.WithHasHeader().WithData(data).Srender()
	if err != nil {
		return errdefs.UnknownWithMsg("display stages: %v", err)
	}
	_, err = fmt.Fprintln(o.outWriter, output)
	if err != nil {
		return errdefs.UnknownWithMsg("display stages: %v", err)
	}
	return nil
}
