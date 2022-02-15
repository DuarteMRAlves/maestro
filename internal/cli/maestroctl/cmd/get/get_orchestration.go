package get

import (
	"context"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/cli/maestroctl/client"
	util2 "github.com/DuarteMRAlves/maestro/internal/cli/maestroctl/cmd/util"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"io"
	"sort"
	"time"
)

// OrchestrationOpts stores the necessary information to execute a get
// orchestration command.
type OrchestrationOpts struct {
	// address for the maestro server
	maestro string

	name  string
	phase string

	// Output for the cobra.command
	outWriter io.Writer
}

func NewCmdGetOrchestration() *cobra.Command {
	o := &OrchestrationOpts{}

	cmd := &cobra.Command{
		Use:     "orchestration",
		Short:   "list one or more orchestrations",
		Args:    cobra.MaximumNArgs(1),
		Aliases: []string{"orchestrations"},
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
func (o *OrchestrationOpts) addFlags(cmd *cobra.Command) {
	util2.AddMaestroFlag(cmd, &o.maestro)

	cmd.Flags().StringVar(&o.phase, "phase", "", "phase to search")
}

// complete fills any remaining information for the runner that is not specified
// by the flags.
func (o *OrchestrationOpts) complete(
	cmd *cobra.Command,
	args []string,
) error {
	if len(args) == 1 {
		o.name = args[0]
	}
	o.outWriter = cmd.OutOrStdout()
	return nil
}

// validate checks if the user options are compatible and the command can
// be executed
func (o *OrchestrationOpts) validate() error {
	if o.phase != "" {
		phase := api.OrchestrationPhase(o.phase)
		switch phase {
		case
			api.OrchestrationPending,
			api.OrchestrationRunning,
			api.OrchestrationSucceeded,
			api.OrchestrationFailed:
			// Do nothing
		default:
			return errdefs.InvalidArgumentWithMsg("unknown phase: %v", phase)
		}
	}
	return nil
}

// run executes the get link command
func (o *OrchestrationOpts) run() error {
	req := &pb.GetOrchestrationRequest{
		Name:  o.name,
		Phase: o.phase,
	}

	conn, err := grpc.Dial(o.maestro, grpc.WithInsecure())
	if err != nil {
		return errdefs.UnavailableWithMsg("create connection: %v", err)
	}
	defer conn.Close()

	c := client.New(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	orchestrations, err := c.GetOrchestration(ctx, req)
	if err != nil {
		return err
	}
	return o.displayOrchestrations(orchestrations)
}

func (o *OrchestrationOpts) displayOrchestrations(
	orchestrations []*pb.Orchestration,
) error {
	sort.Slice(
		orchestrations,
		func(i, j int) bool {
			return orchestrations[i].Name < orchestrations[j].Name
		},
	)
	numOrchestrations := len(orchestrations)
	// Add space for all assets plus the header
	data := make([][]string, 0, numOrchestrations+1)

	headers := []string{NameText, PhaseText}
	data = append(data, headers)

	for _, orchestration := range orchestrations {
		orchestrationData := []string{orchestration.Name, orchestration.Phase}
		data = append(data, orchestrationData)
	}

	output, err := pterm.DefaultTable.WithHasHeader().WithData(data).Srender()
	if err != nil {
		return errdefs.UnknownWithMsg("display orchestrations: %v", err)
	}
	_, err = fmt.Fprintln(o.outWriter, output)
	if err != nil {
		return errdefs.UnknownWithMsg("display orchestrations: %v", err)
	}
	return nil
}
