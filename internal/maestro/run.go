package maestro

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/arrays"
	"github.com/DuarteMRAlves/maestro/internal/compiled"
	"github.com/DuarteMRAlves/maestro/internal/execute"
	"github.com/DuarteMRAlves/maestro/internal/grpcw"
	"github.com/DuarteMRAlves/maestro/internal/logs"
	"github.com/DuarteMRAlves/maestro/internal/repr"
	"github.com/DuarteMRAlves/maestro/internal/retry"
	"github.com/DuarteMRAlves/maestro/internal/yaml"
	"github.com/spf13/cobra"
)

type configVersion string

const (
	v0 configVersion = "v0"
	v1 configVersion = "v1"
)

type RunOpts struct {
	files        []string
	pipelineName string
	v0           bool
	v1           bool
	verbose      bool

	outWriter io.Writer
	version   configVersion
	logger    logs.Logger
}

func NewRunCmd() *cobra.Command {
	var opts RunOpts

	cmd := cobra.Command{
		Use:                   "run [OPTIONS] [PIPELINE]",
		DisableFlagsInUseLine: true,
		Short:                 "Execute a single pipeline",
		Long: `Execute a single pipeline from configuration files.

If no pipeline is specified, the configuration files should only contain
a single pipeline, that will be executed.`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			if err = opts.complete(cmd, args); err != nil {
				opts.logger.Infof("fatal: %s\n", err)
				os.Exit(1)
			}
			if err = opts.validate(); err != nil {
				opts.logger.Infof("fatal: %s\n", err)
				os.Exit(1)
			}
			if err = opts.run(); err != nil {
				opts.logger.Infof("fatal: %s\n", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().BoolVar(&opts.v0, "v0", false, "use version 0 for config yaml format")
	cmd.Flags().BoolVar(&opts.v1, "v1", false, "use version 1 for config yaml format")
	cmd.Flags().StringArrayVarP(&opts.files, "file", "f", nil, "config files")
	cmd.Flags().BoolVarP(&opts.verbose, "verbose", "v", false, "increase verbosity")

	return &cmd
}

func (opts *RunOpts) complete(cmd *cobra.Command, args []string) error {
	opts.outWriter = cmd.OutOrStdout()
	opts.logger = logs.NewWithOutput(opts.outWriter, opts.verbose)
	if len(args) > 1 {
		return errors.New("too many arguments: expected at most one positional argument")
	}
	if len(args) == 1 {
		opts.pipelineName = args[0]
	}
	if opts.v0 && opts.v1 {
		return errors.New("v0 and v1 options are incompatible")
	}
	// Defaults to v1
	opts.version = v1
	if opts.v0 {
		opts.version = v0
	}
	return nil
}

func (opts *RunOpts) validate() error {
	if len(opts.files) == 0 {
		return errors.New("specify at least one configuration file")
	}
	if opts.version == v0 && len(opts.files) > 1 {
		return errors.New("only one configuration file allowed for v0 file specification")
	}
	return nil
}

func (opts *RunOpts) run() error {
	var (
		pipelineCfg *api.Pipeline
		err         error
		backoff     retry.ExponentialBackoff
	)
	switch opts.version {
	case v0:
		opts.logger.Debugf("read v0 from file %s", opts.files[0])
		pipelineCfg, err = yaml.ReadV0(opts.files[0])
		if err != nil {
			return err
		}
	case v1:
		var pipelines []*api.Pipeline
		opts.logger.Debugf("read v1 from files %s", opts.files)
		pipelines, err = yaml.ReadV1(opts.files...)
		if err != nil {
			return err
		}
		pipelineCfg, err = opts.pipelineToRun(pipelines...)
		if err != nil {
			return err
		}
	default:
		// Should never happen if command was completed and validated.
		return fmt.Errorf(
			"unknown config version: expected %s or %s but found %s", v0, v1, opts.version,
		)
	}
	opts.logger.Infof("Pipeline Config:\n%s", repr.Pipeline(pipelineCfg))

	r, err := grpcw.NewReflectionResolver(time.Minute, backoff, opts.logger)
	if err != nil {
		return err
	}
	compilationCtx := compiled.NewContext(r)
	compiledPipeline, err := compiled.New(compilationCtx, pipelineCfg)
	if err != nil {
		return fmt.Errorf("compile %s: %w", pipelineCfg.Name, err)
	}
	b := execute.NewBuilder(opts.logger)
	execution, err := b(compiledPipeline)
	if err != nil {
		return fmt.Errorf("build execution %s: %w", pipelineCfg.Name, err)
	}

	errs := make(chan error, 1)
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		opts.logger.Infof("Received signal: %v", sig)
		errs <- execution.Stop()
	}()

	execution.Start()

	err = <-errs
	opts.logger.Debugf("Execution terminated with error: %s", err)
	return err
}

func (opts *RunOpts) pipelineToRun(available ...*api.Pipeline) (*api.Pipeline, error) {
	if opts.pipelineName != "" {
		pred := func(v *api.Pipeline) bool {
			return v.Name == opts.pipelineName
		}
		available = arrays.Filter(pred, available...)
	}
	switch len(available) {
	case 0:
		var err error
		if opts.pipelineName != "" {
			err = fmt.Errorf("pipeline %s not found", opts.pipelineName)
		} else {
			err = errors.New("no pipelines defined")
		}
		return nil, err
	case 1:
		return available[0], nil
	default:
		names := arrays.Map(
			func(o *api.Pipeline) string { return o.Name },
			available...,
		)
		err := fmt.Errorf("only one pipeline can be executed but found %s", names)
		return nil, err
	}
}
