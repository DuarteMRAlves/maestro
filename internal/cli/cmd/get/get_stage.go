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
	"google.golang.org/grpc"
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

	addressFlag  = "address"
	addressUsage = "address to search"
)

type StageOpts struct {
	// address for the maestro server
	maestro string

	name    string
	asset   string
	service string
	method  string
	address string

	// Output for the cobra.Command to be executed in order to verify outputs.
	outWriter io.Writer
}

func NewCmdGetStage() *cobra.Command {
	o := &StageOpts{}

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
func (o *StageOpts) addFlags(cmd *cobra.Command) {
	util.AddMaestroFlag(cmd, &o.maestro)

	cmd.Flags().StringVarP(&o.asset, assetFlag, assetShort, "", assetUsage)
	cmd.Flags().StringVarP(
		&o.service,
		serviceFlag,
		serviceShort,
		"",
		serviceUsage)
	cmd.Flags().StringVarP(&o.method, methodFlag, methodShort, "", methodUsage)
	cmd.Flags().StringVar(&o.address, addressFlag, "", addressUsage)
}

// complete fills any remaining information for the runner that is not specified
// by the flags.
func (o *StageOpts) complete(cmd *cobra.Command, args []string) error {
	if len(args) == 1 {
		o.name = args[0]
	}
	o.outWriter = cmd.OutOrStdout()
	return nil
}

// validate checks if the user options are compatible and the command can
// be executed
func (o *StageOpts) validate() error {
	return nil
}

// run executes the get link command
func (o *StageOpts) run() error {
	query := &pb.Stage{
		Name:    o.name,
		Asset:   o.asset,
		Service: o.service,
		Method:  o.method,
		Address: o.address,
	}

	conn, err := grpc.Dial(o.maestro, grpc.WithInsecure())
	if err != nil {
		return errdefs.UnavailableWithMsg("create connection: %v", err)
	}
	defer conn.Close()

	c := client.New(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stages, err := c.GetStage(ctx, query)
	if err != nil {
		return err
	}
	return o.displayStages(stages)
}

func (o *StageOpts) displayStages(stages []*pb.Stage) error {
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
