package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var roodCmd = &cobra.Command{
	Use:   "maestro-cli",
	Short: "maestro-cli is a command line interface to communicate with maestro",
}

func Execute() {
	if err := roodCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
