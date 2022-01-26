package server

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/logs"
	"github.com/DuarteMRAlves/maestro/internal/testutil"
	"github.com/dgraph-io/badger/v3"
	"gotest.tools/v3/assert"
	"testing"
)

func TestServer_CreateOrchestration(t *testing.T) {
	const name = "orchestration-name"
	tests := []struct {
		name string
		req  *api.CreateOrchestrationRequest
	}{
		{
			name: "correct with name",
			req: &api.CreateOrchestrationRequest{
				Name: name,
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

				err = s.CreateOrchestration(test.req)
				assert.NilError(t, err, "create orchestration error")
			},
		)
	}
}

func TestServer_CreateOrchestration_NilConfig(t *testing.T) {
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

	err = s.CreateOrchestration(nil)
	assert.Assert(
		t,
		errdefs.IsInvalidArgument(err),
		"error is not InvalidArgument",
	)
	expectedMsg := "'req' is nil"
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateOrchestration_InvalidName(t *testing.T) {
	tests := []struct {
		name string
		req  *api.CreateOrchestrationRequest
	}{
		{
			name: "empty name",
			req: &api.CreateOrchestrationRequest{
				Name: "",
				Links: []api.LinkName{
					testutil.LinkNameForNum(0),
					testutil.LinkNameForNum(1),
					testutil.LinkNameForNum(2),
				},
			},
		},
		{
			name: "invalid characters in name",
			req: &api.CreateOrchestrationRequest{
				Name: "?orchestration-name",
				Links: []api.LinkName{
					testutil.LinkNameForNum(0),
					testutil.LinkNameForNum(1),
					testutil.LinkNameForNum(2),
				},
			},
		},
		{
			name: "invalid character sequence",
			req: &api.CreateOrchestrationRequest{
				Name: "invalid//name",
				Links: []api.LinkName{
					testutil.LinkNameForNum(0),
					testutil.LinkNameForNum(1),
					testutil.LinkNameForNum(2),
				},
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

				err = s.CreateOrchestration(test.req)
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

func TestServer_CreateOrchestration_AlreadyExists(t *testing.T) {
	var err error
	const name = "orchestration-name"
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

	req := &api.CreateOrchestrationRequest{
		Name: name,
		Links: []api.LinkName{
			testutil.LinkNameForNum(0),
			testutil.LinkNameForNum(1),
		},
	}

	err = s.CreateOrchestration(req)
	assert.NilError(t, err, "first creation has an error")
	err = s.CreateOrchestration(req)
	assert.Assert(t, errdefs.IsAlreadyExists(err), "error is not AlreadyExists")
	expectedMsg := fmt.Sprintf("orchestration '%v' already exists", name)
	assert.Error(t, err, expectedMsg)
}
