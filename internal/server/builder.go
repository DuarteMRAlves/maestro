package server

import (
	pb "github.com/DuarteMRAlves/maestro/api/pb"
	ipb "github.com/DuarteMRAlves/maestro/internal/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/blueprint"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/link"
	"github.com/DuarteMRAlves/maestro/internal/stage"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Builder struct {
	grpcActive bool
	grpcOpts   []grpc.ServerOption

	logger *zap.Logger
}

func NewBuilder() *Builder {
	return &Builder{
		grpcActive: false,
	}
}

func (b *Builder) WithGrpc() *Builder {
	b.grpcActive = true
	return b
}

func (b *Builder) WithGrpcOpts(opts ...grpc.ServerOption) *Builder {
	b.grpcOpts = opts
	return b
}

func (b *Builder) WithLogger(logger *zap.Logger) *Builder {
	b.logger = logger
	return b
}

func (b *Builder) Build() (*Server, error) {
	err := b.complete()
	if err != nil {
		return nil, err
	}
	s := &Server{}
	initStores(s)
	if b.grpcActive {
		activateGrpc(s, b)
	}
	s.logger = b.logger
	return s, nil
}

// complete fills any values required to build the server with default options
func (b *Builder) complete() error {
	var err error
	if b.logger == nil {
		b.logger, err = zap.NewProduction()
		if err != nil {
			return errdefs.UnknownWithMsg("build: setup logger: %v", err)
		}
	}
	return nil
}

func initStores(s *Server) {
	s.assetStore = asset.NewStore()
	s.stageStore = stage.NewStore()
	s.linkStore = link.NewStore()
	s.blueprintStore = blueprint.NewStore()
}

func activateGrpc(s *Server, b *Builder) {
	grpcServer := grpc.NewServer(b.grpcOpts...)

	assetManagementServer := ipb.NewAssetManagementServer(s)
	stageManagementServer := ipb.NewStageManagementServer(s)
	linkManagementServer := ipb.NewLinkManagementServer(s)
	blueprintManagementServer := ipb.NewBlueprintManagementServer(s)

	pb.RegisterAssetManagementServer(grpcServer, assetManagementServer)
	pb.RegisterStageManagementServer(grpcServer, stageManagementServer)
	pb.RegisterLinkManagementServer(grpcServer, linkManagementServer)
	pb.RegisterBlueprintManagementServer(grpcServer, blueprintManagementServer)
	s.grpcServer = grpcServer
}
