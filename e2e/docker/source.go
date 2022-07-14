package docker

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

type sourceServer struct {
	UnimplementedSourceServer
	counter int64
}

func (s *sourceServer) Generate(
	_ context.Context, _ *emptypb.Empty,
) (*Message, error) {
	delay := time.Duration(rand.Int63n(5))
	time.Sleep(delay * time.Millisecond)
	val := atomic.AddInt64(&s.counter, 1)
	return &Message{Val: val}, nil
}

func ServeSource() (net.Addr, func()) {
	s := grpc.NewServer()
	RegisterSourceServer(s, &sourceServer{counter: 0})
	reflection.Register(s)
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(fmt.Sprintf("source listen: %s", err))
	}
	log.Printf("source server listening at %v", lis.Addr())
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(fmt.Sprintf("source serve: %v", err))
		}
	}()
	return lis.Addr(), s.Stop
}
