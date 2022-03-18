package create

import (
	"context"
	"fmt"
	util2 "github.com/DuarteMRAlves/maestro/_old/internal/cli/maestroctl/cmd/util"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"io"
	"time"
)

// Options store the flags defined by the user when executing the create
// command and then executes the command.
type Options struct {
	// address for the maestro server
	maestro string

	// names of the resources to create
	names []string
	files []string

	// Output for the cobra.Command to be executed.
	outWriter io.Writer
}

func NewCmdCreate() *cobra.Command {
	o := &Options{}

	cmd := &cobra.Command{
		Use:                   "create -f FILENAME ... [FLAGS]",
		DisableFlagsInUseLine: true,
		Short:                 "Create resources from files",
		Long:                  "Create resources from files in yaml format.",
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			err = o.complete(cmd, args)
			if err != nil {
				util2.WriteOut(cmd, util2.DisplayMsgFromError(err))
				return
			}
			err = o.validate()
			if err != nil {
				util2.WriteOut(cmd, util2.DisplayMsgFromError(err))
				return
			}
			err = o.run()
			if err != nil {
				util2.WriteOut(cmd, util2.DisplayMsgFromError(err))
			}
		},
	}

	o.addFlags(cmd)

	return cmd
}

// addFlags adds the necessary flags to the cobra.Command instance that will
// parse the command line arguments and run the command
func (o *Options) addFlags(cmd *cobra.Command) {
	util2.AddMaestroFlag(cmd, &o.maestro)
	util2.AddFilesFlag(cmd, &o.files, "files to create one or more resources")
}

// complete fills any remaining information that is required to execute the
// create command.
func (o *Options) complete(cmd *cobra.Command, args []string) error {
	o.names = args
	o.outWriter = cmd.OutOrStdout()
	return nil
}

// validate verifies if the user inputs are valid and there are no conflits
func (o *Options) validate() error {
	// In create, we only accept files
	if len(o.files) == 0 {
		return errdefs.InvalidArgumentWithMsg("please specify input files")
	}
	return nil
}

// run executes the Create command
func (o *Options) run() error {
	parsed, err := ParseFiles(o.files)
	if err != nil {
		return err
	}
	if err = IsValidKinds(parsed); err != nil {
		return err
	}

	assets, err := collectAssetSpecs(parsed)
	if err != nil {
		return errdefs.PrependMsg(err, "create")
	}
	orchestrations, err := collectOrchestrationSpecs(parsed)
	if err != nil {
		return errdefs.PrependMsg(err, "create")
	}
	stages, err := collectStageSpecs(parsed)
	if err != nil {
		return errdefs.PrependMsg(err, "create")
	}
	links, err := collectLinkSpecs(parsed)
	if err != nil {
		return errdefs.PrependMsg(err, "create")
	}

	if len(o.names) != 0 {
		seen := make(map[string]bool)
		for _, n := range o.names {
			seen[n] = false
		}
		assets = filterAssets(assets, seen)
		orchestrations = filterOrchestrations(orchestrations, seen)
		stages = filterStages(stages, seen)
		links = filterLinks(links, seen)
		for n, s := range seen {
			if !s {
				_, err = fmt.Fprintf(
					o.outWriter,
					"warning: resource '%s' not founc",
					n,
				)
				if err != nil {
					return errdefs.UnknownWithMsg("create: %v", err)
				}
			}
		}
	}

	conn, err := grpc.Dial(o.maestro, grpc.WithInsecure())
	if err != nil {
		return errdefs.UnavailableWithMsg("create connection: %v", err)
	}
	defer conn.Close()

	archStub := pb.NewArchitectureManagementClient(conn)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second,
	)
	defer cancel()

	for _, req := range assets {
		a := &pb.CreateAssetRequest{
			Name:  string(req.Name),
			Image: req.Image,
		}
		if _, err = archStub.CreateAsset(ctx, a); err != nil {
			return util2.ErrorFromGrpcError(err)
		}
	}
	for _, req := range orchestrations {
		pbReq := &pb.CreateOrchestrationRequest{
			Name: string(req.Name),
		}
		if _, err = archStub.CreateOrchestration(ctx, pbReq); err != nil {
			return util2.ErrorFromGrpcError(err)
		}
	}
	for _, req := range stages {
		pbReq := &pb.CreateStageRequest{
			Name:          req.Name,
			Asset:         req.Asset,
			Service:       req.Service,
			Rpc:           req.Rpc,
			Address:       req.Address,
			Host:          req.Host,
			Port:          req.Port,
			Orchestration: string(req.Orchestration),
		}
		if _, err = archStub.CreateStage(ctx, pbReq); err != nil {
			return util2.ErrorFromGrpcError(err)
		}
	}
	for _, req := range links {
		l := &pb.CreateLinkRequest{
			Name:          string(req.Name),
			SourceStage:   string(req.SourceStage),
			SourceField:   req.SourceField,
			TargetStage:   string(req.TargetStage),
			TargetField:   req.TargetField,
			Orchestration: string(req.Orchestration),
		}
		if _, err = archStub.CreateLink(ctx, l); err != nil {
			return util2.ErrorFromGrpcError(err)
		}
	}

	return nil
}

func filterAssets(specs []*AssetSpec, seen map[string]bool) []*AssetSpec {
	filtered := make([]*AssetSpec, 0)
	for _, s := range specs {
		name := s.Name
		_, exists := seen[name]
		if exists {
			filtered = append(filtered, s)
			seen[name] = true
		}
	}
	return filtered
}

func filterOrchestrations(
	specs []*OrchestrationSpec,
	seen map[string]bool,
) []*OrchestrationSpec {
	filtered := make([]*OrchestrationSpec, 0)
	for _, s := range specs {
		name := s.Name
		_, exists := seen[name]
		if exists {
			filtered = append(filtered, s)
			seen[name] = true
		}
	}
	return filtered
}

func filterStages(specs []*StageSpec, seen map[string]bool) []*StageSpec {
	filtered := make([]*StageSpec, 0)
	for _, s := range specs {
		name := s.Name
		_, exists := seen[name]
		if exists {
			filtered = append(filtered, s)
			seen[name] = true
		}
	}
	return filtered
}

func filterLinks(specs []*LinkSpec, seen map[string]bool) []*LinkSpec {
	filtered := make([]*LinkSpec, 0)
	for _, s := range specs {
		name := s.Name
		_, exists := seen[name]
		if exists {
			filtered = append(filtered, s)
			seen[name] = true
		}
	}
	return filtered
}
