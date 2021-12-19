package mock

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	"google.golang.org/grpc"
)

// MaestroServer offers a simple aggregation of the available mock services
// to run in a single grpc server.
type MaestroServer struct {
	AssetManagementServer         pb.AssetManagementServer
	StageManagementServer         pb.StageManagementServer
	LinkManagementServer          pb.LinkManagementServer
	OrchestrationManagementServer pb.OrchestrationManagementServer
}

func (m *MaestroServer) GrpcServer() *grpc.Server {
	s := grpc.NewServer()

	if m.AssetManagementServer != nil {
		pb.RegisterAssetManagementServer(s, m.AssetManagementServer)
	}

	if m.StageManagementServer != nil {
		pb.RegisterStageManagementServer(s, m.StageManagementServer)
	}

	if m.LinkManagementServer != nil {
		pb.RegisterLinkManagementServer(s, m.LinkManagementServer)
	}

	if m.OrchestrationManagementServer != nil {
		pb.RegisterOrchestrationManagementServer(
			s,
			m.OrchestrationManagementServer)
	}

	return s
}
