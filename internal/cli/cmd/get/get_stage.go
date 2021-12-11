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
}

func NewCmdGetStage() *cobra.Command {
	o := &GetStageOpts{}

	cmd := &cobra.Command{
		Use:     "stage",
		Short:   "list one or more stages",
		Args:    cobra.MaximumNArgs(1),
		Aliases: []string{"stages"},
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
func (o *GetStageOpts) complete(args []string) error {
	if len(args) == 1 {
		o.name = args[0]
	}
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
		return err
	}
	stages := make([]*pb.Stage, 0)
	for {
		s, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		stages = append(stages, s)
	}
	return displayStages(stages)
}

func displayStages(stages []*pb.Stage) error {
	numStages := len(stages)
	names := make([]string, 0, numStages)
	assets := make([]string, 0, numStages)
	services := make([]string, 0, numStages)
	methods := make([]string, 0, numStages)

	for _, s := range stages {
		names = append(names, s.Name)
		assets = append(assets, s.Asset)
		services = append(services, s.Service)
		methods = append(methods, s.Method)
	}

	t := table.NewBuilder().
		WithPadding(colPad).
		WithMinColSize(minColSize).
		Build()

	if err := t.AddColumn(NameText, names); err != nil {
		return err
	}
	if err := t.AddColumn(AssetText, assets); err != nil {
		return err
	}
	if err := t.AddColumn(ServiceText, services); err != nil {
		return err
	}
	if err := t.AddColumn(MethodText, methods); err != nil {
		return err
	}
	fmt.Print(t)
	return nil
}
