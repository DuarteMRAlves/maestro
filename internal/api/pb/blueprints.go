package pb

import (
	"context"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/blueprint"
	"github.com/DuarteMRAlves/maestro/internal/encoding/protobuff"
	"google.golang.org/protobuf/types/known/emptypb"
)

type blueprintManagementServer struct {
	pb.UnimplementedBlueprintManagementServer
	api api.InternalAPI
}

func NewBlueprintManagementServer(
	api api.InternalAPI,
) pb.BlueprintManagementServer {
	return &blueprintManagementServer{api: api}
}

func (s *blueprintManagementServer) Create(
	_ context.Context,
	pbBlueprint *pb.Blueprint,
) (*emptypb.Empty, error) {

	var (
		bp      *blueprint.Blueprint
		err     error
		grpcErr error = nil
	)

	if bp, err = protobuff.UnmarshalBlueprint(pbBlueprint); err != nil {
		return &emptypb.Empty{}, GrpcErrorFromError(err)
	}
	err = s.api.CreateBlueprint(bp)
	if err != nil {
		grpcErr = GrpcErrorFromError(err)
	}
	return &emptypb.Empty{}, grpcErr
}

func (s *blueprintManagementServer) Get(
	pbQuery *pb.Blueprint,
	stream pb.BlueprintManagement_GetServer,
) error {

	var query *blueprint.Blueprint
	var err error

	if query, err = protobuff.UnmarshalBlueprint(pbQuery); err != nil {
		return err
	}

	blueprints := s.api.GetBlueprint(query)
	for _, a := range blueprints {
		pbBlueprint, err := protobuff.MarshalBlueprint(a)
		if err != nil {
			return err
		}
		err = stream.Send(pbBlueprint)
		if err != nil {
			return err
		}
	}
	return nil
}
