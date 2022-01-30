package pb

import (
	"context"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
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
	pbReq *pb.CreateOrchestrationRequest,
) (*emptypb.Empty, error) {

	var (
		req     api.CreateOrchestrationRequest
		err     error
		grpcErr error = nil
	)

	UnmarshalCreateOrchestrationRequest(&req, pbReq)
	err = s.api.CreateOrchestration(&req)
	if err != nil {
		grpcErr = GrpcErrorFromError(err)
	}
	return &emptypb.Empty{}, grpcErr
}

func (s *orchestrationManagementServer) Get(
	pbReq *pb.GetOrchestrationRequest,
	stream pb.OrchestrationManagement_GetServer,
) error {

	var (
		req api.GetOrchestrationRequest
		err error
	)

	UnmarshalGetOrchestrationRequest(&req, pbReq)
	orchestrations, err := s.api.GetOrchestration(&req)
	if err != nil {
		return GrpcErrorFromError(err)
	}
	for _, a := range orchestrations {
		pbOrchestration, err := MarshalOrchestration(a)
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
