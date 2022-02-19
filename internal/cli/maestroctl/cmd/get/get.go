package get

import "github.com/spf13/cobra"

func NewCmdGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [COMMAND]",
		Short: "Display one or more resources",
		Long:  "Display relevant information related to one or more resources.",
	}

	// Subcommands
	cmd.AddCommand(NewCmdGetAsset())
	cmd.AddCommand(NewCmdGetStage())
	cmd.AddCommand(NewCmdGetLink())
	cmd.AddCommand(NewCmdGetOrchestration())

	return cmd
}
