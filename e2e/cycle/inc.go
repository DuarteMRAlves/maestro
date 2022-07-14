package cycle

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

type incServer struct {
	UnimplementedIncServer
}

func (s *incServer) Inc(
	_ context.Context, m *ValMessage,
) (*ValMessage, error) {
	delay := time.Duration(rand.Int63n(5))
	time.Sleep(delay * time.Millisecond)
	return &ValMessage{Val: m.Val + 1}, nil
}

func ServeInc() (net.Addr, func()) {
	s := grpc.NewServer()
	RegisterIncServer(s, &incServer{})
	reflection.Register(s)
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(fmt.Sprintf("inc listen: %s", err))
	}
	log.Printf("inc server listening at %v", lis.Addr())
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(fmt.Sprintf("inc serve: %v", err))
		}
	}()
	return lis.Addr(), s.Stop
}
