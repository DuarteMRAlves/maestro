package create

import (
	"context"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/cli/maestroctl/cmd/util"
	"github.com/DuarteMRAlves/maestro/internal/cli/maestroctl/resources"
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
		Use:   "create",
		Short: "create resources of a given type",
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			err = o.complete(cmd, args)
			if err != nil {
				util.WriteOut(cmd, util.DisplayMsgFromError(err))
				return
			}
			err = o.validate()
			if err != nil {
				util.WriteOut(cmd, util.DisplayMsgFromError(err))
				return
			}
			err = o.run()
			if err != nil {
				util.WriteOut(cmd, util.DisplayMsgFromError(err))
			}
		},
	}

	o.addFlags(cmd)

	return cmd
}

// addFlags adds the necessary flags to the cobra.Command instance that will
// parse the command line arguments and run the command
func (o *Options) addFlags(cmd *cobra.Command) {
	util.AddMaestroFlag(cmd, &o.maestro)
	util.AddFilesFlag(cmd, &o.files, "files to create one or more resources")
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
	parsed, err := resources.ParseFiles(o.files)
	if err != nil {
		return err
	}
	if err = resources.IsValidKinds(parsed); err != nil {
		return err
	}

	assets, err := resources.FilterCreateAssetRequests(parsed)
	if err != nil {
		return errdefs.PrependMsg(err, "create")
	}
	orchestrations, err := resources.FilterCreateOrchestrationRequests(parsed)
	if err != nil {
		return errdefs.PrependMsg(err, "create")
	}
	stages, err := resources.FilterCreateStageRequests(parsed)
	if err != nil {
		return errdefs.PrependMsg(err, "create")
	}
	links, err := resources.FilterCreateLinkRequests(parsed)
	if err != nil {
		return errdefs.PrependMsg(err, "create")
	}

	if len(o.names) != 0 {
		seen := make(map[string]bool)
		for _, n := range o.names {
			seen[n] = false
		}
		assets = o.filterAssets(assets, seen)
		orchestrations = o.filterOrchestrations(orchestrations, seen)
		stages = o.filterStages(stages, seen)
		links = o.filterLinks(links, seen)
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
	execStub := pb.NewExecutionManagementClient(conn)

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
			return util.ErrorFromGrpcError(err)
		}
	}
	for _, req := range orchestrations {
		pbReq := &pb.CreateOrchestrationRequest{
			Name: string(req.Name),
		}
		if _, err = archStub.CreateOrchestration(ctx, pbReq); err != nil {
			return util.ErrorFromGrpcError(err)
		}
	}
	for _, req := range stages {
		pbReq := &pb.CreateStageRequest{
			Name:    string(req.Name),
			Asset:   string(req.Asset),
			Service: req.Service,
			Rpc:     req.Rpc,
			Address: req.Address,
			Host:    req.Host,
			Port:    req.Port,
		}
		if _, err = archStub.CreateStage(ctx, pbReq); err != nil {
			return util.ErrorFromGrpcError(err)
		}
	}
	for _, req := range links {
		l := &pb.CreateLinkRequest{
			Name:        string(req.Name),
			SourceStage: string(req.SourceStage),
			SourceField: req.SourceField,
			TargetStage: string(req.TargetStage),
			TargetField: req.TargetField,
		}
		if _, err = archStub.CreateLink(ctx, l); err != nil {
			return util.ErrorFromGrpcError(err)
		}
	}

	for _, req := range orchestrations {
		startReq := &pb.StartExecutionRequest{Orchestration: string(req.Name)}
		if _, err = execStub.Start(ctx, startReq); err != nil {
			return util.ErrorFromGrpcError(err)
		}
	}

	return nil
}

func (o *Options) filterAssets(
	requests []*api.CreateAssetRequest,
	seen map[string]bool,
) []*api.CreateAssetRequest {
	filtered := make([]*api.CreateAssetRequest, 0)
	for _, req := range requests {
		name := string(req.Name)
		_, exists := seen[name]
		if exists {
			filtered = append(filtered, req)
			seen[name] = true
		}
	}
	return filtered
}

func (o *Options) filterOrchestrations(
	requests []*api.CreateOrchestrationRequest,
	seen map[string]bool,
) []*api.CreateOrchestrationRequest {
	filtered := make([]*api.CreateOrchestrationRequest, 0)
	for _, req := range requests {
		name := string(req.Name)
		_, exists := seen[name]
		if exists {
			filtered = append(filtered, req)
			seen[name] = true
		}
	}
	return filtered
}

func (o *Options) filterStages(
	requests []*api.CreateStageRequest,
	seen map[string]bool,
) []*api.CreateStageRequest {
	filtered := make([]*api.CreateStageRequest, 0)
	for _, req := range requests {
		name := string(req.Name)
		_, exists := seen[name]
		if exists {
			filtered = append(filtered, req)
			seen[name] = true
		}
	}
	return filtered
}

func (o *Options) filterLinks(
	requests []*api.CreateLinkRequest,
	seen map[string]bool,
) []*api.CreateLinkRequest {
	filtered := make([]*api.CreateLinkRequest, 0)
	for _, req := range requests {
		name := string(req.Name)
		_, exists := seen[name]
		if exists {
			filtered = append(filtered, req)
			seen[name] = true
		}
	}
	return filtered
}
