package docker

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type transformServer struct {
	UnimplementedTransformServer
}

func (s *transformServer) Process(
	_ context.Context, m *Message,
) (*Message, error) {
	delay := time.Duration(rand.Int63n(5))
	time.Sleep(delay * time.Millisecond)
	return &Message{Val: m.Val * 2}, nil
}

func ServeTransform() (net.Addr, func()) {
	s := grpc.NewServer()
	RegisterTransformServer(s, &transformServer{})
	reflection.Register(s)
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(fmt.Sprintf("transform listen: %s", err))
	}
	log.Printf("transform server listening at %v", lis.Addr())
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(fmt.Sprintf("transform serve: %v", err))
		}
	}()
	return lis.Addr(), s.Stop
}
