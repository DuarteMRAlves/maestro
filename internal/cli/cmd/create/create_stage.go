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
	assetFlag  = "asset"
	assetShort = "a"
	assetUsage = "name of the asset the stage executes (required)"

	serviceFlag  = "service"
	serviceShort = "s"
	serviceUsage = "name of the grpc service to call " +
		"(if not specified the asset must only have one service)"

	methodFlag  = "method"
	methodShort = "m"
	methodUsage = "name of the grpc method to call " +
		"(if not specified the service must only have on method to execute)"
)

var createStageOpts = &struct {
	asset   string
	service string
	method  string
}{}

func NewCmdCreateStage() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stage name",
		Short: "create a new stage",
		Args:  cobra.ExactArgs(1),
		Run:   runCreateStage,
	}

	cmd.Flags().StringVarP(
		&createStageOpts.asset,
		assetFlag,
		assetShort,
		"",
		assetUsage)
	cmd.MarkFlagRequired(assetFlag)
	cmd.Flags().StringVarP(
		&createStageOpts.service,
		serviceFlag,
		serviceShort,
		"",
		serviceUsage)
	cmd.Flags().StringVarP(
		&createStageOpts.method,
		methodFlag,
		methodShort,
		"",
		methodUsage)

	return cmd
}

func runCreateStage(_ *cobra.Command, args []string) {
	conn := client.NewConnection(createOpts.addr)
	defer conn.Close()

	c := pb.NewStageManagementClient(conn)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second)
	defer cancel()

	_, err := c.Create(
		ctx,
		&pb.Stage{
			Name:    args[0],
			Asset:   createStageOpts.asset,
			Service: createStageOpts.service,
			Method:  createStageOpts.method,
		})
	if err != nil {
		fmt.Printf("unable to create stage: %v\n", err)
		return
	}
	fmt.Printf("stage '%v' created\n", args[0])
}
