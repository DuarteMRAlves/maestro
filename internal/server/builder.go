package server

import (
	apipb "github.com/DuarteMRAlves/maestro/internal/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/orchestration"
	"github.com/dgraph-io/badger/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Builder struct {
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
	s.execManager = orchestration.NewManager(b.logger)
}

func activateGrpc(s *Server, b *Builder) {
	grpcServer := grpc.NewServer(b.grpcOpts...)
	m := apipb.ServerManagement{
		CreateAsset:         s.CreateAsset,
		GetAssets:           s.GetAsset,
		CreateOrchestration: s.CreateOrchestration,
		GetOrchestrations:   s.GetOrchestration,
		CreateStage:         s.CreateStage,
		GetStages:           s.GetStage,
		CreateLink:          s.CreateLink,
		GetLinks:            s.GetLink,
		StartExecution:      s.StartExecution,
	}
	apipb.RegisterServices(grpcServer, m)
	s.grpcServer = grpcServer
}
