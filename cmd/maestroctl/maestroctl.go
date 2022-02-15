package main

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/cli/maestroctl/cmd"
)

func main() {
	rootCmd := cmd.NewCmdRoot()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
