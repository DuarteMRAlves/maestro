package create

import (
	"context"
	"github.com/DuarteMRAlves/maestro/internal/cli/client"
	"github.com/DuarteMRAlves/maestro/internal/cli/resources"
	"github.com/DuarteMRAlves/maestro/internal/cli/util"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"time"
)

// LinkOpts stores the flags for the CreateLink command and executes it
type LinkOpts struct {
	// address for the maestro server
	maestro string

	name        string
	sourceStage string
	sourceField string
	targetStage string
	targetField string
}

// NewCmdCreateLink returns a new command that creates a link from command line
// arguments
func NewCmdCreateLink() *cobra.Command {
	o := &LinkOpts{}

	cmd := &cobra.Command{
		Use:   "link name",
		Short: "create a new link",
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
// run the CreateLink command
func (o *LinkOpts) addFlags(cmd *cobra.Command) {
	util.AddMaestroFlag(cmd, &o.maestro)

	cmd.Flags().StringVar(
		&o.sourceStage,
		"source-stage",
		"",
		"name of the source stage to link (required)")
	cmd.Flags().StringVar(&o.sourceField,
		"source-field",
		"",
		"field in the source message to use. If not specified the whole "+
			"message is used.")
	cmd.Flags().StringVar(
		&o.targetStage,
		"target-stage",
		"",
		"name of the target stage to link (required)")
	cmd.Flags().StringVar(
		&o.targetField,
		"target-field",
		"",
		"field in the target message to set. If not specified the entire "+
			"message is used.")
}

// complete fills any remaining information necessary to run the command that is
// not specified by the user flags and is in the positional arguments
func (o *LinkOpts) complete(args []string) error {
	if len(args) == 1 {
		o.name = args[0]
	}
	return nil
}

// validate verifies if the user options are valid and all necessary information
// for the command to run is present
func (o *LinkOpts) validate() error {
	if o.name == "" {
		return errdefs.InvalidArgumentWithMsg("please specify a link name")
	}
	if o.sourceStage == "" {
		return errdefs.InvalidArgumentWithMsg("please specify a source stage")
	}
	if o.targetStage == "" {
		return errdefs.InvalidArgumentWithMsg("please specify a target stage")
	}
	return nil
}

// run executes a CreateLink command with the specified options.
// It assumes the options were previously validated.
func (o *LinkOpts) run() error {
	link := &resources.LinkSpec{
		Name:        o.name,
		SourceStage: o.sourceStage,
		SourceField: o.sourceField,
		TargetStage: o.targetStage,
		TargetField: o.targetField,
	}

	conn, err := grpc.Dial(o.maestro, grpc.WithInsecure())
	if err != nil {
		return errdefs.UnavailableWithMsg("create connection: %v", err)
	}
	defer conn.Close()

	c := client.New(conn)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second)
	defer cancel()
	return c.CreateLink(ctx, link)
}
