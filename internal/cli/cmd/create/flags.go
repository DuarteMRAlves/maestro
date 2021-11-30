package create

import (
	"github.com/spf13/cobra"
)

const (
	addrFull    = "addr"
	addrDefault = "localhost:50051"
	addrHelp    = "address to connect to the maestro server"

	fileFull  = "file"
	fileShort = "f"
	fileHelp  = "files to create one or more resources"
)

func addAddrFlag(cmd *cobra.Command, value *string, usage string) {
	cmd.Flags().StringVar(
		value,
		addrFull,
		addrDefault,
		usage)
}

func addFilesFlag(cmd *cobra.Command, value *[]string, usage string) {
	cmd.Flags().StringSliceVarP(
		value,
		fileFull,
		fileShort,
		nil,
		usage)
}
