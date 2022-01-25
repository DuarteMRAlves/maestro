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
	CreateOrchestrationFn func(ctx context.Context, config *pb.Orchestration) (
		*emptypb.Empty,
		error,
	)
	GetOrchestrationFn func(
		query *pb.Orchestration,
		stream pb.OrchestrationManagement_GetServer,
	) error
}

func (s *MockOrchestrationManagementServer) Create(
	ctx context.Context,
	config *pb.Orchestration,
) (*emptypb.Empty, error) {
	if s.CreateOrchestrationFn != nil {
		return s.CreateOrchestrationFn(ctx, config)
	}
	return &emptypb.Empty{}, fmt.Errorf(
		"method CreateOrchestration not configured but called with config %v",
		config,
	)
}

func (s *MockOrchestrationManagementServer) Get(
	query *pb.Orchestration,
	stream pb.OrchestrationManagement_GetServer,
) error {
	if s.GetOrchestrationFn != nil {
		return s.GetOrchestrationFn(query, stream)
	}
	return fmt.Errorf(
		"method GetOrchestration not configured but called with query %v",
		query,
	)
}

// MockStageManagementServer is a mocking of a stage management server to be used
// during tests. By default, methods are not implemented and raise an error.
//
// The server can be initialized with a specific function for each grpc method
// that will be called.
type MockStageManagementServer struct {
	pb.UnimplementedStageManagementServer
	CreateStageFn func(ctx context.Context, config *pb.Stage) (
		*emptypb.Empty,
		error,
	)
	GetStageFn func(
		query *pb.Stage,
		stream pb.StageManagement_GetServer,
	) error
}

func (s *MockStageManagementServer) Create(
	ctx context.Context,
	config *pb.Stage,
) (*emptypb.Empty, error) {
	if s.CreateStageFn != nil {
		return s.CreateStageFn(ctx, config)
	}
	return &emptypb.Empty{}, fmt.Errorf(
		"method CreateStage not configured but called with config %v",
		config,
	)
}

func (s *MockStageManagementServer) Get(
	query *pb.Stage,
	stream pb.StageManagement_GetServer,
) error {
	if s.GetStageFn != nil {
		return s.GetStageFn(query, stream)
	}
	return fmt.Errorf(
		"method GetStage not configured but called with query %v",
		query,
	)
}

// MockLinkManagementServer is a mocking of a link management server to be used
// during tests. By default, methods are not implemented and raise an error.
//
// The server can be initialized with a specific function for each grpc method
// that will be called.
type MockLinkManagementServer struct {
	pb.UnimplementedLinkManagementServer
	CreateLinkFn func(ctx context.Context, config *pb.Link) (
		*emptypb.Empty,
		error,
	)
	GetLinkFn func(
		query *pb.Link,
		stream pb.LinkManagement_GetServer,
	) error
}

func (s *MockLinkManagementServer) Create(
	ctx context.Context,
	config *pb.Link,
) (*emptypb.Empty, error) {
	if s.CreateLinkFn != nil {
		return s.CreateLinkFn(ctx, config)
	}
	return &emptypb.Empty{}, fmt.Errorf(
		"method CreateLink not configured but called with config %v",
		config,
	)
}

func (s *MockLinkManagementServer) Get(
	query *pb.Link,
	stream pb.LinkManagement_GetServer,
) error {
	if s.GetLinkFn != nil {
		return s.GetLinkFn(query, stream)
	}
	return fmt.Errorf(
		"method GetLink not configured but called with query %v",
		query,
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
