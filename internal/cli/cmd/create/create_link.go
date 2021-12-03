package create

import (
	"context"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/cli/client"
	"github.com/DuarteMRAlves/maestro/internal/cli/resources"
	"github.com/DuarteMRAlves/maestro/internal/cli/util"
	"github.com/spf13/cobra"
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

var createLinkOpts = &struct {
	addr        string
	files       []string
	sourceStage string
	sourceField string
	targetStage string
	targetField string
}{}

func NewCmdCreateLink() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "link name",
		Short: "create a new link",
		Args:  cobra.ExactArgs(1),
		Run:   runCreateLink,
	}

	addAddrFlag(cmd, &createLinkOpts.addr, addrHelp)
	addFilesFlag(cmd, &createLinkOpts.files, addrHelp)

	cmd.Flags().StringVar(
		&createLinkOpts.sourceStage,
		sourceStageFlag,
		"",
		sourceStageUsage)

	cmd.Flags().StringVar(
		&createLinkOpts.sourceField,
		sourceFieldFlag,
		"",
		sourceFieldUsage)

	cmd.Flags().StringVar(
		&createLinkOpts.targetStage,
		targetStageFlag,
		"",
		targetStageUsage)

	cmd.Flags().StringVar(
		&createLinkOpts.targetField,
		targetFieldFlag,
		"",
		targetFieldUsage)

	return cmd
}

func runCreateLink(cmd *cobra.Command, args []string) {
	if cmd.Flag(fileFull).Changed {
		util.WarnArgsIgnore(args, "creating from file")
		util.WarnFlagsIgnore(
			cmd,
			[]string{
				sourceStageFlag,
				sourceFieldFlag,
				targetStageFlag,
				targetFieldFlag,
			},
			"creating from file")
		err := createFromFiles(
			createLinkOpts.files,
			createStageOpts.addr,
			resources.LinkKind)
		if err != nil {
			fmt.Printf("unable to create link: %v", err)
		}
	} else {
		if len(args) != 1 {
			fmt.Printf("unable to create link: expected name argument")
			return
		}
		err := util.VerifyFlagsChanged(
			cmd,
			[]string{sourceStageFlag, targetStageFlag})
		if err != nil {
			fmt.Printf("unable to create link: %v", err)
		}
		link := &resources.LinkResource{
			Name:        args[0],
			SourceStage: createLinkOpts.sourceStage,
			SourceField: createLinkOpts.sourceField,
			TargetStage: createLinkOpts.targetStage,
			TargetField: createLinkOpts.targetField,
		}

		ctx, cancel := context.WithTimeout(
			context.Background(),
			time.Second)
		defer cancel()
		err = client.CreateLink(ctx, link, createLinkOpts.addr)
		if err != nil {
			fmt.Printf("unable to create link: %v\n", err)
			return
		}
		fmt.Printf("link '%v' created\n", args[0])
	}

}
