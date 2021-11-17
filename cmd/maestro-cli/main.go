package main

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/cli/cmd"
	"os"
)

func main() {
	rootCmd := cmd.NewCmdRoot()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
