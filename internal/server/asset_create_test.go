package server

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/logs"
	"github.com/dgraph-io/badger/v3"
	"gotest.tools/v3/assert"
	"testing"
)

func TestServer_CreateAsset(t *testing.T) {
	const name = "asset-name"

	tests := []struct {
		name string
		req  *api.CreateAssetRequest
	}{
		{
			name: "correct nil image",
			req: &api.CreateAssetRequest{
				Name: name,
			},
		},
		{
			name: "correct with empty image",
			req: &api.CreateAssetRequest{
				Name:  name,
				Image: "",
			},
		},
		{
			name: "correct with image",
			req: &api.CreateAssetRequest{
				Name:  name,
				Image: "image-name",
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
				err = s.CreateAsset(test.req)
				assert.NilError(t, err, "create asset error")
			},
		)
	}
}

func TestServer_CreateAsset_NilConfig(t *testing.T) {
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

	err = s.CreateAsset(nil)
	assert.Assert(
		t,
		errdefs.IsInvalidArgument(err),
		"error is not InvalidArgument",
	)
	expectedMsg := "'req' is nil"
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateAsset_InvalidName(t *testing.T) {
	tests := []struct {
		name string
		req  *api.CreateAssetRequest
	}{
		{
			name: "empty name",
			req: &api.CreateAssetRequest{
				Name: "",
			},
		},
		{
			name: "invalid characters in name",
			req: &api.CreateAssetRequest{
				Name: "%some-name%",
			},
		},
		{
			name: "invalid character sequence",
			req: &api.CreateAssetRequest{
				Name: "invalid..name",
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

				err = s.CreateAsset(test.req)
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

func TestServer_CreateAsset_AlreadyExists(t *testing.T) {
	var err error
	const name = "asset-name"

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

	req := &api.CreateAssetRequest{
		Name: name,
	}

	err = s.CreateAsset(req)
	assert.NilError(t, err, "first creation has an error")
	err = s.CreateAsset(req)
	assert.Assert(t, errdefs.IsAlreadyExists(err), "error is not NotFound")
	expectedMsg := fmt.Sprintf("asset '%v' already exists", name)
	assert.Error(t, err, expectedMsg)
}
