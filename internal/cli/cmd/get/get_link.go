package get

import (
	"context"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/cli"
	"github.com/DuarteMRAlves/maestro/internal/cli/display/table"
	"github.com/spf13/cobra"
	"io"
	"log"
	"time"
)

const (
	sourceStageFlag  = "source-stage"
	sourceStageUsage = "name of the source stage to search"

	sourceFieldFlag  = "source-field"
	sourceFieldUsage = "field in the source message search"

	targetStageFlag  = "target-stage"
	targetStageUsage = "name of the target stage to search"

	targetFieldFlag  = "target-field"
	targetFieldUsage = "field in the target message search"
)

var getLinkOpts = &struct {
	sourceStage string
	sourceField string
	targetStage string
	targetField string
}{}

func NewCmdGetLink() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "link",
		Short:   "list one or more links",
		Args:    cobra.MaximumNArgs(1),
		Aliases: []string{"stages"},
		Run:     runGetLink,
	}

	cmd.Flags().StringVar(
		&getLinkOpts.sourceStage,
		sourceStageFlag,
		"",
		sourceStageUsage)

	cmd.Flags().StringVar(
		&getLinkOpts.sourceField,
		sourceFieldFlag,
		"",
		sourceFieldUsage)

	cmd.Flags().StringVar(
		&getLinkOpts.targetStage,
		targetStageFlag,
		"",
		targetStageUsage)

	cmd.Flags().StringVar(
		&getLinkOpts.targetField,
		targetFieldFlag,
		"",
		targetFieldUsage)

	return cmd
}

func runGetLink(_ *cobra.Command, args []string) {
	query := &pb.Link{
		SourceStage: getLinkOpts.sourceStage,
		SourceField: getLinkOpts.sourceField,
		TargetStage: getLinkOpts.targetStage,
		TargetField: getLinkOpts.targetField,
	}

	if len(args) == 1 {
		query.Name = args[0]
	}

	conn := cli.NewConnection(createOpts.addr)
	defer conn.Close()

	c := pb.NewLinkManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := c.Get(ctx, query)
	if err != nil {
		log.Fatalf("list links: %v", err)
	}
	links := make([]*pb.Link, 0)
	for {
		l, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("list links: %v", err)
		}
		links = append(links, l)
	}
	if err := displayLinks(links); err != nil {
		log.Fatalf("list links: %v", err)
	}
}

func displayLinks(links []*pb.Link) error {
	numLinks := len(links)
	names := make([]string, 0, numLinks)
	sourceStages := make([]string, 0, numLinks)
	sourceFields := make([]string, 0, numLinks)
	targetStages := make([]string, 0, numLinks)
	targetFields := make([]string, 0, numLinks)

	for _, l := range links {
		names = append(names, l.Name)
		sourceStages = append(sourceStages, l.SourceStage)
		sourceFields = append(sourceFields, l.SourceField)
		targetStages = append(targetStages, l.TargetStage)
		targetFields = append(targetFields, l.TargetField)
	}

	t := table.NewBuilder().
		WithPadding(colPad).
		WithMinColSize(minColSize).
		Build()

	if err := t.AddColumn(NameText, names); err != nil {
		return err
	}
	if err := t.AddColumn(SourceStageText, sourceStages); err != nil {
		return err
	}
	if err := t.AddColumn(SourceFieldText, sourceFields); err != nil {
		return err
	}
	if err := t.AddColumn(TargetStageText, targetStages); err != nil {
		return err
	}
	if err := t.AddColumn(TargetFieldText, targetFields); err != nil {
		return err
	}
	fmt.Print(t)
	return nil
}
