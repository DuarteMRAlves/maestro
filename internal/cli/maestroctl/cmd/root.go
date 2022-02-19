package cmd

import (
	"github.com/DuarteMRAlves/maestro/internal/cli/maestroctl/cmd/attach"
	"github.com/DuarteMRAlves/maestro/internal/cli/maestroctl/cmd/create"
	"github.com/DuarteMRAlves/maestro/internal/cli/maestroctl/cmd/get"
	"github.com/DuarteMRAlves/maestro/internal/cli/maestroctl/cmd/start"
	"github.com/spf13/cobra"
)

func NewCmdRoot() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "maestroctl",
		Short: "maestroctl controls the Maestro server",
		Long:  "maestroctl is a command line interface to control a Maestro server.",
	}

	// Subcommands
	rootCmd.AddCommand(attach.NewCmdAttach())
	rootCmd.AddCommand(create.NewCmdCreate())
	rootCmd.AddCommand(get.NewCmdGet())
	rootCmd.AddCommand(start.NewCmdStart())
	return rootCmd
}
