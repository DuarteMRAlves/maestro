package pb

import (
	"context"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/encoding/protobuff"
	"google.golang.org/protobuf/types/known/emptypb"
)

type stageManagementServer struct {
	pb.UnimplementedStageManagementServer
	api api.InternalAPI
}

func NewStageManagementServer(api api.InternalAPI) pb.StageManagementServer {
	return &stageManagementServer{api: api}
}

func (s *stageManagementServer) Create(
	_ context.Context,
	pbReq *pb.CreateStageRequest,
) (*emptypb.Empty, error) {

	var req api.CreateStageRequest
	var err error
	var grpcErr error = nil

	protobuff.UnmarshalCreateStageRequest(&req, pbReq)
	err = s.api.CreateStage(&req)
	if err != nil {
		grpcErr = GrpcErrorFromError(err)
	}
	return &emptypb.Empty{}, grpcErr
}

func (s *stageManagementServer) Get(
	pbReq *pb.GetStageRequest,
	stream pb.StageManagement_GetServer,
) error {

	var (
		req api.GetStageRequest
		err error
	)

	protobuff.UnmarshalGetStageRequest(&req, pbReq)
	stages, err := s.api.GetStage(&req)
	if err != nil {
		return GrpcErrorFromError(err)
	}
	for _, s := range stages {
		pbStage, err := protobuff.MarshalStage(s)
		if err != nil {
			return err
		}
		stream.Send(pbStage)
	}
	return nil
}
