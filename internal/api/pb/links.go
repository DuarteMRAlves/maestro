package pb

import (
	"context"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/encoding/protobuff"
	"github.com/DuarteMRAlves/maestro/internal/link"
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

	var link *link.Link
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
		query *link.Link
		err   error
	)

	if query, err = protobuff.UnmarshalLink(pbQuery); err != nil {
		return GrpcErrorFromError(err)
	}

	links := s.api.GetLink(query)
	for _, l := range links {
		pbLink, err := protobuff.MarshalLink(l)
		if err != nil {
			return err
		}
		stream.Send(pbLink)
	}
	return nil
}
