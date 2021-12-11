package get

import "github.com/spf13/cobra"

const (
	addrFull    = "addr"
	addrDefault = "localhost:50051"
	addrHelp    = "address to connect to the maestro server"
)

func addAddrFlag(cmd *cobra.Command, value *string) {
	cmd.Flags().StringVar(
		value,
		addrFull,
		addrDefault,
		addrHelp)
}
