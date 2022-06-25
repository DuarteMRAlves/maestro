package grpcw

import (
	"context"
	"fmt"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

type blockingReflectionStream struct {
	stream grpc_reflection_v1alpha.ServerReflection_ServerReflectionInfoClient
	mu     sync.Mutex
}

func newBlockingReflectionStream(
	ctx context.Context, conn grpc.ClientConnInterface,
) (*blockingReflectionStream, error) {
	var s blockingReflectionStream
	stub := grpc_reflection_v1alpha.NewServerReflectionClient(conn)
	stream, err := stub.ServerReflectionInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("create grpc reflection stream: %w", err)
	}
	s.stream = stream
	s.mu = sync.Mutex{}
	return &s, nil
}

func (s *blockingReflectionStream) listServiceNames() ([]string, error) {
	if s == nil || s.stream == nil {
		return nil, fmt.Errorf("not connected")
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	req := &grpc_reflection_v1alpha.ServerReflectionRequest{
		MessageRequest: &grpc_reflection_v1alpha.ServerReflectionRequest_ListServices{
			// this field value is ignored by grpc
			ListServices: "*",
		},
	}
	if err := s.stream.Send(req); err != nil {
		return nil, fmt.Errorf("send list services request: %w", err)
	}
	rep, err := s.stream.Recv()
	if err != nil {
		return nil, fmt.Errorf("recv list services response: %w", err)
	}

	switch r := rep.MessageResponse.(type) {
	case *grpc_reflection_v1alpha.ServerReflectionResponse_ListServicesResponse:
		names := make([]string, 0, len(r.ListServicesResponse.Service))
		for _, s := range r.ListServicesResponse.Service {
			names = append(names, s.Name)
		}
		return names, nil
	default:
		return nil, fmt.Errorf("invalid list services response type: %q", r)
	}
}

func (s *blockingReflectionStream) filesForSymbol(symb string) ([][]byte, error) {
	if s == nil || s.stream == nil {
		return nil, fmt.Errorf("not connected")
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	req := &grpc_reflection_v1alpha.ServerReflectionRequest{
		MessageRequest: &grpc_reflection_v1alpha.ServerReflectionRequest_FileContainingSymbol{
			FileContainingSymbol: symb,
		},
	}

	if err := s.stream.Send(req); err != nil {
		return nil, fmt.Errorf("send find file containing symbol request: %w", err)
	}
	rep, err := s.stream.Recv()
	if err != nil {
		return nil, fmt.Errorf("recv find file containing symbol response: %w", err)
	}
	switch r := rep.MessageResponse.(type) {
	case *grpc_reflection_v1alpha.ServerReflectionResponse_FileDescriptorResponse:
		fds := make([][]byte, 0, len(r.FileDescriptorResponse.FileDescriptorProto))
		fds = append(fds, r.FileDescriptorResponse.FileDescriptorProto...)
		return fds, nil
	default:
		return nil, fmt.Errorf("invalid find file containing symbol response type: %q", r)
	}
}
