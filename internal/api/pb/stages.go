package pb

import (
	"context"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
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
	pbStage *pb.Stage,
) (*emptypb.Empty, error) {

	var stage *apitypes.Stage
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

func (s *stageManagementServer) Get(
	pbQuery *pb.Stage,
	stream pb.StageManagement_GetServer,
) error {

	var (
		query *apitypes.Stage
		err   error
	)

	if query, err = protobuff.UnmarshalStage(pbQuery); err != nil {
		return GrpcErrorFromError(err)
	}

	stages := s.api.GetStage(query)
	for _, s := range stages {
		pbStage, err := protobuff.MarshalStage(s)
		if err != nil {
			return err
		}
		stream.Send(pbStage)
	}
	return nil
}
