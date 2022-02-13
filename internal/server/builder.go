package server

import (
	apipb "github.com/DuarteMRAlves/maestro/internal/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/arch"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/exec"
	"github.com/DuarteMRAlves/maestro/internal/rpc"
	"github.com/dgraph-io/badger/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Builder struct {
	reflectionManager rpc.Manager

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

func (b *Builder) WithReflectionManager(m rpc.Manager) *Builder {
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
	err = b.initManagers(s)
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
		b.reflectionManager = rpc.NewManager()
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

func (b *Builder) initManagers(s *Server) error {
	var err error
	s.reflectionManager = b.reflectionManager
	storageManagerCtx := arch.NewDefaultContext(s.db, s.reflectionManager)
	s.archManager, err = arch.NewManager(storageManagerCtx)
	if err != nil {
		return errdefs.PrependMsg(err, "init managers:")
	}
	s.execManager = exec.NewManager(s.reflectionManager)
	return nil
}

func activateGrpc(s *Server, b *Builder) {
	grpcServer := grpc.NewServer(b.grpcOpts...)
	apipb.RegisterServices(grpcServer, s)
	s.grpcServer = grpcServer
}
