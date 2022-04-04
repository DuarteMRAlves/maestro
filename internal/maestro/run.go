package maestro

import (
	"errors"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/internal/arrays"
	"github.com/DuarteMRAlves/maestro/internal/create"
	"github.com/DuarteMRAlves/maestro/internal/execute"
	"github.com/DuarteMRAlves/maestro/internal/grpc"
	"github.com/DuarteMRAlves/maestro/internal/logs"
	"github.com/DuarteMRAlves/maestro/internal/mapstore"
	"github.com/DuarteMRAlves/maestro/internal/retry"
	"github.com/DuarteMRAlves/maestro/internal/yaml"
	"github.com/spf13/cobra"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
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
		resources yaml.ResourceSet
		err       error
		backoff   retry.ExponentialBackoff
	)
	switch opts.version {
	case v0:
		opts.logger.Debugf("read v0 from file %s", opts.files[0])
		resources, err = yaml.ReadV0(opts.files[0])
	case v1:
		opts.logger.Debugf("read v1 from files %s", opts.files)
		resources, err = yaml.ReadV1(opts.files...)
	default:
		// Should never happen if command was completed and validated.
		return fmt.Errorf(
			"unknown config version: expected %s or %s but found %s", v0, v1, opts.version,
		)
	}
	if err != nil {
		return err
	}

	pipelineStore := make(mapstore.Pipelines, len(resources.Pipelines))
	stageStore := make(mapstore.Stages, len(resources.Stages))
	linkStore := make(mapstore.Links, len(resources.Links))

	createPipeline := create.Pipeline(pipelineStore)
	createStage := create.Stage(stageStore, pipelineStore)
	createLink := create.Link(linkStore, stageStore, pipelineStore)

	for _, o := range resources.Pipelines {
		// FIXME: remove online execution hardcoded.
		if err := createPipeline(o.Name, internal.OnlineExecution); err != nil {
			return err
		}
	}
	for _, s := range resources.Stages {
		m := internal.NewMethodContext(
			s.Method.Address, s.Method.Service, s.Method.Method,
		)
		if err := createStage(s.Name, m, s.Pipeline); err != nil {
			return err
		}
	}
	for _, l := range resources.Links {
		s := internal.NewLinkEndpoint(l.Source.Stage, l.Source.Field)
		t := internal.NewLinkEndpoint(l.Target.Stage, l.Target.Field)
		if err := createLink(l.Name, s, t, l.Pipeline); err != nil {
			return err
		}
	}

	availablePipelines := arrays.Map(
		func(o yaml.Pipeline) internal.PipelineName { return o.Name },
		resources.Pipelines...,
	)

	pipelineName, err := opts.pipelineToRun(availablePipelines...)
	if err != nil {
		return err
	}

	pipeline, err := pipelineStore.Load(pipelineName)
	if err != nil {
		return err
	}

	r := grpc.NewReflectionMethodLoader(5*time.Minute, backoff, opts.logger)
	b := execute.NewBuilder(stageStore, linkStore, r, opts.logger)
	execution, err := b(pipeline)
	if err != nil {
		return err
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

func (opts *RunOpts) pipelineToRun(
	available ...internal.PipelineName,
) (internal.PipelineName, error) {
	if opts.pipelineName != "" {
		pred := func(v internal.PipelineName) bool {
			return v.Unwrap() == opts.pipelineName
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
		return internal.PipelineName{}, err
	case 1:
		return available[0], nil
	default:
		err := fmt.Errorf("only one pipeline can be executed but found %s", available)
		return internal.PipelineName{}, err
	}
}
