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

const (
	sourceStageFlag  = "source-stage"
	sourceStageUsage = "name of the source stage to link (required)"

	sourceFieldFlag  = "source-field"
	sourceFieldUsage = "field in the source message to use. " +
		"If not specified the whole message is used."

	targetStageFlag  = "target-stage"
	targetStageUsage = "name of the target stage to link (required)"

	targetFieldFlag  = "target-field"
	targetFieldUsage = "field in the target message to set. " +
		"If not specified the entire message is used."
)

// CreateLinkOptions stores the flags for the CreateLink command and executes it
type CreateLinkOptions struct {
	addr string

	name        string
	sourceStage string
	sourceField string
	targetStage string
	targetField string
}

// NewCmdCreateLink returns a new command that creates a link from command line
// arguments
func NewCmdCreateLink() *cobra.Command {
	o := &CreateLinkOptions{}

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
func (o *CreateLinkOptions) addFlags(cmd *cobra.Command) {
	util.AddAddrFlag(cmd, &o.addr)

	cmd.Flags().StringVar(&o.sourceStage, sourceStageFlag, "", sourceStageUsage)
	cmd.Flags().StringVar(&o.sourceField, sourceFieldFlag, "", sourceFieldUsage)
	cmd.Flags().StringVar(&o.targetStage, targetStageFlag, "", targetStageUsage)
	cmd.Flags().StringVar(&o.targetField, targetFieldFlag, "", targetFieldUsage)
}

// complete fills any remaining information necessary to run the command that is
// not specified by the user flags and is in the positional arguments
func (o *CreateLinkOptions) complete(args []string) error {
	if len(args) == 1 {
		o.name = args[0]
	}
	return nil
}

// validate verifies if the user options are valid and all necessary information
// for the command to run is present
func (o *CreateLinkOptions) validate() error {
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

// CreateLinkOptions runs a CreateLink command with the specified options.
// It assumes the options were previously validated.
func (o *CreateLinkOptions) run() error {
	link := &resources.LinkSpec{
		Name:        o.name,
		SourceStage: o.sourceStage,
		SourceField: o.sourceField,
		TargetStage: o.targetStage,
		TargetField: o.targetField,
	}

	conn, err := grpc.Dial(o.addr, grpc.WithInsecure())
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
