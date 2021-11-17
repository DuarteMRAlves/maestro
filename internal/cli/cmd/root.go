package cmd

import (
	"github.com/DuarteMRAlves/maestro/internal/cli/cmd/create"
	"github.com/DuarteMRAlves/maestro/internal/cli/cmd/get"
	"github.com/spf13/cobra"
)

const shortDescription = "maestro-cli is a command line interface to " +
	"communicate with maestro"

func NewCmdRoot() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "maestro-cli",
		Short: shortDescription,
	}

	// Subcommands
	rootCmd.AddCommand(create.NewCmdCreate())
	rootCmd.AddCommand(get.NewCmdGet())
	return rootCmd
}
