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

// GetOrchestrationOptions stores the necessary information to execute a get
// orchestration command.
type GetOrchestrationOptions struct {
	addr string

	name string

	// Output for the cobra.command
	outWriter io.Writer
}

func NewCmdGetOrchestration() *cobra.Command {
	o := &GetOrchestrationOptions{}

	cmd := &cobra.Command{
		Use:     "orchestration",
		Short:   "list one or more orchestrations",
		Args:    cobra.MaximumNArgs(1),
		Aliases: []string{"orchestrations"},
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
func (o *GetOrchestrationOptions) addFlags(cmd *cobra.Command) {
	util.AddAddrFlag(cmd, &o.addr)
}

// complete fills any remaining information for the runner that is not specified
// by the flags.
func (o *GetOrchestrationOptions) complete(
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
func (o *GetOrchestrationOptions) validate() error {
	return nil
}

// run executes the get link command
func (o *GetOrchestrationOptions) run() error {
	query := &pb.Orchestration{
		Name: o.name,
	}

	conn, err := grpc.Dial(o.addr, grpc.WithInsecure())
	if err != nil {
		return errdefs.UnavailableWithMsg("create connection: %v", err)
	}
	defer conn.Close()

	c := client.New(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	orchestrations, err := c.GetOrchestration(ctx, query)
	if err != nil {
		return err
	}
	return o.displayOrchestrations(orchestrations)
}

func (o *GetOrchestrationOptions) displayOrchestrations(
	orchestrations []*pb.Orchestration,
) error {
	sort.Slice(
		orchestrations,
		func(i, j int) bool {
			return orchestrations[i].Name < orchestrations[j].Name
		})
	numOrchestrations := len(orchestrations)
	// Add space for all assets plus the header
	data := make([][]string, 0, numOrchestrations+1)

	headers := []string{NameText}
	data = append(data, headers)

	for _, orchestration := range orchestrations {
		orchestrationData := []string{orchestration.Name}
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
