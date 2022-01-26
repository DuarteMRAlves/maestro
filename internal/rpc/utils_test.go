package rpc

import (
	"context"
	"fmt"
	"github.com/DuarteMRAlves/maestro/tests/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gotest.tools/v3/assert"
	"net"
	"testing"
)

var dummyErr = fmt.Errorf("dummy error")

type testService struct {
	pb.UnimplementedTestServiceServer
	pb.UnimplementedExtraServiceServer
}

func (s *testService) Unary(
	ctx context.Context,
	request *pb.Request,
) (*pb.Reply, error) {

	if request.StringField == "error" {
		return nil, dummyErr
	} else {
		return replyFromRequest(request), nil
	}
}

func replyFromRequest(request *pb.Request) *pb.Reply {
	doubleField := float64(len(request.StringField))
	for _, val := range request.RepeatedField {
		doubleField += float64(val)
	}

	innerMsg := &pb.InnerMessage{RepeatedString: []string{}}
	for _, inner := range request.RepeatedInnerMsg {
		repeatedString := ""
		for _, str := range inner.RepeatedString {
			repeatedString += str
		}
		innerMsg.RepeatedString = append(
			innerMsg.RepeatedString,
			repeatedString,
		)
	}
	return &pb.Reply{
		DoubleField: doubleField,
		InnerMsg:    innerMsg,
	}
}

func startServer(
	t *testing.T,
	lis net.Listener,
	reflectionFlag bool,
) *grpc.Server {
	testServer := grpc.NewServer()
	pb.RegisterTestServiceServer(testServer, &testService{})
	pb.RegisterExtraServiceServer(testServer, &testService{})

	if reflectionFlag {
		reflection.Register(testServer)
	}

	go func() {
		err := testServer.Serve(lis)
		assert.NilError(t, err, "test server error")
	}()
	return testServer
}
