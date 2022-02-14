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
	AssetManagementServer         pb.AssetManagementServer
	StageManagementServer         pb.StageManagementServer
	LinkManagementServer          pb.LinkManagementServer
	OrchestrationManagementServer pb.OrchestrationManagementServer
	ExecutionManagementServer     pb.ExecutionManagementServer
}

func (m *MockMaestroServer) GrpcServer() *grpc.Server {
	s := grpc.NewServer()

	if m.AssetManagementServer != nil {
		pb.RegisterAssetManagementServer(s, m.AssetManagementServer)
	}

	if m.StageManagementServer != nil {
		pb.RegisterStageManagementServer(s, m.StageManagementServer)
	}

	if m.LinkManagementServer != nil {
		pb.RegisterLinkManagementServer(s, m.LinkManagementServer)
	}

	if m.OrchestrationManagementServer != nil {
		pb.RegisterOrchestrationManagementServer(
			s,
			m.OrchestrationManagementServer,
		)
	}

	if m.ExecutionManagementServer != nil {
		pb.RegisterExecutionManagementServer(s, m.ExecutionManagementServer)
	}

	return s
}

// MockOrchestrationManagementServer is a mocking of an orchestration management
// server to be used during tests. By default, methods are not implemented and
// raise an error.
//
// The server can be initialized with a specific function for each grpc method
// that will be called.
type MockOrchestrationManagementServer struct {
	pb.UnimplementedOrchestrationManagementServer
	CreateOrchestrationFn func(
		context.Context,
		*pb.CreateOrchestrationRequest,
	) (*emptypb.Empty, error)
	GetOrchestrationFn func(
		*pb.GetOrchestrationRequest,
		pb.OrchestrationManagement_GetServer,
	) error
}

func (s *MockOrchestrationManagementServer) Create(
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

func (s *MockOrchestrationManagementServer) Get(
	req *pb.GetOrchestrationRequest,
	stream pb.OrchestrationManagement_GetServer,
) error {
	if s.GetOrchestrationFn != nil {
		return s.GetOrchestrationFn(req, stream)
	}
	return fmt.Errorf(
		"method GetOrchestration not configured but called with req %v",
		req,
	)
}

// MockStageManagementServer is a mocking of a stage management server to be used
// during tests. By default, methods are not implemented and raise an error.
//
// The server can be initialized with a specific function for each grpc method
// that will be called.
type MockStageManagementServer struct {
	pb.UnimplementedStageManagementServer
	CreateStageFn func(context.Context, *pb.CreateStageRequest) (
		*emptypb.Empty,
		error,
	)
	GetStageFn func(*pb.GetStageRequest, pb.StageManagement_GetServer) error
}

func (s *MockStageManagementServer) Create(
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

func (s *MockStageManagementServer) Get(
	req *pb.GetStageRequest,
	stream pb.StageManagement_GetServer,
) error {
	if s.GetStageFn != nil {
		return s.GetStageFn(req, stream)
	}
	return fmt.Errorf(
		"method GetStage not configured but called with req %v",
		req,
	)
}

// MockLinkManagementServer is a mocking of a link management server to be used
// during tests. By default, methods are not implemented and raise an error.
//
// The server can be initialized with a specific function for each grpc method
// that will be called.
type MockLinkManagementServer struct {
	pb.UnimplementedLinkManagementServer
	CreateLinkFn func(context.Context, *pb.CreateLinkRequest) (
		*emptypb.Empty,
		error,
	)
	GetLinkFn func(*pb.GetLinkRequest, pb.LinkManagement_GetServer) error
}

func (s *MockLinkManagementServer) Create(
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

func (s *MockLinkManagementServer) Get(
	req *pb.GetLinkRequest,
	stream pb.LinkManagement_GetServer,
) error {
	if s.GetLinkFn != nil {
		return s.GetLinkFn(req, stream)
	}
	return fmt.Errorf(
		"method GetLink not configured but called with req %v",
		req,
	)
}

// MockAssetManagementServer is a mocking of an asset management server to be used
// during tests. By default, methods are not implemented and raise an error.
//
// The server can be initialized with a specific function for each grpc method
// that will be called.
type MockAssetManagementServer struct {
	pb.UnimplementedAssetManagementServer
	CreateAssetFn func(
		context.Context,
		*pb.CreateAssetRequest,
	) (*emptypb.Empty, error)
	GetAssetFn func(*pb.GetAssetRequest, pb.AssetManagement_GetServer) error
}

func (s *MockAssetManagementServer) Create(
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

func (s *MockAssetManagementServer) Get(
	req *pb.GetAssetRequest,
	stream pb.AssetManagement_GetServer,
) error {
	if s.GetAssetFn != nil {
		return s.GetAssetFn(req, stream)
	}
	return fmt.Errorf(
		"method GetAsset not configured but called with request %v",
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
