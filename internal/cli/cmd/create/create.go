package create

import (
	"context"
	"github.com/DuarteMRAlves/maestro/internal/cli/client"
	"github.com/DuarteMRAlves/maestro/internal/cli/resources"
	"github.com/DuarteMRAlves/maestro/internal/cli/util"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/spf13/cobra"
	"time"
)

// Flags store the flags defined by the user when executing the create
// command.
type Flags struct {
	addr string

	files []string
}

func NewCmdCreate() *cobra.Command {
	flags := &Flags{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create resources of a given type",
		Args:  cobra.MaximumNArgs(0),
		Run: func(cmd *cobra.Command, _ []string) {
			err := validate(flags)
			if err != nil {
				util.WriteOut(cmd, util.DisplayMsgFromError(err))
				return
			}
			err = execute(flags)
			if err != nil {
				util.WriteOut(cmd, util.DisplayMsgFromError(err))
			}
		},
	}

	addFlags(cmd, flags)

	// Subcommands
	cmd.AddCommand(NewCmdCreateAsset())
	cmd.AddCommand(NewCmdCreateStage())
	cmd.AddCommand(NewCmdCreateLink())

	return cmd
}

func addFlags(cmd *cobra.Command, flags *Flags) {
	addAddrFlag(cmd, &flags.addr, addrHelp)
	addFilesFlag(cmd, &flags.files, fileHelp)
}

func validate(flags *Flags) error {
	// In create, we only accept files
	if len(flags.files) == 0 {
		return errdefs.InvalidArgumentWithMsg("please specify input files")
	}
	return nil
}

func execute(flags *Flags) error {
	parsed, err := resources.ParseFiles(flags.files)
	if err != nil {
		return err
	}
	if err = resources.IsValidKinds(parsed); err != nil {
		return err
	}

	assets := resources.FilterAssets(parsed)
	stages := resources.FilterStages(parsed)
	links := resources.FilterLinks(parsed)

	orderedResourcesSize := len(assets) + len(stages) + len(links)

	resourcesByKind := make([]*resources.Resource, 0, orderedResourcesSize)

	resourcesByKind = append(resourcesByKind, assets...)
	resourcesByKind = append(resourcesByKind, stages...)
	resourcesByKind = append(resourcesByKind, links...)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second)
	defer cancel()

	for _, r := range resourcesByKind {
		if err := client.CreateResource(ctx, r, flags.addr); err != nil {
			return err
		}
	}
	return nil
}
