package util

import (
	"fmt"
	"github.com/spf13/cobra"
)

func WriteOut(cmd *cobra.Command, format string, args ...interface{}) {
	_, err := fmt.Fprintf(cmd.OutOrStdout(), fmt.Sprintf(format, args...))
	if err != nil {
		fmt.Printf("unable to write to cmd output: %v", err)
	}
}
