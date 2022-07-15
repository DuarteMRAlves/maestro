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

type sumServer struct {
	UnimplementedSumServer
	fn func(msg *SumMessage)
}

func (s *sumServer) Sum(
	_ context.Context, m *SumMessage,
) (*ValMessage, error) {
	s.fn(m)
	delay := time.Duration(rand.Int63n(5))
	time.Sleep(delay * time.Millisecond)
	return &ValMessage{Val: m.Inc.Val + m.Counter.Val}, nil
}

func ServeSum(fn func(msg *SumMessage)) (net.Addr, func()) {
	s := grpc.NewServer()
	RegisterSumServer(s, &sumServer{fn: fn})
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
