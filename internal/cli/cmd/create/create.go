package create

import (
	"fmt"
	"github.com/spf13/cobra"
)

var createOpts = struct {
	addr  string
	files []string
}{}

func NewCmdCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create resources of a given type",
		Args:  cobra.MaximumNArgs(0),
		Run:   RunCreate,
	}

	cmd.PersistentFlags().StringVar(
		&createOpts.addr,
		addrFull,
		addrDefault,
		addrHelp)

	cmd.PersistentFlags().StringSliceVarP(
		&createOpts.files,
		fileFull,
		fileShort,
		fileDefault,
		fileHelp)

	// Subcommands
	cmd.AddCommand(NewCmdCreateAsset())
	cmd.AddCommand(NewCmdCreateStage())
	cmd.AddCommand(NewCmdCreateLink())

	return cmd
}

func RunCreate(_ *cobra.Command, _ []string) {
	err := createFromFiles(createOpts.files, createOpts.addr, "")
	if err != nil {
		fmt.Printf("unable to create resources: %v\n", err)
	}
}
