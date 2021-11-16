package pb

import (
	"context"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/blueprint"
	"github.com/DuarteMRAlves/maestro/internal/encoding/protobuff"
	"github.com/DuarteMRAlves/maestro/internal/identifier"
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
) (*pb.Id, error) {

	var (
		bp  *blueprint.Blueprint
		id  identifier.Id
		err error
	)

	if bp, err = protobuff.UnmarshalBlueprint(pbBlueprint); err != nil {
		return emptyIdPb, err
	}
	if id, err = s.api.CreateBlueprint(bp); err != nil {
		return emptyIdPb, err
	}
	return protobuff.MarshalID(id), nil
}
