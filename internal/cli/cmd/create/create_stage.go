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
		"(if not specified the service must only have on method to run)"
)

// CreateStageOptions stores the flags for the CreateStage command and executes it
type CreateStageOptions struct {
	addr string

	name    string
	asset   string
	service string
	method  string
}

// NewCmdCreateStage returns a new command that creates a stage from command
// line arguments
func NewCmdCreateStage() *cobra.Command {
	o := &CreateStageOptions{}

	cmd := &cobra.Command{
		Use:   "stage name",
		Short: "create a new stage",
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
// run the CreateStage command
func (o *CreateStageOptions) addFlags(cmd *cobra.Command) {
	addAddrFlag(cmd, &o.addr, addrHelp)

	cmd.Flags().StringVarP(&o.asset, assetFlag, assetShort, "", assetUsage)
	cmd.Flags().StringVarP(
		&o.service,
		serviceFlag,
		serviceShort,
		"",
		serviceUsage)
	cmd.Flags().StringVarP(&o.method, methodFlag, methodShort, "", methodUsage)
}

// complete fills any remaining information necessary to run the command that is
// not specified by the user flags and is in the positional arguments
func (o *CreateStageOptions) complete(args []string) error {
	if len(args) == 1 {
		o.name = args[0]
	}
	return nil
}

// validate verifies if the user options are valid and all necessary information
// for the command to run is present
func (o *CreateStageOptions) validate() error {
	if o.name == "" {
		return errdefs.InvalidArgumentWithMsg("please specify a stage name")
	}
	if o.asset == "" {
		return errdefs.InvalidArgumentWithMsg("please specify an asset")
	}
	return nil
}

func (o *CreateStageOptions) run() error {
	stage := &resources.StageResource{
		Name:    o.name,
		Asset:   o.asset,
		Service: o.service,
		Method:  o.method,
	}
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second)
	defer cancel()

	return client.CreateStage(ctx, stage, o.addr)
}
