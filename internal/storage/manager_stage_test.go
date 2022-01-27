package storage

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/rpc"
	"github.com/DuarteMRAlves/maestro/internal/util"
	"github.com/dgraph-io/badger/v3"
	"gotest.tools/v3/assert"
	"testing"
)

func TestManager_CreateStage(t *testing.T) {
	const name api.StageName = "Stage-Name"
	tests := []struct {
		name     string
		req      *api.CreateStageRequest
		expected *api.Stage
	}{
		{
			name: "required parameters",
			req: &api.CreateStageRequest{
				Name: name,
			},
			expected: &api.Stage{
				Name:    name,
				Phase:   api.StagePending,
				Service: "",
				Rpc:     "",
				Address: fmt.Sprintf(
					"%s:%d",
					defaultStageHost,
					defaultStagePort,
				),
				Orchestration: api.OrchestrationName("default"),
				Asset:         api.AssetName(""),
			},
		},
		{
			name: "custom address",
			req: &api.CreateStageRequest{
				Name:    name,
				Address: "some-address",
			},
			expected: &api.Stage{
				Name:          name,
				Phase:         api.StagePending,
				Service:       "",
				Rpc:           "",
				Address:       "some-address",
				Orchestration: api.OrchestrationName("default"),
				Asset:         api.AssetName(""),
			},
		},
		{
			name: "custom host and port",
			req: &api.CreateStageRequest{
				Name: name,
				Host: "some-host",
				Port: 12345,
			},
			expected: &api.Stage{
				Name:          name,
				Phase:         api.StagePending,
				Service:       "",
				Rpc:           "",
				Address:       "some-host:12345",
				Orchestration: api.OrchestrationName("default"),
				Asset:         api.AssetName(""),
			},
		},
		{
			name: "all parameters",
			req: &api.CreateStageRequest{
				Name:          name,
				Service:       "Service",
				Rpc:           "Rpc",
				Address:       "Address",
				Host:          "",
				Port:          0,
				Orchestration: util.OrchestrationNameForNum(0),
				Asset:         util.AssetNameForNum(0),
			},
			expected: &api.Stage{
				Name:          name,
				Phase:         api.StagePending,
				Service:       "Service",
				Rpc:           "Rpc",
				Address:       "Address",
				Orchestration: util.OrchestrationNameForNum(0),
				Asset:         util.AssetNameForNum(0),
			},
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				testCreateStage(t, test.req, test.expected)
			},
		)
	}
}

func testCreateStage(
	t *testing.T,
	req *api.CreateStageRequest,
	expected *api.Stage,
) {
	var stored api.Stage
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	assert.NilError(t, err, "db creation")
	defer db.Close()

	m, err := NewManager(NewDefaultContext(db, rpc.NewManager()))
	assert.NilError(t, err, "manager creation")

	err = db.Update(
		func(txn *badger.Txn) error {
			helper := NewTxnHelper(txn)
			a := &api.Asset{Name: util.AssetNameForNum(0)}
			if err := helper.SaveAsset(a); err != nil {
				return err
			}
			o := &api.Orchestration{Name: util.OrchestrationNameForNum(0)}
			if err := helper.SaveOrchestration(o); err != nil {
				return err
			}
			return nil
		},
	)
	assert.NilError(t, err, "setup db error")

	err = db.Update(
		func(txn *badger.Txn) error {
			return m.CreateStage(txn, req)
		},
	)
	assert.NilError(t, err, "create error not nil")
	err = db.View(
		func(txn *badger.Txn) error {
			helper := TxnHelper{txn: txn}
			return helper.LoadStage(&stored, req.Name)
		},
	)
	assert.NilError(t, err, "load error")
	assert.Equal(t, expected.Name, stored.Name, "name not equal")
	assert.Equal(t, expected.Phase, stored.Phase, "phase not equal")
	assert.Equal(t, expected.Service, stored.Service, "service not equal")
	assert.Equal(t, expected.Rpc, stored.Rpc, "rpc not equal")
	assert.Equal(t, expected.Address, stored.Address, "address not equal")
	assert.Equal(
		t,
		expected.Orchestration,
		stored.Orchestration,
		"orchestration not equal",
	)
	assert.Equal(t, expected.Asset, stored.Asset, "asset not equal")
}

func TestManager_CreateStage_Error(t *testing.T) {
	tests := []struct {
		name            string
		req             *api.CreateStageRequest
		assertErrTypeFn func(error) bool
		expectedErrMsg  string
	}{
		{
			name:            "nil config",
			req:             nil,
			assertErrTypeFn: errdefs.IsInvalidArgument,
			expectedErrMsg:  "'req' is nil",
		},
		{
			name:            "empty name",
			req:             &api.CreateStageRequest{Name: ""},
			assertErrTypeFn: errdefs.IsInvalidArgument,
			expectedErrMsg:  "invalid name ''",
		},
		{
			name:            "invalid name",
			req:             &api.CreateStageRequest{Name: "some#name"},
			assertErrTypeFn: errdefs.IsInvalidArgument,
			expectedErrMsg:  "invalid name 'some#name'",
		},
		{
			name:            "stage already exists",
			req:             &api.CreateStageRequest{Name: "duplicate"},
			assertErrTypeFn: errdefs.IsAlreadyExists,
			expectedErrMsg:  "stage 'duplicate' already exists",
		},
		{
			name: "orchestration not found",
			req: &api.CreateStageRequest{
				Name:          "some-stage",
				Orchestration: "unknown",
			},
			assertErrTypeFn: errdefs.IsNotFound,
			expectedErrMsg:  "orchestration 'unknown' not found",
		},
		{
			name: "asset not found",
			req: &api.CreateStageRequest{
				Name:  "some-stage",
				Asset: "unknown",
			},
			assertErrTypeFn: errdefs.IsNotFound,
			expectedErrMsg:  "asset 'unknown' not found",
		},
		{
			name: "stage and host specified",
			req: &api.CreateStageRequest{
				Name:    "some-stage",
				Address: "some-address",
				Host:    "some-host",
			},
			assertErrTypeFn: errdefs.IsInvalidArgument,
			expectedErrMsg:  "Cannot simultaneously specify address and host for stage",
		},
		{
			name: "stage and port specified",
			req: &api.CreateStageRequest{
				Name:    "some-stage",
				Address: "some-address",
				Port:    23456,
			},
			assertErrTypeFn: errdefs.IsInvalidArgument,
			expectedErrMsg:  "Cannot simultaneously specify address and port for stage",
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				testCreateStageError(
					t,
					test.req,
					test.assertErrTypeFn,
					test.expectedErrMsg,
				)
			},
		)
	}
}

func testCreateStageError(
	t *testing.T,
	req *api.CreateStageRequest,
	assertErrTypeFn func(error) bool,
	expectedErrMsg string,
) {
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	assert.NilError(t, err, "db creation")
	defer db.Close()

	m, err := NewManager(NewDefaultContext(db, rpc.NewManager()))
	assert.NilError(t, err, "manager creation")

	err = db.Update(
		func(txn *badger.Txn) error {
			helper := NewTxnHelper(txn)
			a := &api.Asset{Name: util.AssetNameForNum(0)}
			if err := helper.SaveAsset(a); err != nil {
				return err
			}
			o := &api.Orchestration{Name: util.OrchestrationNameForNum(0)}
			if err := helper.SaveOrchestration(o); err != nil {
				return err
			}
			s := &api.Stage{Name: "duplicate"}
			if err := helper.SaveStage(s); err != nil {
				return err
			}
			return nil
		},
	)
	assert.NilError(t, err, "setup db error")

	err = db.Update(
		func(txn *badger.Txn) error {
			return m.CreateStage(txn, req)
		},
	)
	assert.Assert(t, assertErrTypeFn(err), "wrong error type")
	assert.Equal(t, expectedErrMsg, err.Error(), "wrong error message")
}
