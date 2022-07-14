package splitmerge

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

type sinkServer struct {
	UnimplementedSinkServer
	fn func(msg *Compose)
}

func (s *sinkServer) Collect(
	_ context.Context, m *Compose,
) (*emptypb.Empty, error) {
	s.fn(m)
	delay := time.Duration(rand.Int63n(5))
	time.Sleep(delay * time.Millisecond)
	return &emptypb.Empty{}, nil
}

func ServeSink(fn func(msg *Compose)) (net.Addr, func()) {
	s := grpc.NewServer()
	RegisterSinkServer(s, &sinkServer{fn: fn})
	reflection.Register(s)
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(fmt.Sprintf("sink listen: %s", err))
	}
	log.Printf("sink server listening at %v", lis.Addr())
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(fmt.Sprintf("sink serve: %v", err))
		}
	}()
	return lis.Addr(), s.Stop
}
