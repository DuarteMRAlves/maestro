package execution

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/rpc"
	"github.com/DuarteMRAlves/maestro/internal/storage"
	"github.com/DuarteMRAlves/maestro/internal/util"
	"github.com/DuarteMRAlves/maestro/tests/pb"
	"github.com/dgraph-io/badger/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gotest.tools/v3/assert"
	"net"
	"testing"
)

func TestBuilder_build(t *testing.T) {
	var err error

	lis1 := util.NewTestListener(t)
	addr1 := lis1.Addr().String()
	server1 := startLinear1(t, lis1)
	defer server1.Stop()

	lis2 := util.NewTestListener(t)
	addr2 := lis2.Addr().String()
	server2 := startLinear2(t, lis2)
	defer server2.Stop()

	rpcManager := rpc.NewManager()

	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	assert.NilError(t, err, "create db")
	defer db.Close()

	err = initDb(db, addr1, addr2)
	assert.NilError(t, err, "init db")

	err = db.View(
		func(txn *badger.Txn) error {
			builder := newBuilder(txn, rpcManager)
			_, err = builder.withOrchestration("default").build()
			return err
		},
	)
	assert.NilError(t, err, "build error")
}

type LinearService1 struct {
	pb.UnimplementedLinear1Server
}

type LinearService2 struct {
	pb.UnimplementedLinear2Server
}

func startLinear1(t *testing.T, lis net.Listener) *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterLinear1Server(s, &LinearService1{})
	reflection.Register(s)

	go func() {
		err := s.Serve(lis)
		assert.NilError(t, err, "linear 1 server error")
	}()
	return s
}

func startLinear2(t *testing.T, lis net.Listener) *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterLinear2Server(s, &LinearService2{})
	reflection.Register(s)

	go func() {
		err := s.Serve(lis)
		assert.NilError(t, err, "linear 2 server error")
	}()
	return s
}

func initDb(db *badger.DB, addr1 string, addr2 string) error {
	var err error

	o := &api.Orchestration{
		Name:   "default",
		Phase:  api.OrchestrationPending,
		Stages: []api.StageName{"linear1", "linear2"},
		Links:  []api.LinkName{"link-1-2"},
	}

	linear1 := &api.Stage{
		Name:          "linear1",
		Phase:         api.StagePending,
		Service:       "pb.Linear1",
		Rpc:           "Process",
		Address:       addr1,
		Orchestration: "default",
		Asset:         "",
	}

	linear2 := &api.Stage{
		Name:          "linear2",
		Phase:         api.StagePending,
		Service:       "pb.Linear2",
		Rpc:           "Process",
		Address:       addr2,
		Orchestration: "default",
		Asset:         "",
	}

	l := &api.Link{
		Name:          "link-1-2",
		SourceStage:   "linear1",
		SourceField:   "",
		TargetStage:   "linear2",
		TargetField:   "",
		Orchestration: "default",
	}

	return db.Update(
		func(txn *badger.Txn) error {
			helper := storage.NewTxnHelper(txn)
			err = helper.SaveOrchestration(o)
			if err != nil {
				return err
			}
			err = helper.SaveStage(linear1)
			if err != nil {
				return err
			}
			err = helper.SaveStage(linear2)
			if err != nil {
				return err
			}
			err = helper.SaveLink(l)
			if err != nil {
				return err
			}
			return nil
		},
	)
}
