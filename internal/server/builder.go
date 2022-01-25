package server

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	ipb "github.com/DuarteMRAlves/maestro/internal/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/execution"
	"github.com/DuarteMRAlves/maestro/internal/reflection"
	"github.com/DuarteMRAlves/maestro/internal/storage"
	"github.com/dgraph-io/badger/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Builder struct {
	reflectionManager reflection.Manager

	grpcActive bool
	grpcOpts   []grpc.ServerOption

	// db is a key-value store database to persist state across multiple
	// executions of the server and to ensure consistency.
	db *badger.DB

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

func (b *Builder) WithDb(db *badger.DB) *Builder {
	b.db = db
	return b
}

func (b *Builder) WithLogger(logger *zap.Logger) *Builder {
	b.logger = logger
	return b
}

func (b *Builder) WithReflectionManager(m reflection.Manager) *Builder {
	b.reflectionManager = m
	return b
}

func (b *Builder) Build() (*Server, error) {
	var err error
	err = b.complete()
	if err != nil {
		return nil, err
	}
	err = b.validate()
	if err != nil {
		return nil, err
	}
	s := &Server{}
	s.logger = b.logger
	s.db = b.db
	b.initManagers(s)
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
	if b.reflectionManager == nil {
		b.reflectionManager = reflection.NewManager()
	}
	return nil
}

// validate checks if all necessary preconditions to build the server are
// fulfilled.
func (b *Builder) validate() error {
	if b.db == nil {
		return errdefs.FailedPreconditionWithMsg("no database specified")
	}
	return nil
}

func (b *Builder) initManagers(s *Server) {
	s.reflectionManager = b.reflectionManager
	s.storageManager = storage.NewManager(s.reflectionManager)
	s.flowManager = execution.NewManager(s.reflectionManager)
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
		orchestrationManagementServer,
	)
	s.grpcServer = grpcServer
}
