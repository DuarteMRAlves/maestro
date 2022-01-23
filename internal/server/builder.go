package server

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	ipb "github.com/DuarteMRAlves/maestro/internal/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/flow"
	"github.com/DuarteMRAlves/maestro/internal/orchestration"
	"github.com/dgraph-io/badger/v3"
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
	s.logger = b.logger
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	if err != nil {
		return nil, errdefs.UnknownWithMsg("initialize db: %v", err)
	}
	s.db = db
	initStores(s)
	initManagers(s)
	if b.grpcActive {
		activateGrpc(s, b)
	}
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
}

func initManagers(s *Server) {
	s.orchestrationManager = orchestration.NewManager(s.assetStore)
	s.flowManager = flow.NewManager()
}

func activateGrpc(s *Server, b *Builder) {
	grpcServer := grpc.NewServer(b.grpcOpts...)

	assetManagementServer := ipb.NewAssetManagementServer(s)
	stageManagementServer := ipb.NewStageManagementServer(s)
	linkManagementServer := ipb.NewLinkManagementServer(s)
	orchestrationManagementServer := ipb.NewOrchestrationManagementServer(s)

	pb.RegisterAssetManagementServer(grpcServer, assetManagementServer)
	pb.RegisterStageManagementServer(grpcServer, stageManagementServer)
	pb.RegisterLinkManagementServer(grpcServer, linkManagementServer)
	pb.RegisterOrchestrationManagementServer(
		grpcServer,
		orchestrationManagementServer)
	s.grpcServer = grpcServer
}
