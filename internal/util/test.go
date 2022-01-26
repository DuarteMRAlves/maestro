package util

import (
	"context"
	"fmt"
	"github.com/DuarteMRAlves/maestro/tests/pb"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gotest.tools/v3/assert"
	"net"
	"testing"
)

// NewTestListener returns a new Listener on an empty port where the server
// can run. The client should connect to the address using the Addr method.
func NewTestListener(t *testing.T) net.Listener {
	lis, err := net.Listen("tcp", "localhost:0")
	assert.NilError(t, err, "failed to listen")
	return lis
}

type testService struct {
	pb.UnimplementedTestServiceServer
	pb.UnimplementedExtraServiceServer
}

func (s *testService) Unary(
	ctx context.Context,
	request *pb.Request,
) (*pb.Reply, error) {

	if request.StringField == "error" {
		return nil, fmt.Errorf("dummy error")
	}
	return replyFromRequest(request), nil
}

func StartTestServer(
	t *testing.T,
	lis net.Listener,
	registerTest bool,
	registerExtra bool,
) *grpc.Server {
	testServer := grpc.NewServer()
	if registerTest {
		pb.RegisterTestServiceServer(testServer, &testService{})
	}
	if registerExtra {
		pb.RegisterExtraServiceServer(testServer, &testService{})
	}

	reflection.Register(testServer)

	go func() {
		err := testServer.Serve(lis)
		assert.NilError(t, err, "test server error")
	}()
	return testServer
}

func AssertUnaryRequest(t *testing.T, req *pb.Request, rep *pb.Reply) {
	expected := replyFromRequest(req)
	opts := cmpopts.IgnoreUnexported(pb.Reply{}, pb.InnerMessage{})
	assert.DeepEqual(t, expected, rep, opts)
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
