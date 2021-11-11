package server

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"google.golang.org/grpc"
	"net"
)

const grpcNotConfigured = "grpc server not configured"

type Server struct {
	api        api.InternalAPI
	grpcServer *grpc.Server
}

func (s *Server) ServeGrpc(lis net.Listener) error {
	if ok, err := assert.Status(s.grpcServer != nil, grpcNotConfigured); !ok {
		return err
	}
	return s.grpcServer.Serve(lis)
}
