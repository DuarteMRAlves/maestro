package pb

import (
	"context"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/encoding/protobuff"
	"github.com/DuarteMRAlves/maestro/internal/stage"
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
	pbStage *pb.Stage,
) (*emptypb.Empty, error) {

	var stage *stage.Stage
	var err error
	var grpcErr error = nil

	if stage, err = protobuff.UnmarshalStage(pbStage); err != nil {
		return &emptypb.Empty{}, GrpcErrorFromError(err)
	}
	err = s.api.CreateStage(stage)
	if err != nil {
		grpcErr = GrpcErrorFromError(err)
	}
	return &emptypb.Empty{}, grpcErr
}
