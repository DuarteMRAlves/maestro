package pb

import (
	"context"
	"errors"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/blueprint"
	"github.com/DuarteMRAlves/maestro/internal/encoding/protobuff"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		grpcErr error
	)

	if bp, err = protobuff.UnmarshalBlueprint(pbBlueprint); err != nil {
		return &emptypb.Empty{}, status.Error(codes.Unknown, err.Error())
	}
	err = s.api.CreateBlueprint(bp)
	if err != nil {
		var alreadyExists blueprint.AlreadyExists
		if errors.As(err, &alreadyExists) {
			grpcErr = status.Error(codes.AlreadyExists, alreadyExists.Error())
		} else {
			grpcErr = status.Error(codes.Unknown, err.Error())
		}
	}
	return &emptypb.Empty{}, grpcErr
}
