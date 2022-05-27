package maestro

import "github.com/spf13/cobra"

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "maestro COMMAND [OPTIONS]",
		Short: "maestro is a tool to execute grpc pipelines",
	}

	cmd.AddCommand(NewRunCmd(), NewConvertCmd())
	return cmd
}
