package pb

import (
	"context"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
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
	pbLink *pb.Link,
) (*emptypb.Empty, error) {

	var link *apitypes.Link
	var err error
	var grpcErr error = nil

	if link, err = protobuff.UnmarshalLink(pbLink); err != nil {
		return &emptypb.Empty{}, GrpcErrorFromError(err)
	}
	err = s.api.CreateLink(link)
	if err != nil {
		grpcErr = GrpcErrorFromError(err)
	}
	return &emptypb.Empty{}, grpcErr
}

func (s *linkManagementServer) Get(
	pbQuery *pb.Link,
	stream pb.LinkManagement_GetServer,
) error {

	var (
		query *apitypes.Link
		err   error
	)

	if query, err = protobuff.UnmarshalLink(pbQuery); err != nil {
		return GrpcErrorFromError(err)
	}

	links, err := s.api.GetLink(query)
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
