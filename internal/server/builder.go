package server

import (
	pb "github.com/DuarteMRAlves/maestro/api/pb"
	ipb "github.com/DuarteMRAlves/maestro/internal/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/blueprint"
	"github.com/DuarteMRAlves/maestro/internal/stage"
	"google.golang.org/grpc"
)

type Builder struct {
	grpcActive bool
	grpcOpts   []grpc.ServerOption
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

func (b *Builder) Build() *Server {
	s := &Server{}
	initStores(s)
	if b.grpcActive {
		activateGrpc(s, b)
	}
	return s
}

func initStores(s *Server) {
	s.assetStore = asset.NewStore()
	s.stageStore = stage.NewStore()
	s.blueprintStore = blueprint.NewStore()
}

func activateGrpc(s *Server, b *Builder) {
	grpcServer := grpc.NewServer(b.grpcOpts...)

	assetManagementServer := ipb.NewAssetManagementServer(s)
	blueprintManagementServer := ipb.NewBlueprintManagementServer(s)

	pb.RegisterAssetManagementServer(grpcServer, assetManagementServer)
	pb.RegisterBlueprintManagementServer(grpcServer, blueprintManagementServer)
	s.grpcServer = grpcServer
}
