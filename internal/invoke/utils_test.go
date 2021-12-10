package invoke

import (
	"context"
	"fmt"
	"github.com/DuarteMRAlves/maestro/tests/pb"
	"google.golang.org/grpc"
	"gotest.tools/v3/assert"
	"net"
	"testing"
)

var dummyErr = fmt.Errorf("dummy error")

type service struct {
	pb.UnimplementedTestServiceServer
	pb.UnimplementedExtraServiceServer
}

func (s *service) Unary(
	ctx context.Context,
	request *pb.Request,
) (*pb.Reply, error) {

	if request.StringField == "error" {
		return nil, dummyErr
	} else {
		return replyFromRequest(request), nil
	}
}

func startServer(t *testing.T, addr string) *grpc.Server {
	lis, err := net.Listen("tcp", addr)
	assert.NilError(t, err, "listen error")

	testServer := grpc.NewServer()
	pb.RegisterTestServiceServer(testServer, &service{})
	pb.RegisterExtraServiceServer(testServer, &service{})

	go func() {
		err = testServer.Serve(lis)
		assert.NilError(t, err, "test server error")
	}()
	return testServer
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
			repeatedString)
	}
	return &pb.Reply{
		DoubleField: doubleField,
		InnerMsg:    innerMsg,
	}
}
