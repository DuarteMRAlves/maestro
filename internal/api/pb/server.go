package pb

import (
	"context"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/events"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
)

type CreateAsset func(*api.CreateAssetRequest) error
type GetAssets func(*api.GetAssetRequest) ([]*api.Asset, error)

type CreateStage func(*api.CreateStageRequest) error
type GetStages func(*api.GetStageRequest) ([]*api.Stage, error)

type CreateLink func(*api.CreateLinkRequest) error
type GetLinks func(*api.GetLinkRequest) ([]*api.Link, error)

type CreateOrchestration func(*api.CreateOrchestrationRequest) error
type GetOrchestrations func(*api.GetOrchestrationRequest) (
	[]*api.Orchestration,
	error,
)

type StartExecution func(*api.StartExecutionRequest) error
type AttachExecution func(*api.AttachExecutionRequest) (
	*events.Subscription,
	error,
)

type ServerManagement struct {
	CreateAsset CreateAsset
	GetAssets   GetAssets

	CreateOrchestration CreateOrchestration
	GetOrchestrations   GetOrchestrations

	CreateStage CreateStage
	GetStages   GetStages

	CreateLink CreateLink
	GetLinks   GetLinks

	StartExecution  StartExecution
	AttachExecution AttachExecution
}

// RegisterServices registers all grpc services into a grpc server.
func RegisterServices(s *grpc.Server, m ServerManagement) {
	archServ := &architectureService{
		createAsset:         m.CreateAsset,
		getAssets:           m.GetAssets,
		createOrchestration: m.CreateOrchestration,
		getOrchestrations:   m.GetOrchestrations,
		createStage:         m.CreateStage,
		getStages:           m.GetStages,
		createLink:          m.CreateLink,
		getLinks:            m.GetLinks,
	}
	pb.RegisterArchitectureManagementServer(s, archServ)
	execServ := &executionsService{
		startExecution:  m.StartExecution,
		attachExecution: m.AttachExecution,
	}
	pb.RegisterExecutionManagementServer(s, execServ)
}

type architectureService struct {
	pb.UnimplementedArchitectureManagementServer

	createAsset CreateAsset
	getAssets   GetAssets

	createOrchestration CreateOrchestration
	getOrchestrations   GetOrchestrations

	createStage CreateStage
	getStages   GetStages

	createLink CreateLink
	getLinks   GetLinks
}

func (s *architectureService) CreateAsset(
	_ context.Context,
	pbReq *pb.CreateAssetRequest,
) (*emptypb.Empty, error) {

	var (
		req api.CreateAssetRequest
	)
	var err error
	var grpcErr error = nil

	UnmarshalCreateAssetRequest(&req, pbReq)
	err = s.createAsset(&req)
	if err != nil {
		grpcErr = GrpcErrorFromError(err)
	}
	return &emptypb.Empty{}, grpcErr
}

func (s *architectureService) GetAsset(
	pbQuery *pb.GetAssetRequest,
	stream pb.ArchitectureManagement_GetAssetServer,
) error {

	var (
		query api.GetAssetRequest
		err   error
	)

	UnmarshalGetAssetRequest(&query, pbQuery)

	assets, err := s.getAssets(&query)
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

func (s *architectureService) CreateOrchestration(
	_ context.Context,
	pbReq *pb.CreateOrchestrationRequest,
) (*emptypb.Empty, error) {

	var (
		req     api.CreateOrchestrationRequest
		err     error
		grpcErr error = nil
	)

	UnmarshalCreateOrchestrationRequest(&req, pbReq)
	err = s.createOrchestration(&req)
	if err != nil {
		grpcErr = GrpcErrorFromError(err)
	}
	return &emptypb.Empty{}, grpcErr
}

func (s *architectureService) GetOrchestration(
	pbReq *pb.GetOrchestrationRequest,
	stream pb.ArchitectureManagement_GetOrchestrationServer,
) error {

	var (
		req api.GetOrchestrationRequest
		err error
	)

	UnmarshalGetOrchestrationRequest(&req, pbReq)
	orchestrations, err := s.getOrchestrations(&req)
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

func (s *architectureService) CreateStage(
	_ context.Context,
	pbReq *pb.CreateStageRequest,
) (*emptypb.Empty, error) {

	var req api.CreateStageRequest
	var err error
	var grpcErr error = nil

	UnmarshalCreateStageRequest(&req, pbReq)
	err = s.createStage(&req)
	if err != nil {
		grpcErr = GrpcErrorFromError(err)
	}
	return &emptypb.Empty{}, grpcErr
}

func (s *architectureService) GetStage(
	pbReq *pb.GetStageRequest,
	stream pb.ArchitectureManagement_GetStageServer,
) error {

	var (
		req api.GetStageRequest
		err error
	)

	UnmarshalGetStageRequest(&req, pbReq)
	stages, err := s.getStages(&req)
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

func (s *architectureService) CreateLink(
	_ context.Context,
	pbReq *pb.CreateLinkRequest,
) (*emptypb.Empty, error) {

	var (
		req     api.CreateLinkRequest
		err     error
		grpcErr error = nil
	)

	UnmarshalCreateLinkRequest(&req, pbReq)
	err = s.createLink(&req)
	if err != nil {
		grpcErr = GrpcErrorFromError(err)
	}
	return &emptypb.Empty{}, grpcErr
}

func (s *architectureService) GetLink(
	pbReq *pb.GetLinkRequest,
	stream pb.ArchitectureManagement_GetLinkServer,
) error {

	var (
		req api.GetLinkRequest
		err error
	)

	UnmarshalGetLinkRequest(&req, pbReq)
	links, err := s.getLinks(&req)
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

type executionsService struct {
	pb.UnimplementedExecutionManagementServer

	startExecution  StartExecution
	attachExecution AttachExecution
}

func (s *executionsService) Start(
	_ context.Context,
	pbReq *pb.StartExecutionRequest,
) (*emptypb.Empty, error) {
	var (
		req api.StartExecutionRequest
		err error
	)

	UnmarshalStartExecutionRequest(&req, pbReq)
	err = s.startExecution(&req)
	if err != nil {
		return nil, GrpcErrorFromError(err)
	}
	return &emptypb.Empty{}, nil
}

func (s *executionsService) Attach(stream pb.ExecutionManagement_AttachServer) error {
	received := make(chan *api.AttachExecutionRequest)
	errs1 := make(chan error)
	errs2 := make(chan error)
	errs := make(chan error)
	go func() {
		for {
			pbReq, err := stream.Recv()
			if err != nil {
				if err != io.EOF {
					errs1 <- err
				}
				close(received)
				close(errs1)
				return
			}
			req := &api.AttachExecutionRequest{}
			UnmarshalAttachExecutionRequest(req, pbReq)
			received <- req
		}
	}()
	go func() {
		req, open := <-received
		if !open {
			close(errs2)
			return
		}
		sub, err := s.attachExecution(req)
		if err != nil {
			errs2 <- GrpcErrorFromError(err)
			close(errs2)
			return
		}
		for _, event := range sub.Hist {
			pbEvent := &pb.Event{}
			MarshalEvent(pbEvent, event)
			err = stream.Send(pbEvent)
			if err != nil {
				errs2 <- err
				close(errs2)
				return
			}
		}
		for {
			select {
			case <-received:
				// TODO: Unsubscribe
				close(errs2)
				return
			case event := <-sub.Future:
				pbEvent := &pb.Event{}
				MarshalEvent(pbEvent, event)
				err = stream.Send(pbEvent)
				if err != nil {
					errs2 <- err
					close(errs2)
					return
				}
			}
		}
	}()
	go func() {
		defer close(errs)
		for errs1 != nil && errs2 != nil {
			select {
			case err, open := <-errs1:
				if !open {
					errs1 = nil
					continue
				}
				errs <- err
			case err, open := <-errs2:
				if !open {
					errs2 = nil
					continue
				}
				errs <- err
			}
		}
	}()
	err, open := <-errs
	if !open {
		return nil
	}
	return err
}
