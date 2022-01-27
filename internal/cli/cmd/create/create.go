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

// Options store the flags defined by the user when executing the create
// command and then executes the command.
type Options struct {
	// address for the maestro server
	maestro string

	files []string
}

func NewCmdCreate() *cobra.Command {
	o := &Options{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create resources of a given type",
		Args:  cobra.MaximumNArgs(0),
		Run: func(cmd *cobra.Command, _ []string) {
			err := o.validate()
			if err != nil {
				util.WriteOut(cmd, util.DisplayMsgFromError(err))
				return
			}
			err = o.run()
			if err != nil {
				util.WriteOut(cmd, util.DisplayMsgFromError(err))
			}
		},
	}

	o.addFlags(cmd)

	// Subcommands
	cmd.AddCommand(NewCmdCreateAsset())
	cmd.AddCommand(NewCmdCreateStage())
	cmd.AddCommand(NewCmdCreateLink())
	cmd.AddCommand(NewCmdCreateOrchestration())

	return cmd
}

// addFlags adds the necessary flags to the cobra.Command instance that will
// parse the command line arguments and run the command
func (o *Options) addFlags(cmd *cobra.Command) {
	util.AddMaestroFlag(cmd, &o.maestro)
	util.AddFilesFlag(cmd, &o.files, "files to create one or more resources")
}

// validate verifies if the user inputs are valid and there are no conflits
func (o *Options) validate() error {
	// In create, we only accept files
	if len(o.files) == 0 {
		return errdefs.InvalidArgumentWithMsg("please specify input files")
	}
	return nil
}

// run executes the Create command
func (o *Options) run() error {
	parsed, err := resources.ParseFiles(o.files)
	if err != nil {
		return err
	}
	if err = resources.IsValidKinds(parsed); err != nil {
		return err
	}

	assets := resources.FilterAssets(parsed)
	stages := resources.FilterStages(parsed)
	links := resources.FilterLinks(parsed)
	orchestrations := resources.FilterOrchestrations(parsed)

	orderedResourcesSize :=
		len(assets) + len(stages) + len(links) + len(orchestrations)

	resourcesByKind := make([]*resources.Resource, 0, orderedResourcesSize)

	resourcesByKind = append(resourcesByKind, assets...)
	resourcesByKind = append(resourcesByKind, orchestrations...)
	resourcesByKind = append(resourcesByKind, stages...)
	resourcesByKind = append(resourcesByKind, links...)

	conn, err := grpc.Dial(o.maestro, grpc.WithInsecure())
	if err != nil {
		return errdefs.UnavailableWithMsg("create connection: %v", err)
	}
	defer conn.Close()

	c := client.New(conn)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second,
	)
	defer cancel()

	for _, r := range resourcesByKind {
		if err := c.CreateResource(ctx, r); err != nil {
			return err
		}
	}
	return nil
}
