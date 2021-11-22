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
	assetFlag  = "asset"
	assetShort = "a"
	assetUsage = "asset name to search"

	serviceFlag  = "service"
	serviceShort = "s"
	serviceUsage = "service name to search"

	methodFlag  = "method"
	methodShort = "m"
	methodUsage = "method name to search"
)

var getStageOpts = &struct {
	asset   string
	service string
	method  string
}{}

func NewCmdGetStage() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stage",
		Short:   "list one or more stages",
		Args:    cobra.MaximumNArgs(1),
		Aliases: []string{"stages"},
		Run:     runGetStage,
	}
	cmd.Flags().StringVarP(
		&getStageOpts.asset,
		assetFlag,
		assetShort,
		"",
		assetUsage)
	cmd.Flags().StringVarP(
		&getStageOpts.service,
		serviceFlag,
		serviceShort,
		"",
		serviceUsage)
	cmd.Flags().StringVarP(
		&getStageOpts.method,
		methodFlag,
		methodShort,
		"",
		methodUsage)

	return cmd
}

func runGetStage(_ *cobra.Command, args []string) {
	query := &pb.Stage{
		Asset:   getStageOpts.asset,
		Service: getStageOpts.service,
		Method:  getStageOpts.method,
	}

	if len(args) == 1 {
		query.Name = args[0]
	}

	conn := cli.NewConnection(createOpts.addr)
	defer conn.Close()

	c := pb.NewStageManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := c.Get(ctx, query)
	if err != nil {
		log.Fatalf("list stages: %v", err)
	}
	stages := make([]*pb.Stage, 0)
	for {
		s, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("list stages: %v", err)
		}
		stages = append(stages, s)
	}
	if err := displayStages(stages); err != nil {
		log.Fatalf("list stages: %v", err)
	}
}

func displayStages(stages []*pb.Stage) error {
	numStages := len(stages)
	names := make([]string, 0, numStages)
	assets := make([]string, 0, numStages)
	services := make([]string, 0, numStages)
	methods := make([]string, 0, numStages)

	for _, s := range stages {
		names = append(names, s.Name)
		assets = append(assets, s.Asset)
		services = append(services, s.Service)
		methods = append(methods, s.Method)
	}

	t := table.NewBuilder().
		WithPadding(colPad).
		WithMinColSize(minColSize).
		Build()

	if err := t.AddColumn(NameText, names); err != nil {
		return err
	}
	if err := t.AddColumn(AssetText, assets); err != nil {
		return err
	}
	if err := t.AddColumn(ServiceText, services); err != nil {
		return err
	}
	if err := t.AddColumn(MethodText, methods); err != nil {
		return err
	}
	fmt.Print(t)
	return nil
}
