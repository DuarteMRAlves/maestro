package maestro

import (
	"errors"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/yaml"
	"github.com/spf13/cobra"
	"io"
	"log"
	"path"
)

type ConvertOpts struct {
	inFile  string
	outFile string

	outW io.Writer
}

func NewConvertCmd() *cobra.Command {
	var opts ConvertOpts

	cmd := cobra.Command{
		Use:                   "convert [OPTIONS] v0-file",
		DisableFlagsInUseLine: true,
		Short:                 "Convert configs from version 0 to version 1",
		Long:                  "Convert a configuration file from version 0 to version 1.",
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			if err = opts.complete(cmd, args); err != nil {
				if _, writeErr := fmt.Fprintln(opts.outW, err); writeErr != nil {
					log.Fatalf("write error at convert command: %s\n", writeErr)
				}
			}
			if err = opts.run(); err != nil {
				if _, writeErr := fmt.Fprintln(opts.outW, err); writeErr != nil {
					log.Fatalf("write error at convert command: %s\n", writeErr)
				}
			}
		},
	}

	cmd.Flags().StringVarP(
		&opts.outFile, "output", "o", "", "name for the output file (defaults to conv-<in-file>",
	)

	return &cmd
}

func (o *ConvertOpts) complete(cmd *cobra.Command, args []string) error {
	o.outW = cmd.OutOrStdout()
	if len(args) != 1 {
		return errors.New("invalid number of arguments: expected one with source file")
	}
	o.inFile = args[0]

	if o.outFile == "" {
		dir, f := path.Split(o.inFile)
		conv := fmt.Sprintf("conv-%s", f)
		o.outFile = path.Join(dir, conv)
	}
	return nil
}

func (o *ConvertOpts) run() error {
	resources, err := yaml.ReadV0(o.inFile)
	if err != nil {
		return err
	}
	return yaml.WriteV1(resources, o.outFile, 0644)
}
