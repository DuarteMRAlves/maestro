package client

import (
	"google.golang.org/grpc"
	"log"
)

const connError = "Unable to connect to maestro server at '%v'. " +
	"Is the server running?"

func NewConnection(addr string) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf(connError, addr)
	}
	return conn
}
