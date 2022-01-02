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

	addressFlag  = "address"
	addressUsage = "the address where the stage service is running"
)

// StageOpts stores the flags for the CreateStage command and executes it
type StageOpts struct {
	// address for the maestro server
	maestro string

	name    string
	asset   string
	service string
	method  string
	address string
	host    string
	port    int32
}

// NewCmdCreateStage returns a new command that creates a stage from command
// line arguments
func NewCmdCreateStage() *cobra.Command {
	o := &StageOpts{}

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
func (o *StageOpts) addFlags(cmd *cobra.Command) {
	util.AddMaestroFlag(cmd, &o.maestro)

	cmd.Flags().StringVarP(&o.asset, assetFlag, assetShort, "", assetUsage)
	cmd.Flags().StringVarP(
		&o.service,
		serviceFlag,
		serviceShort,
		"",
		serviceUsage)
	cmd.Flags().StringVarP(&o.method, methodFlag, methodShort, "", methodUsage)
	cmd.Flags().StringVar(&o.address, addressFlag, "", addressUsage)
	cmd.Flags().StringVar(&o.host, "host", "", "host where service is running")
	cmd.Flags().Int32Var(&o.port, "port", 0, "port where service is running")
}

// complete fills any remaining information necessary to run the command that is
// not specified by the user flags and is in the positional arguments
func (o *StageOpts) complete(args []string) error {
	if len(args) == 1 {
		o.name = args[0]
	}
	return nil
}

// validate verifies if the user options are valid and all necessary information
// for the command to run is present
func (o *StageOpts) validate() error {
	if o.name == "" {
		return errdefs.InvalidArgumentWithMsg("please specify a stage name")
	}
	if o.address != "" && o.host != "" {
		return errdefs.InvalidArgumentWithMsg(
			"address and host options are incompatible")
	}
	if o.address != "" && o.port != 0 {
		return errdefs.InvalidArgumentWithMsg(
			"address and port options are incompatible")
	}
	return nil
}

func (o *StageOpts) run() error {
	stage := &resources.StageSpec{
		Name:    o.name,
		Asset:   o.asset,
		Service: o.service,
		Method:  o.method,
		Address: o.address,
		Host:    o.host,
		Port:    o.port,
	}

	conn, err := grpc.Dial(o.maestro, grpc.WithInsecure())
	if err != nil {
		return errdefs.UnavailableWithMsg("create connection: %v", err)
	}
	defer conn.Close()

	c := client.New(conn)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second)
	defer cancel()

	return c.CreateStage(ctx, stage)
}
