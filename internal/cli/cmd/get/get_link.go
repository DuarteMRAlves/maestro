package get

import (
	"context"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/cli/client"
	"github.com/DuarteMRAlves/maestro/internal/cli/util"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"io"
	"sort"
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

type GetLinkOptions struct {
	addr string

	name        string
	sourceStage string
	sourceField string
	targetStage string
	targetField string

	// Output for the cobra.Command to be executed in order to verify outputs.
	outWriter io.Writer
}

func NewCmdGetLink() *cobra.Command {
	o := &GetLinkOptions{}

	cmd := &cobra.Command{
		Use:     "link",
		Short:   "list one or more links",
		Args:    cobra.MaximumNArgs(1),
		Aliases: []string{"links"},
		Run: func(cmd *cobra.Command, args []string) {
			err := o.complete(cmd, args)
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
// execute
func (o *GetLinkOptions) addFlags(cmd *cobra.Command) {
	addAddrFlag(cmd, &o.addr)

	cmd.Flags().StringVar(&o.sourceStage, sourceStageFlag, "", sourceStageUsage)
	cmd.Flags().StringVar(&o.sourceField, sourceFieldFlag, "", sourceFieldUsage)
	cmd.Flags().StringVar(&o.targetStage, targetStageFlag, "", targetStageUsage)
	cmd.Flags().StringVar(&o.targetField, targetFieldFlag, "", targetFieldUsage)
}

// complete fills any remaining information for the runner that is not specified
// by the flags.
func (o *GetLinkOptions) complete(cmd *cobra.Command, args []string) error {
	if len(args) == 1 {
		o.name = args[0]
	}
	o.outWriter = cmd.OutOrStdout()
	return nil
}

// validate checks if the user options are compatible and the command can
// be executed
func (o *GetLinkOptions) validate() error {
	return nil
}

// run executes the get link command
func (o *GetLinkOptions) run() error {
	query := &pb.Link{
		Name:        o.name,
		SourceStage: o.sourceStage,
		SourceField: o.sourceField,
		TargetStage: o.targetStage,
		TargetField: o.targetField,
	}

	conn := client.NewConnection(o.addr)
	defer conn.Close()

	c := pb.NewLinkManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := c.Get(ctx, query)
	if err != nil {
		return client.ErrorFromGrpcError(err)
	}
	links := make([]*pb.Link, 0)
	for {
		l, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return client.ErrorFromGrpcError(err)
		}
		links = append(links, l)
	}
	return o.displayLinks(links)
}

func (o *GetLinkOptions) displayLinks(links []*pb.Link) error {
	sort.Slice(
		links,
		func(i, j int) bool {
			return links[i].Name < links[j].Name
		})
	numLinks := len(links)
	// Add space for all assets plus the header
	data := make([][]string, 0, numLinks+1)
	headers := []string{
		NameText,
		SourceStageText,
		SourceFieldText,
		TargetStageText,
		TargetFieldText,
	}
	data = append(data, headers)
	for _, l := range links {
		linkData := []string{
			l.Name,
			l.SourceStage,
			l.SourceField,
			l.TargetStage,
			l.TargetField,
		}
		data = append(data, linkData)
	}
	output, err := pterm.DefaultTable.WithHasHeader().WithData(data).Srender()
	if err != nil {
		return errdefs.UnknownWithMsg("display links: %v", err)
	}
	_, err = fmt.Fprintln(o.outWriter, output)
	if err != nil {
		return errdefs.UnknownWithMsg("display links: %v", err)
	}
	return nil
}
