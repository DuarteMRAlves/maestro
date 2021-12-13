package main

import (
	"github.com/DuarteMRAlves/maestro/internal/server"
	"go.uber.org/zap"
	"net"
)

func main() {
	logger, err := zap.NewProduction()
	// Should never happen
	if err != nil {
		panic(err)
	}
	sugar := logger.Sugar()
	address := "localhost:50051"
	lis, err := net.Listen("tcp", address)
	if err != nil {
		sugar.Fatal("Failed to listen.", "err", err)
	}
	sugar.Infof("Server listening at: %v", lis.Addr())

	s, err := server.NewBuilder().WithGrpc().Build()
	if err != nil {
		sugar.Fatalf("build server: %v", err)
	}

	if err := s.ServeGrpc(lis); err != nil {
		sugar.Fatal("Failed to serve.", "err", err)
	}
}
