package pb

import (
	"context"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// RegisterServices registers all grpc services into a grpc server.
func RegisterServices(s *grpc.Server, api api.InternalAPI) {
	pb.RegisterAssetManagementServer(s, &assetsService{api: api})
	pb.RegisterStageManagementServer(s, &stagesService{api: api})
	pb.RegisterLinkManagementServer(s, &linksService{api: api})
	pb.RegisterOrchestrationManagementServer(
		s,
		&orchestrationsService{api: api},
	)
}

type assetsService struct {
	pb.UnimplementedAssetManagementServer
	api api.InternalAPI
}

func (s *assetsService) Create(
	_ context.Context,
	pbReq *pb.CreateAssetRequest,
) (*emptypb.Empty, error) {

	var (
		req api.CreateAssetRequest
	)
	var err error
	var grpcErr error = nil

	UnmarshalCreateAssetRequest(&req, pbReq)
	err = s.api.CreateAsset(&req)
	if err != nil {
		grpcErr = GrpcErrorFromError(err)
	}
	return &emptypb.Empty{}, grpcErr
}

func (s *assetsService) Get(
	pbQuery *pb.GetAssetRequest,
	stream pb.AssetManagement_GetServer,
) error {

	var (
		query api.GetAssetRequest
		err   error
	)

	UnmarshalGetAssetRequest(&query, pbQuery)

	assets, err := s.api.GetAsset(&query)
	if err != nil {
		return GrpcErrorFromError(err)
	}
	for _, a := range assets {
		pbAsset, err := MarshalAsset(a)
		if err != nil {
			return err
		}
		stream.Send(pbAsset)
	}
	return nil
}

type orchestrationsService struct {
	pb.UnimplementedOrchestrationManagementServer
	api api.InternalAPI
}

func (s *orchestrationsService) Create(
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

func (s *orchestrationsService) Get(
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

type stagesService struct {
	pb.UnimplementedStageManagementServer
	api api.InternalAPI
}

func (s *stagesService) Create(
	_ context.Context,
	pbReq *pb.CreateStageRequest,
) (*emptypb.Empty, error) {

	var req api.CreateStageRequest
	var err error
	var grpcErr error = nil

	UnmarshalCreateStageRequest(&req, pbReq)
	err = s.api.CreateStage(&req)
	if err != nil {
		grpcErr = GrpcErrorFromError(err)
	}
	return &emptypb.Empty{}, grpcErr
}

func (s *stagesService) Get(
	pbReq *pb.GetStageRequest,
	stream pb.StageManagement_GetServer,
) error {

	var (
		req api.GetStageRequest
		err error
	)

	UnmarshalGetStageRequest(&req, pbReq)
	stages, err := s.api.GetStage(&req)
	if err != nil {
		return GrpcErrorFromError(err)
	}
	for _, s := range stages {
		pbStage, err := MarshalStage(s)
		if err != nil {
			return err
		}
		stream.Send(pbStage)
	}
	return nil
}

type linksService struct {
	pb.UnimplementedLinkManagementServer
	api api.InternalAPI
}

func (s *linksService) Create(
	_ context.Context,
	pbReq *pb.CreateLinkRequest,
) (*emptypb.Empty, error) {

	var (
		req     api.CreateLinkRequest
		err     error
		grpcErr error = nil
	)

	UnmarshalCreateLinkRequest(&req, pbReq)
	err = s.api.CreateLink(&req)
	if err != nil {
		grpcErr = GrpcErrorFromError(err)
	}
	return &emptypb.Empty{}, grpcErr
}

func (s *linksService) Get(
	pbReq *pb.GetLinkRequest,
	stream pb.LinkManagement_GetServer,
) error {

	var (
		req api.GetLinkRequest
		err error
	)

	UnmarshalGetLinkRequest(&req, pbReq)
	links, err := s.api.GetLink(&req)
	if err != nil {
		return GrpcErrorFromError(err)
	}
	for _, l := range links {
		pbLink, err := MarshalLink(l)
		if err != nil {
			return err
		}
		stream.Send(pbLink)
	}
	return nil
}
