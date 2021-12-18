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

// CreateBlueprintOptions stores the necessary information to execute the create
// blueprint command.
type CreateBlueprintOptions struct {
	addr string

	name  string
	links []string
}

// NewCmdCreateBlueprint returns a new command that create a blueprint from
// command line arguments.
func NewCmdCreateBlueprint() *cobra.Command {
	o := &CreateBlueprintOptions{}

	cmd := &cobra.Command{
		Use:   "blueprint name",
		Short: "create a new blueprint",
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
func (o *CreateBlueprintOptions) addFlags(cmd *cobra.Command) {
	addAddrFlag(cmd, &o.addr, addrHelp)

	cmd.Flags().StringSliceVar(
		&o.links,
		"link",
		nil,
		"links to include in the blueprint")
}

// complete fills any remaining information necessary to run the command that is
// not specified by the user flags and is in the positional arguments
func (o *CreateBlueprintOptions) complete(args []string) error {
	if len(args) == 1 {
		o.name = args[0]
	}
	return nil
}

// validate verifies if the user options are valid and all necessary information
// for the command to run is present
func (o *CreateBlueprintOptions) validate() error {
	if o.name == "" {
		return errdefs.InvalidArgumentWithMsg("please specify a blueprint name")
	}
	if len(o.links) == 0 {
		return errdefs.InvalidArgumentWithMsg(
			"please specify at least one link",
		)
	}
	return nil
}

// run executes the create blueprint command with the specified options.
// It assumes the options were previously validated.
func (o *CreateBlueprintOptions) run() error {
	bp := &resources.BlueprintSpec{
		Name:  o.name,
		Links: o.links,
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
	return c.CreateBlueprint(ctx, bp)
}
