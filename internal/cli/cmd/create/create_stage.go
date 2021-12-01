package create

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/cli/client"
	"github.com/DuarteMRAlves/maestro/internal/cli/resources"
	"github.com/DuarteMRAlves/maestro/internal/cli/util"
	"github.com/spf13/cobra"
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
	addr    string
	files   []string
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

	addAddrFlag(cmd, &createStageOpts.addr, addrHelp)
	addFilesFlag(cmd, &createStageOpts.files, fileHelp)

	cmd.Flags().StringVarP(
		&createStageOpts.asset,
		assetFlag,
		assetShort,
		"",
		assetUsage)
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

func runCreateStage(cmd *cobra.Command, args []string) {
	var stage *pb.Stage

	if cmd.Flag(fileFull).Changed {
		util.WarnArgsIgnore(args, "creating from file")
		util.WarnFlagsIgnore(
			cmd,
			[]string{assetFlag, serviceFlag, methodFlag},
			"creating from file")
		err := createFromFiles(
			createStageOpts.files,
			createStageOpts.addr,
			resources.StageKind)
		if err != nil {
			fmt.Printf("unable to create resources: %v\n", err)
		}
	} else {
		if len(args) != 1 {
			fmt.Printf("unable to create stage: expected name argument")
			return
		}
		err := util.VerifyFlagsChanged(cmd, []string{assetFlag})
		if err != nil {
			fmt.Printf("unable to create stage: %v", err)
			return
		}
		stage = &pb.Stage{
			Name:    args[0],
			Asset:   createStageOpts.asset,
			Service: createStageOpts.service,
			Method:  createStageOpts.method,
		}
		err = client.CreateStage(stage, createStageOpts.addr)
		if err != nil {
			fmt.Printf("unable to create stage: %v\n", err)
			return
		}
		fmt.Printf("stage '%v' created\n", args[0])
	}
}
