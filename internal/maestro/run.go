package maestro

import (
	"errors"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/yaml"
	"github.com/spf13/cobra"
	"io"
	"log"
)

type configVersion string

const (
	v0 configVersion = "v0"
	v1 configVersion = "v1"
)

type RunOpts struct {
	files      []string
	versionArg string
	orchName   string

	outWriter io.Writer
	version   configVersion
}

func NewRunCmd() *cobra.Command {
	var opts RunOpts

	cmd := cobra.Command{
		Use:                   "run [OPTIONS] [ORCHESTRATION]",
		DisableFlagsInUseLine: true,
		Short:                 "Execute a single orchestration",
		Long: `Execute a single orchestration from configuration files.

If no orchestration is specified, the configuration files should only contain
a single orchestration, that will be executed.`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			if err = opts.complete(cmd, args); err != nil {
				if _, writeErr := fmt.Fprintln(opts.outWriter, err); writeErr != nil {
					log.Fatalf("write error at run command: %s\n", writeErr)
				}
				return
			}
			if err = opts.validate(); err != nil {
				if _, writeErr := fmt.Fprintln(opts.outWriter, err); writeErr != nil {
					log.Fatalf("write error at run command: %s\n", writeErr)
				}
				return
			}
			if err = opts.run(); err != nil {
				if _, writeErr := fmt.Fprintln(opts.outWriter, err); writeErr != nil {
					log.Fatalf("write error at run command: %s\n", writeErr)
				}
				return
			}
		},
	}

	cmd.Flags().StringVar(
		&opts.versionArg, "conf-version", string(v1), "version for the config yaml format",
	)
	cmd.Flags().StringArrayVarP(&opts.files, "file", "f", nil, "config files")

	return &cmd
}

func (o *RunOpts) complete(cmd *cobra.Command, args []string) error {
	o.outWriter = cmd.OutOrStdout()
	if len(args) > 1 {
		return errors.New("too many arguments: expected at most one")
	}
	if len(args) == 1 {
		o.orchName = args[0]
	}
	o.version = configVersion(o.versionArg)
	switch o.version {
	case v0:
	case v1:
	default:
		return fmt.Errorf(
			"unknown config version: expected %s or %s but found %s", v0, v1, o.versionArg,
		)
	}
	return nil
}

func (o *RunOpts) validate() error {
	if len(o.files) == 0 {
		return errors.New("specify at least one configuration file")
	}
	if o.version == v0 && len(o.files) > 1 {
		return errors.New("only one configuration file allowed for v0 file specification")
	}
	return nil
}

func (o *RunOpts) run() error {
	var (
		resources yaml.ResourceSet
		err       error
	)
	switch o.version {
	case v0:
		resources, err = yaml.ReadV0(o.files[0])
	case v1:
		resources, err = yaml.ReadV1(o.files...)
	default:
		return fmt.Errorf(
			"unknown config version: expected %s or %s but found %s", v0, v1, o.versionArg,
		)
	}
	if err != nil {
		return err
	}
	fmt.Println(resources)
	return nil
}
