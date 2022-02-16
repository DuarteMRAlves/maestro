package pb

import (
	"context"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
)

// RegisterServices registers all grpc services into a grpc server.
func RegisterServices(s *grpc.Server, api api.InternalAPI) {
	pb.RegisterArchitectureManagementServer(s, &architectureService{api: api})
	pb.RegisterExecutionManagementServer(s, &executionsService{api: api})
}

type architectureService struct {
	pb.UnimplementedArchitectureManagementServer
	api api.InternalAPI
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
	err = s.api.CreateAsset(&req)
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
	err = s.api.CreateOrchestration(&req)
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

func (s *architectureService) CreateStage(
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

func (s *architectureService) GetStage(
	pbReq *pb.GetStageRequest,
	stream pb.ArchitectureManagement_GetStageServer,
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
	err = s.api.CreateLink(&req)
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

type executionsService struct {
	pb.UnimplementedExecutionManagementServer
	api api.InternalAPI
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
	err = s.api.StartExecution(&req)
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
		sub, err := s.api.AttachExecution(req)
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
