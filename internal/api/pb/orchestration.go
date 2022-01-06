package pb

import (
	"context"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/encoding/protobuff"
	"google.golang.org/protobuf/types/known/emptypb"
)

type orchestrationManagementServer struct {
	pb.UnimplementedOrchestrationManagementServer
	api api.InternalAPI
}

func NewOrchestrationManagementServer(
	api api.InternalAPI,
) pb.OrchestrationManagementServer {
	return &orchestrationManagementServer{api: api}
}

func (s *orchestrationManagementServer) Create(
	_ context.Context,
	pbOrchestration *pb.Orchestration,
) (*emptypb.Empty, error) {

	var (
		o       *apitypes.Orchestration
		err     error
		grpcErr error = nil
	)

	if o, err = protobuff.UnmarshalOrchestration(pbOrchestration); err != nil {
		return &emptypb.Empty{}, GrpcErrorFromError(err)
	}
	err = s.api.CreateOrchestration(o)
	if err != nil {
		grpcErr = GrpcErrorFromError(err)
	}
	return &emptypb.Empty{}, grpcErr
}

func (s *orchestrationManagementServer) Get(
	pbQuery *pb.Orchestration,
	stream pb.OrchestrationManagement_GetServer,
) error {

	var query *apitypes.Orchestration
	var err error

	if query, err = protobuff.UnmarshalOrchestration(pbQuery); err != nil {
		return err
	}

	orchestrations := s.api.GetOrchestration(query)
	for _, a := range orchestrations {
		pbOrchestration, err := protobuff.MarshalOrchestration(a)
		if err != nil {
			return err
		}
		err = stream.Send(pbOrchestration)
		if err != nil {
			return err
		}
	}
	return nil
}
