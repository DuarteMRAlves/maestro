package get

import "github.com/spf13/cobra"

func NewCmdGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "list one or more resources of a given type",
	}

	// Subcommands
	cmd.AddCommand(NewCmdGetAsset())
	cmd.AddCommand(NewCmdGetStage())
	cmd.AddCommand(NewCmdGetLink())

	return cmd
}
