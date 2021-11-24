package get

import "github.com/spf13/cobra"

const (
	defaultAddr = "localhost:50051"
	addrUsage   = "Address to connect to the maestro server"
)

var createOpts = struct {
	addr string
}{}

func NewCmdGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "list one or more resources of a given type",
	}

	cmd.PersistentFlags().StringVar(
		&createOpts.addr,
		"addr",
		defaultAddr,
		addrUsage)

	// Subcommands
	cmd.AddCommand(NewCmdGetAsset())
	cmd.AddCommand(NewCmdGetStage())
	cmd.AddCommand(NewCmdGetLink())

	return cmd
}
