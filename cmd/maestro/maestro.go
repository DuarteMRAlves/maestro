package main

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/maestro"
	"os"
)

func main() {
	cmd := maestro.RootCmd()
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
