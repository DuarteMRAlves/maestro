package pb

import (
	"context"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// MockMaestroServer offers a simple aggregation of the available mock services
// to run in a single grpc server.
type MockMaestroServer struct {
	ArchitectureManagementServer pb.ArchitectureManagementServer
	ExecutionManagementServer    pb.ExecutionManagementServer
}

func (m *MockMaestroServer) GrpcServer() *grpc.Server {
	s := grpc.NewServer()

	if m.ArchitectureManagementServer != nil {
		pb.RegisterArchitectureManagementServer(
			s,
			m.ArchitectureManagementServer,
		)
	}

	if m.ExecutionManagementServer != nil {
		pb.RegisterExecutionManagementServer(s, m.ExecutionManagementServer)
	}

	return s
}

// MockArchitectureManagementServer is a mocking of an architecture management
// server to be used during tests. By default, methods are not implemented and
// raise an error.
//
// The server can be initialized with a specific function for each grpc method
// that will be called.
type MockArchitectureManagementServer struct {
	pb.UnimplementedArchitectureManagementServer

	CreateAssetFn func(
		context.Context,
		*pb.CreateAssetRequest,
	) (*emptypb.Empty, error)
	GetAssetFn func(
		*pb.GetAssetRequest,
		pb.ArchitectureManagement_GetAssetServer,
	) error

	CreateOrchestrationFn func(
		context.Context,
		*pb.CreateOrchestrationRequest,
	) (*emptypb.Empty, error)
	GetOrchestrationFn func(
		*pb.GetOrchestrationRequest,
		pb.ArchitectureManagement_GetOrchestrationServer,
	) error

	CreateStageFn func(context.Context, *pb.CreateStageRequest) (
		*emptypb.Empty,
		error,
	)
	GetStageFn func(
		*pb.GetStageRequest,
		pb.ArchitectureManagement_GetStageServer,
	) error

	CreateLinkFn func(context.Context, *pb.CreateLinkRequest) (
		*emptypb.Empty,
		error,
	)
	GetLinkFn func(
		*pb.GetLinkRequest,
		pb.ArchitectureManagement_GetLinkServer,
	) error
}

func (s *MockArchitectureManagementServer) CreateAsset(
	ctx context.Context,
	req *pb.CreateAssetRequest,
) (*emptypb.Empty, error) {
	if s.CreateAssetFn != nil {
		return s.CreateAssetFn(ctx, req)
	}
	return &emptypb.Empty{}, fmt.Errorf(
		"method CreateAsset not configured but called with request %v",
		req,
	)
}

func (s *MockArchitectureManagementServer) GetAsset(
	req *pb.GetAssetRequest,
	stream pb.ArchitectureManagement_GetAssetServer,
) error {
	if s.GetAssetFn != nil {
		return s.GetAssetFn(req, stream)
	}
	return fmt.Errorf(
		"method GetAsset not configured but called with request %v",
		req,
	)
}

func (s *MockArchitectureManagementServer) CreateOrchestration(
	ctx context.Context,
	req *pb.CreateOrchestrationRequest,
) (*emptypb.Empty, error) {
	if s.CreateOrchestrationFn != nil {
		return s.CreateOrchestrationFn(ctx, req)
	}
	return &emptypb.Empty{}, fmt.Errorf(
		"method CreateOrchestration not configured but called with req %v",
		req,
	)
}

func (s *MockArchitectureManagementServer) GetOrchestration(
	req *pb.GetOrchestrationRequest,
	stream pb.ArchitectureManagement_GetOrchestrationServer,
) error {
	if s.GetOrchestrationFn != nil {
		return s.GetOrchestrationFn(req, stream)
	}
	return fmt.Errorf(
		"method GetOrchestration not configured but called with req %v",
		req,
	)
}

func (s *MockArchitectureManagementServer) CreateStage(
	ctx context.Context,
	req *pb.CreateStageRequest,
) (*emptypb.Empty, error) {
	if s.CreateStageFn != nil {
		return s.CreateStageFn(ctx, req)
	}
	return &emptypb.Empty{}, fmt.Errorf(
		"method CreateStage not configured but called with req %v",
		req,
	)
}

func (s *MockArchitectureManagementServer) GetStage(
	req *pb.GetStageRequest,
	stream pb.ArchitectureManagement_GetStageServer,
) error {
	if s.GetStageFn != nil {
		return s.GetStageFn(req, stream)
	}
	return fmt.Errorf(
		"method GetStage not configured but called with req %v",
		req,
	)
}

func (s *MockArchitectureManagementServer) CreateLink(
	ctx context.Context,
	req *pb.CreateLinkRequest,
) (*emptypb.Empty, error) {
	if s.CreateLinkFn != nil {
		return s.CreateLinkFn(ctx, req)
	}
	return &emptypb.Empty{}, fmt.Errorf(
		"method CreateLink not configured but called with req %v",
		req,
	)
}

func (s *MockArchitectureManagementServer) GetLink(
	req *pb.GetLinkRequest,
	stream pb.ArchitectureManagement_GetLinkServer,
) error {
	if s.GetLinkFn != nil {
		return s.GetLinkFn(req, stream)
	}
	return fmt.Errorf(
		"method GetLink not configured but called with req %v",
		req,
	)
}

// MockExecutionManagementServer is a mocking of an execution management
// server to be used during tests. By default, methods are not implemented and
// raise an error.
//
// The server can be initialized with a specific function for each grpc method
// that will be called.
type MockExecutionManagementServer struct {
	pb.UnimplementedExecutionManagementServer
	StartExecutionFn func(
		context.Context,
		*pb.StartExecutionRequest,
	) (*emptypb.Empty, error)
}

func (s *MockExecutionManagementServer) Start(
	ctx context.Context,
	req *pb.StartExecutionRequest,
) (*emptypb.Empty, error) {
	if s.StartExecutionFn != nil {
		return s.StartExecutionFn(ctx, req)
	}
	return &emptypb.Empty{}, fmt.Errorf(
		"method StartExecution not configured but called with req %v",
		req,
	)
}
