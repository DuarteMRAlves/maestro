package get

import (
	"context"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/cli/maestroctl/cmd/util"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"io"
	"sort"
	"time"
)

type StageOpts struct {
	// address for the maestro server
	maestro string

	name    string
	phase   string
	asset   string
	service string
	rpc     string
	address string

	// Output for the cobra.Command to be executed in order to verify outputs.
	outWriter io.Writer
}

func NewCmdGetStage() *cobra.Command {
	o := &StageOpts{}

	cmd := &cobra.Command{
		Use:                   "stage [NAME] [FLAGS]",
		DisableFlagsInUseLine: true,
		Short:                 "Display stages",
		Long: `Display relevant information related to stages.

The displayed stages can be filtered by specifying flags. When a flag is specified, 
only the stages with the flag value are displayed.

If a name is provided, only that stage is displayed.`,
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

	cmd.Flags().StringVar(&o.phase, "phase", "", "phase to search")
	cmd.Flags().StringVar(&o.asset, "asset", "", "asset name to search")
	cmd.Flags().StringVar(&o.service, "service", "", "service name to search")
	cmd.Flags().StringVar(&o.rpc, "rpc", "", "rpc name to search")
	cmd.Flags().StringVar(&o.address, "address", "", "address to search")
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
	if o.phase != "" {
		phase := api.StagePhase(o.phase)
		switch phase {
		case
			api.StagePending,
			api.StageRunning,
			api.StageFailed,
			api.StageSucceeded:
			// Do nothing
		default:
			return errdefs.InvalidArgumentWithMsg("unknown phase: %v", phase)
		}
	}
	return nil
}

// run executes the get link command
func (o *StageOpts) run() error {
	req := &pb.GetStageRequest{
		Name:    o.name,
		Phase:   o.phase,
		Asset:   o.asset,
		Service: o.service,
		Rpc:     o.rpc,
		Address: o.address,
	}

	conn, err := grpc.Dial(o.maestro, grpc.WithInsecure())
	if err != nil {
		return errdefs.UnavailableWithMsg("create connection: %v", err)
	}
	defer conn.Close()

	stub := pb.NewArchitectureManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := stub.GetStage(ctx, req)
	if err != nil {
		return util.ErrorFromGrpcError(err)
	}
	stages := make([]*pb.Stage, 0)
	for {
		s, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return util.ErrorFromGrpcError(err)
		}
		stages = append(stages, s)
	}

	return o.displayStages(stages)
}

func (o *StageOpts) displayStages(stages []*pb.Stage) error {
	sort.Slice(
		stages,
		func(i, j int) bool {
			return stages[i].Name < stages[j].Name
		},
	)
	numStages := len(stages)
	// Add space for all assets plus the header
	data := make([][]string, 0, numStages+1)

	head := []string{
		NameText,
		PhaseText,
		AssetText,
		ServiceText,
		RpcText,
		AddressText,
	}
	data = append(data, head)
	for _, s := range stages {
		stageData := []string{
			s.Name,
			s.Phase,
			s.Asset,
			s.Service,
			s.Rpc,
			s.Address,
		}
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
