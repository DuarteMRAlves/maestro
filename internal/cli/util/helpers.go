package util

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/spf13/cobra"
)

func VerifyFlagsChanged(cmd *cobra.Command, flags []string) error {
	for _, f := range flags {
		if cmd.Flag(f).Changed {
			return errdefs.InvalidArgumentWithMsg("%v flag required", f)
		}
	}
	return nil
}

func WarnArgsIgnore(args []string, msg string) {
	if len(args) > 0 {
		fmt.Printf("warning: positional arguments ignored (%v)", msg)
	}
}

func WarnFlagsIgnore(cmd *cobra.Command, flags []string, msg string) {
	for _, f := range flags {
		if cmd.Flag(f).Changed {
			fmt.Printf("warning: %v flag ignored (%v)\n", f, msg)
		}
	}
}
