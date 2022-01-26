package server

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/logs"
	"github.com/DuarteMRAlves/maestro/internal/storage"
	"github.com/DuarteMRAlves/maestro/internal/util"
	"github.com/dgraph-io/badger/v3"
	"gotest.tools/v3/assert"
	"testing"
)

func TestServer_CreateStage(t *testing.T) {
	const name = "stage-name"

	tests := []struct {
		name string
		req  *api.CreateStageRequest
	}{
		{
			name: "nil asset, service and rpc",
			req: &api.CreateStageRequest{
				Name:    name,
				Address: "some-address",
			},
		},
		{
			name: "from host and port",
			req: &api.CreateStageRequest{
				Name:    name,
				Asset:   util.AssetNameForNum(0),
				Service: "Service",
				Rpc:     "Method",
				Host:    "host",
				Port:    int32(12345),
			},
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				db, err := badger.Open(
					badger.DefaultOptions("").WithInMemory(true),
				)
				assert.NilError(t, err, "db creation")
				defer db.Close()
				s, err := NewBuilder().
					WithGrpc().
					WithDb(db).
					WithLogger(logs.NewTestLogger(t)).
					Build()
				assert.NilError(t, err, "build server")
				err = db.Update(
					func(txn *badger.Txn) error {
						populateForStages(t, txn)
						return nil
					},
				)
				assert.NilError(t, err, "Populate error")
				err = s.CreateStage(test.req)
				assert.NilError(t, err, "create stage error")
			},
		)
	}
}

func TestServer_CreateStage_NilConfig(t *testing.T) {
	db, err := badger.Open(
		badger.DefaultOptions("").WithInMemory(true),
	)
	assert.NilError(t, err, "db creation")
	defer db.Close()
	s, err := NewBuilder().
		WithGrpc().
		WithDb(db).
		WithLogger(logs.NewTestLogger(t)).
		Build()
	assert.NilError(t, err, "build server")
	err = db.Update(
		func(txn *badger.Txn) error {
			populateForStages(t, txn)
			return nil
		},
	)
	assert.NilError(t, err, "Populate error")

	err = s.CreateStage(nil)
	assert.Assert(
		t,
		errdefs.IsInvalidArgument(err),
		"error is not InvalidArgument",
	)
	expectedMsg := "'req' is nil"
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateStage_InvalidName(t *testing.T) {
	tests := []struct {
		name string
		req  *api.CreateStageRequest
	}{
		{
			name: "empty name",
			req: &api.CreateStageRequest{
				Name:  "",
				Asset: util.AssetNameForNum(0),
			},
		},
		{
			name: "invalid characters in name",
			req: &api.CreateStageRequest{
				Name:  "some@name",
				Asset: util.AssetNameForNum(0),
			},
		},
		{
			name: "invalid character sequence",
			req: &api.CreateStageRequest{
				Name:  "other-/name",
				Asset: util.AssetNameForNum(0),
			},
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				db, err := badger.Open(
					badger.DefaultOptions("").WithInMemory(true),
				)
				assert.NilError(t, err, "db creation")
				defer db.Close()
				s, err := NewBuilder().
					WithGrpc().
					WithDb(db).
					WithLogger(logs.NewTestLogger(t)).
					Build()
				assert.NilError(t, err, "build server")
				err = db.Update(
					func(txn *badger.Txn) error {
						populateForStages(t, txn)
						return nil
					},
				)
				assert.NilError(t, err, "Populate error")

				err = s.CreateStage(test.req)
				assert.Assert(
					t,
					errdefs.IsInvalidArgument(err),
					"error is not InvalidArgument",
				)
				expectedMsg := fmt.Sprintf(
					"invalid name '%v'",
					test.req.Name,
				)
				assert.Error(t, err, expectedMsg)
			},
		)
	}
}

func TestServer_CreateStage_AssetNotFound(t *testing.T) {
	const name = "stage-name"
	db, err := badger.Open(
		badger.DefaultOptions("").WithInMemory(true),
	)
	assert.NilError(t, err, "db creation")
	defer db.Close()
	s, err := NewBuilder().
		WithGrpc().
		WithDb(db).
		WithLogger(logs.NewTestLogger(t)).
		Build()
	assert.NilError(t, err, "build server")
	err = db.Update(
		func(txn *badger.Txn) error {
			populateForStages(t, txn)
			return nil
		},
	)
	assert.NilError(t, err, "Populate error")

	config := &api.CreateStageRequest{
		Name:  name,
		Asset: util.AssetNameForNum(1),
	}

	err = s.CreateStage(config)
	assert.Assert(t, errdefs.IsNotFound(err), "error is not NotFound")
	expectedMsg := fmt.Sprintf(
		"asset '%v' not found",
		util.AssetNameForNum(1),
	)
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateStage_AlreadyExists(t *testing.T) {
	var err error
	const name = "stage-name"

	db, err := badger.Open(
		badger.DefaultOptions("").WithInMemory(true),
	)
	assert.NilError(t, err, "db creation")
	defer db.Close()
	s, err := NewBuilder().
		WithGrpc().
		WithDb(db).
		WithLogger(logs.NewTestLogger(t)).
		Build()
	assert.NilError(t, err, "build server")
	err = db.Update(
		func(txn *badger.Txn) error {
			populateForStages(t, txn)
			return nil
		},
	)
	assert.NilError(t, err, "Populate error")

	config := &api.CreateStageRequest{
		Name:    name,
		Asset:   util.AssetNameForNum(0),
		Service: "Service",
		Rpc:     "Method",
		Address: "address",
	}

	err = s.CreateStage(config)
	assert.NilError(t, err, "first creation has an error")
	err = s.CreateStage(config)
	assert.Assert(t, errdefs.IsAlreadyExists(err), "error is not AlreadyExists")
	expectedMsg := fmt.Sprintf("stage '%v' already exists", name)
	assert.Error(t, err, expectedMsg)
}

func populateForStages(t *testing.T, txn *badger.Txn) {
	assets := []*api.Asset{
		assetForNum(0),
	}
	helper := storage.NewTxnHelper(txn)
	for _, a := range assets {
		assert.NilError(t, helper.SaveAsset(a))
	}
}
