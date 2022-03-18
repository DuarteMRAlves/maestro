package start

import (
	"context"
	util2 "github.com/DuarteMRAlves/maestro/_old/internal/cli/maestroctl/cmd/util"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"io"
	"time"
)

// Options store the flags defined by the user when executing the attach
// command and then executes the command.
type Options struct {
	// address for the maestro server.
	maestro string

	// names of the orchestrations to be executed.
	names []string

	// Output for the cobra.Command to be executed.
	outWriter io.Writer
}

func NewCmdStart() *cobra.Command {
	o := &Options{}

	cmd := &cobra.Command{
		Use:                   "start [ORCHESTRATION...] [FLAGS]",
		DisableFlagsInUseLine: true,
		Short:                 "Start executions for orchestrations",
		Long: `Start executions for the received orchestrations.

If no orchestration is provided the default one is started.`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			err = o.complete(cmd, args)
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
			}
		},
	}

	o.addFlags(cmd)

	return cmd
}

// addFlags adds the necessary flags to the cobra.Command instance that will
// parse the command line arguments and run the command
func (o *Options) addFlags(cmd *cobra.Command) {
	util2.AddMaestroFlag(cmd, &o.maestro)
}

// complete fills any remaining information that is required to execute the
// create command.
func (o *Options) complete(cmd *cobra.Command, args []string) error {
	o.names = args
	o.outWriter = cmd.OutOrStdout()
	return nil
}

// validate verifies if the user inputs are valid and there are no conflicts
func (o *Options) validate() error {
	return nil
}

// run executes the Start command
func (o *Options) run() error {
	conn, err := grpc.Dial(o.maestro, grpc.WithInsecure())
	if err != nil {
		return errdefs.UnavailableWithMsg("create connection: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second,
	)
	defer cancel()

	stub := pb.NewExecutionManagementClient(conn)
	if len(o.names) == 0 {
		// Append empty string for the default orchestration.
		o.names = append(o.names, "")
	}
	for _, name := range o.names {
		startReq := &pb.StartExecutionRequest{Orchestration: name}
		if _, err = stub.Start(ctx, startReq); err != nil {
			return util2.ErrorFromGrpcError(err)
		}
	}
	return nil
}
