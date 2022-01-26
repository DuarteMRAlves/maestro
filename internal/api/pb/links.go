package pb

import (
	"context"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/encoding/protobuff"
	"google.golang.org/protobuf/types/known/emptypb"
)

type linkManagementServer struct {
	pb.UnimplementedLinkManagementServer
	api api.InternalAPI
}

func NewLinkManagementServer(api api.InternalAPI) pb.LinkManagementServer {
	return &linkManagementServer{api: api}
}

func (s *linkManagementServer) Create(
	_ context.Context,
	pbReq *pb.CreateLinkRequest,
) (*emptypb.Empty, error) {

	var (
		req     api.CreateLinkRequest
		err     error
		grpcErr error = nil
	)

	protobuff.UnmarshalCreateLinkRequest(&req, pbReq)
	err = s.api.CreateLink(&req)
	if err != nil {
		grpcErr = GrpcErrorFromError(err)
	}
	return &emptypb.Empty{}, grpcErr
}

func (s *linkManagementServer) Get(
	pbReq *pb.GetLinkRequest,
	stream pb.LinkManagement_GetServer,
) error {

	var (
		req api.GetLinkRequest
		err error
	)

	protobuff.UnmarshalGetLinkRequest(&req, pbReq)
	links, err := s.api.GetLink(&req)
	if err != nil {
		return GrpcErrorFromError(err)
	}
	for _, l := range links {
		pbLink, err := protobuff.MarshalLink(l)
		if err != nil {
			return err
		}
		stream.Send(pbLink)
	}
	return nil
}
