package cmd

import (
	"github.com/DuarteMRAlves/maestro/internal/cli/maestroctl/cmd/attach"
	"github.com/DuarteMRAlves/maestro/internal/cli/maestroctl/cmd/create"
	"github.com/DuarteMRAlves/maestro/internal/cli/maestroctl/cmd/get"
	"github.com/spf13/cobra"
)

const shortDescription = "maestroctl is a command line interface to " +
	"communicate with maestro"

func NewCmdRoot() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "maestroctl",
		Short: shortDescription,
	}

	// Subcommands
	rootCmd.AddCommand(attach.NewCmdAttach())
	rootCmd.AddCommand(create.NewCmdCreate())
	rootCmd.AddCommand(get.NewCmdGet())
	return rootCmd
}
