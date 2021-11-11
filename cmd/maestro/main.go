package main

import (
	"github.com/DuarteMRAlves/maestro/internal/server"
	"log"
	"net"
)

func main() {
	address := "localhost:50051"
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("Server listening at %v", lis.Addr())

	s := server.NewBuilder().WithGrpc().Build()

	if err := s.ServeGrpc(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
