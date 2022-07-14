package cycle

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

type counterServer struct {
	UnimplementedCounterServer
	counter int64
}

func (s *counterServer) Generate(
	_ context.Context, _ *emptypb.Empty,
) (*ValMessage, error) {
	delay := time.Duration(rand.Int63n(5))
	time.Sleep(delay * time.Millisecond)
	val := atomic.AddInt64(&s.counter, 1)
	return &ValMessage{Val: val}, nil
}

func ServeCounter() (net.Addr, func()) {
	s := grpc.NewServer()
	RegisterCounterServer(s, &counterServer{counter: 0})
	reflection.Register(s)
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(fmt.Sprintf("counter listen: %s", err))
	}
	log.Printf("counter server listening at %v", lis.Addr())
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(fmt.Sprintf("counter serve: %v", err))
		}
	}()
	return lis.Addr(), s.Stop
}
