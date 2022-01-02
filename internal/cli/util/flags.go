package util

import "github.com/spf13/cobra"

// AddAddrFlag adds the flag that should be used to ask for the address to
// connect to the maestro server.
func AddAddrFlag(cmd *cobra.Command, value *string) {
	cmd.Flags().StringVar(
		value,
		"addr",
		"localhost:50051",
		"address to connect to the maestro server")
}

// AddFilesFlag adds a flag that should be used to ask for multiple files
// from the user. The usage is a parameter as different commands may use files
// for different purposes, thus requiring different help messages.
func AddFilesFlag(cmd *cobra.Command, value *[]string, usage string) {
	cmd.Flags().StringSliceVarP(value, "file", "f", nil, usage)
}
