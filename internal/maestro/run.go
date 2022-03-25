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
	"github.com/DuarteMRAlves/maestro/internal/yaml"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type configVersion string

const (
	v0 configVersion = "v0"
	v1 configVersion = "v1"
)

type RunOpts struct {
	files    []string
	orchName string
	v0       bool
	v1       bool

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

	cmd.Flags().BoolVar(&opts.v0, "v0", false, "use version 0 for config yaml format")
	cmd.Flags().BoolVar(&opts.v1, "v1", false, "use version 1 for config yaml format")
	cmd.Flags().StringArrayVarP(&opts.files, "file", "f", nil, "config files")

	return &cmd
}

func (opts *RunOpts) complete(cmd *cobra.Command, args []string) error {
	opts.outWriter = cmd.OutOrStdout()
	if len(args) > 1 {
		return errors.New("too many arguments: expected at most one")
	}
	if len(args) == 1 {
		opts.orchName = args[0]
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
	)
	switch opts.version {
	case v0:
		resources, err = yaml.ReadV0(opts.files[0])
	case v1:
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

	orchStore := make(mapstore.Orchestrations, len(resources.Orchestrations))
	stageStore := make(mapstore.Stages, len(resources.Stages))
	linkStore := make(mapstore.Links, len(resources.Links))

	createOrchestration := create.Orchestration(orchStore)
	createStage := create.Stage(stageStore, orchStore)
	createLink := create.Link(linkStore, stageStore, orchStore)

	for _, o := range resources.Orchestrations {
		if err := createOrchestration(o.Name); err != nil {
			return err
		}
	}
	for _, s := range resources.Stages {
		m := internal.NewMethodContext(
			s.Method.Address, s.Method.Service, s.Method.Method,
		)
		if err := createStage(s.Name, m, s.Orchestration); err != nil {
			return err
		}
	}
	for _, l := range resources.Links {
		s := internal.NewLinkEndpoint(l.Source.Stage, l.Source.Field)
		t := internal.NewLinkEndpoint(l.Target.Stage, l.Target.Field)
		if err := createLink(l.Name, s, t, l.Orchestration); err != nil {
			return err
		}
	}

	availableOrchs := arrays.Map(
		func(o yaml.Orchestration) internal.OrchestrationName { return o.Name },
		resources.Orchestrations...,
	)

	orchName, err := opts.orchToRun(availableOrchs...)
	if err != nil {
		return err
	}

	orch, err := orchStore.Load(orchName)
	if err != nil {
		return err
	}

	logger := logs.New(false)
	b := execute.NewBuilder(stageStore, linkStore, grpc.ReflectionMethodLoader, logger)
	execution, err := b(orch)
	if err != nil {
		return err
	}

	errs := make(chan error, 1)
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		errs <- execution.Stop()
	}()

	execution.Start()

	err = <-errs
	return err
}

func (opts *RunOpts) orchToRun(
	available ...internal.OrchestrationName,
) (internal.OrchestrationName, error) {
	if opts.orchName != "" {
		pred := func(v internal.OrchestrationName) bool {
			return v.Unwrap() == opts.orchName
		}
		available = arrays.Filter(pred, available...)
	}
	switch len(available) {
	case 0:
		var err error
		if opts.orchName != "" {
			err = fmt.Errorf("orchestration %s not found", opts.orchName)
		} else {
			err = errors.New("no orchestrations defined")
		}
		return internal.OrchestrationName{}, err
	case 1:
		return available[0], nil
	default:
		err := fmt.Errorf("only one orchestration can be executed but found %s", available)
		return internal.OrchestrationName{}, err
	}
}
