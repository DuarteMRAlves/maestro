package server

import (
	pb "github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
	ipb "github.com/DuarteMRAlves/maestro/internal/api/pb"
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
	s := &Server{api: api.NewInternalAPI()}
	if b.grpcActive {
		activateGrpc(s, b)
	}
	return s
}

func activateGrpc(s *Server, b *Builder) {
	grpcServer := grpc.NewServer(b.grpcOpts...)

	assetManagementServer := ipb.NewAssetManagementServer(s.api)
	blueprintManagementServer := ipb.NewBlueprintManagementServer(s.api)
	
	pb.RegisterAssetManagementServer(grpcServer, assetManagementServer)
	pb.RegisterBlueprintManagementServer(grpcServer, blueprintManagementServer)
	s.grpcServer = grpcServer
}
