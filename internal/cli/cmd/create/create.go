package create

import "github.com/spf13/cobra"

const (
	defaultAddr = "localhost:50051"
	addrUsage   = "Address to connect to the maestro server"
)

var createOpts = struct {
	addr string
}{}

func NewCmdCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create resources of a given type",
	}

	cmd.PersistentFlags().StringVar(
		&createOpts.addr,
		"addr",
		defaultAddr,
		addrUsage)

	// Subcommands
	cmd.AddCommand(NewCmdCreateAsset())
	cmd.AddCommand(NewCmdCreateStage())
	cmd.AddCommand(NewCmdCreateLink())

	return cmd
}
