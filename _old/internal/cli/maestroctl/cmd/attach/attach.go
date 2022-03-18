package attach

import (
	"context"
	"fmt"
	util2 "github.com/DuarteMRAlves/maestro/_old/internal/cli/maestroctl/cmd/util"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"io"
)

// Options store the flags defined by the user when executing the attach
// command and then executes the command.
type Options struct {
	// address for the maestro server
	maestro string

	// name of the orchestration for the execution to attach to
	name string

	// Output for the cobra.Command to be executed.
	outWriter io.Writer
}

func NewCmdAttach() *cobra.Command {
	o := &Options{}

	cmd := &cobra.Command{
		Use:                   "attach [EXECUTION] [FLAGS]",
		DisableFlagsInUseLine: true,
		Short:                 "Attach to a running execution",
		Long: `Attach to the output of a running execution.

If no execution is specified, the default one is used.`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			err = o.complete(cmd, args)
			if err != nil {
				util2.WriteOut(cmd, util2.DisplayMsgFromError(err))
				return
			}
			err = o.validate(args)
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
	if len(args) == 1 {
		o.name = args[0]
	}
	o.outWriter = cmd.OutOrStdout()
	return nil
}

// validate verifies if the user inputs are valid and there are no conflicts
func (o *Options) validate(args []string) error {
	// We only attach to one orchestration
	if len(args) != 1 {
		return errdefs.InvalidArgumentWithMsg("please specify one orchestration")
	}
	return nil
}

// run executes the Create command
func (o *Options) run() error {
	conn, err := grpc.Dial(o.maestro, grpc.WithInsecure())
	if err != nil {
		return errdefs.UnavailableWithMsg("create connection: %v", err)
	}
	defer conn.Close()

	stub := pb.NewExecutionManagementClient(conn)
	stream, err := stub.Attach(context.Background())
	if err != nil {
		return util2.ErrorFromGrpcError(err)
	}
	err = stream.Send(&pb.AttachExecutionRequest{Orchestration: o.name})
	if err != nil {
		return util2.ErrorFromGrpcError(err)
	}
	for {
		event, err := stream.Recv()
		if err != nil {
			return util2.ErrorFromGrpcError(err)
		}
		_, err = fmt.Fprintf(
			o.outWriter,
			"%v: %s",
			event.Timestamp,
			event.Description,
		)
		if err != nil {
			return errdefs.UnknownWithMsg("attach: %v", err)
		}
	}
}
