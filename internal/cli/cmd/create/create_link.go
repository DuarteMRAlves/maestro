package create

import (
	"context"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/cli/client"
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
	cmd.Flags().StringVar(
		&createLinkOpts.sourceStage,
		sourceStageFlag,
		"",
		sourceStageUsage)
	cmd.MarkFlagRequired(sourceStageFlag)

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
	cmd.MarkFlagRequired(targetStageFlag)

	cmd.Flags().StringVar(
		&createLinkOpts.targetField,
		targetFieldFlag,
		"",
		targetFieldUsage)

	return cmd
}

func runCreateLink(cmd *cobra.Command, args []string) {
	conn := client.NewConnection(createOpts.addr)
	defer conn.Close()

	c := pb.NewLinkManagementClient(conn)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second)
	defer cancel()

	_, err := c.Create(
		ctx,
		&pb.Link{
			Name:        args[0],
			SourceStage: createLinkOpts.sourceStage,
			SourceField: createLinkOpts.sourceField,
			TargetStage: createLinkOpts.targetStage,
			TargetField: createLinkOpts.targetField,
		})
	if err != nil {
		fmt.Printf("unable to create link: %v\n", err)
		return
	}
	fmt.Printf("link '%v' created\n", args[0])

}
