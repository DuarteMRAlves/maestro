package server

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	mockreflection "github.com/DuarteMRAlves/maestro/internal/reflection"
	"github.com/DuarteMRAlves/maestro/internal/storage"
	"github.com/DuarteMRAlves/maestro/internal/testutil"
	"github.com/DuarteMRAlves/maestro/tests/pb"
	"github.com/dgraph-io/badger/v3"
	"gotest.tools/v3/assert"
	"sync"
	"testing"
)

func TestServer_CreateLink(t *testing.T) {
	const name = "link-name"
	tests := []struct {
		name string
		req  *api.CreateLinkRequest
	}{
		{
			name: "correct with nil fields",
			req: &api.CreateLinkRequest{
				Name:        name,
				SourceStage: "stage-1",
				TargetStage: "stage-2",
			},
		},
		{
			name: "correct with empty fields",
			req: &api.CreateLinkRequest{
				Name:        name,
				SourceStage: "stage-1",
				SourceField: "",
				TargetStage: "stage-2",
				TargetField: "",
			},
		},
		{
			name: "correct with fields",
			req: &api.CreateLinkRequest{
				Name:        name,
				SourceStage: "stage-1",
				SourceField: "field4",
				TargetStage: "stage-2",
				TargetField: "fieldName4",
			},
		},
		{
			name: "incompatible outer but compatible inner",
			req: &api.CreateLinkRequest{
				Name:        name,
				SourceStage: "stage-1",
				SourceField: "field4",
				TargetStage: "stage-3",
				TargetField: "field4",
			},
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				rpcManager := &mockreflection.MockManager{Rpcs: sync.Map{}}
				db, err := badger.Open(
					badger.DefaultOptions("").WithInMemory(true),
				)
				assert.NilError(t, err, "db creation")
				defer db.Close()
				s, err := NewBuilder().
					WithGrpc().
					WithDb(db).
					WithReflectionManager(rpcManager).
					WithLogger(testutil.NewLogger(t)).
					Build()
				fmt.Println("On populate links")
				assert.NilError(t, err, "build server")
				populateForLinks(t, s, rpcManager)
				fmt.Println("On create link")
				err = s.CreateLink(test.req)
				assert.NilError(t, err, "create link error")
			},
		)
	}
}

func TestServer_CreateLink_NilConfig(t *testing.T) {
	rpcManager := &mockreflection.MockManager{Rpcs: sync.Map{}}
	db, err := badger.Open(
		badger.DefaultOptions("").WithInMemory(true),
	)
	assert.NilError(t, err, "db creation")
	defer db.Close()
	s, err := NewBuilder().
		WithGrpc().
		WithDb(db).
		WithReflectionManager(rpcManager).
		WithLogger(testutil.NewLogger(t)).
		Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s, rpcManager)

	err = s.CreateLink(nil)
	assert.Assert(
		t,
		errdefs.IsInvalidArgument(err),
		"error is not InvalidArgument",
	)
	expectedMsg := "'req' is nil"
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateLink_InvalidName(t *testing.T) {
	tests := []struct {
		name string
		req  *api.CreateLinkRequest
	}{
		{
			name: "empty name",
			req: &api.CreateLinkRequest{
				Name:        "",
				SourceStage: "stage-1",
				TargetStage: "stage-2",
			},
		},
		{
			name: "invalid characters in name",
			req: &api.CreateLinkRequest{
				Name:        "some'character",
				SourceStage: "stage-1",
				TargetStage: "stage-2",
			},
		},
		{
			name: "invalid character sequence",
			req: &api.CreateLinkRequest{
				Name:        "//invalid-name",
				SourceStage: "stage-1",
				TargetStage: "stage-2",
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
					WithLogger(testutil.NewLogger(t)).
					Build()
				assert.NilError(t, err, "build server")

				err = s.CreateLink(test.req)
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

func TestServer_CreateLink_SourceEmpty(t *testing.T) {
	const name = "link-name"
	rpcManager := &mockreflection.MockManager{Rpcs: sync.Map{}}
	db, err := badger.Open(
		badger.DefaultOptions("").WithInMemory(true),
	)
	assert.NilError(t, err, "db creation")
	defer db.Close()
	s, err := NewBuilder().
		WithGrpc().
		WithDb(db).
		WithReflectionManager(rpcManager).
		WithLogger(testutil.NewLogger(t)).
		Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s, rpcManager)

	config := &api.CreateLinkRequest{
		Name:        name,
		SourceStage: "",
		TargetStage: "stage-2",
	}

	err = s.CreateLink(config)
	assert.Assert(
		t,
		errdefs.IsInvalidArgument(err),
		"error is not InvalidArgument",
	)
	assert.Error(t, err, "empty source stage name")
}

func TestServer_CreateLink_TargetEmpty(t *testing.T) {
	const name = "link-name"
	rpcManager := &mockreflection.MockManager{Rpcs: sync.Map{}}
	db, err := badger.Open(
		badger.DefaultOptions("").WithInMemory(true),
	)
	assert.NilError(t, err, "db creation")
	defer db.Close()
	s, err := NewBuilder().
		WithGrpc().
		WithDb(db).
		WithReflectionManager(rpcManager).
		WithLogger(testutil.NewLogger(t)).
		Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s, rpcManager)

	config := &api.CreateLinkRequest{
		Name:        name,
		SourceStage: "stage-2",
		TargetStage: "",
	}

	err = s.CreateLink(config)
	assert.Assert(
		t,
		errdefs.IsInvalidArgument(err),
		"error is not InvalidArgument",
	)
	assert.Error(t, err, "empty target stage name")
}

func TestServer_CreateLink_EqualSourceAndTarget(t *testing.T) {
	const name = "link-name"
	rpcManager := &mockreflection.MockManager{Rpcs: sync.Map{}}
	db, err := badger.Open(
		badger.DefaultOptions("").WithInMemory(true),
	)
	assert.NilError(t, err, "db creation")
	defer db.Close()
	s, err := NewBuilder().
		WithGrpc().
		WithDb(db).
		WithReflectionManager(rpcManager).
		WithLogger(testutil.NewLogger(t)).
		Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s, rpcManager)

	config := &api.CreateLinkRequest{
		Name:        name,
		SourceStage: "stage-1",
		TargetStage: "stage-1",
	}

	err = s.CreateLink(config)
	assert.Assert(t, errdefs.IsInvalidArgument(err), "error is not Invalid Arg")
	assert.Error(t, err, "source and target stages are equal")
}

func TestServer_CreateLink_SourceNotFound(t *testing.T) {
	const name = "link-name"
	rpcManager := &mockreflection.MockManager{Rpcs: sync.Map{}}
	db, err := badger.Open(
		badger.DefaultOptions("").WithInMemory(true),
	)
	assert.NilError(t, err, "db creation")
	defer db.Close()
	s, err := NewBuilder().
		WithGrpc().
		WithDb(db).
		WithReflectionManager(rpcManager).
		WithLogger(testutil.NewLogger(t)).
		Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s, rpcManager)

	config := &api.CreateLinkRequest{
		Name:        name,
		SourceStage: "stage-4",
		TargetStage: "stage-2",
	}

	err = s.CreateLink(config)
	assert.Assert(t, errdefs.IsNotFound(err), "error is not NotFound")
	expectedMsg := "source stage 'stage-4' not found"
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateLink_TargetNotFound(t *testing.T) {
	const name = "link-name"
	rpcManager := &mockreflection.MockManager{Rpcs: sync.Map{}}
	db, err := badger.Open(
		badger.DefaultOptions("").WithInMemory(true),
	)
	assert.NilError(t, err, "db creation")
	defer db.Close()
	s, err := NewBuilder().
		WithGrpc().
		WithDb(db).
		WithReflectionManager(rpcManager).
		WithLogger(testutil.NewLogger(t)).
		Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s, rpcManager)

	config := &api.CreateLinkRequest{
		Name:        name,
		SourceStage: "stage-1",
		TargetStage: "stage-4",
	}

	err = s.CreateLink(config)
	assert.Assert(t, errdefs.IsNotFound(err), "error is not NotFound")
	expectedMsg := "target stage 'stage-4' not found"
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateLink_AlreadyExists(t *testing.T) {
	var err error
	const name = "link-name"
	rpcManager := &mockreflection.MockManager{Rpcs: sync.Map{}}
	db, err := badger.Open(
		badger.DefaultOptions("").WithInMemory(true),
	)
	assert.NilError(t, err, "db creation")
	defer db.Close()
	s, err := NewBuilder().
		WithGrpc().
		WithDb(db).
		WithReflectionManager(rpcManager).
		WithLogger(testutil.NewLogger(t)).
		Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s, rpcManager)

	config := &api.CreateLinkRequest{
		Name:        name,
		SourceStage: "stage-1",
		TargetStage: "stage-2",
	}

	err = s.CreateLink(config)
	assert.NilError(t, err, "first creation has an error")
	err = s.CreateLink(config)
	assert.Assert(t, errdefs.IsAlreadyExists(err), "error is not NotFound")
	expectedMsg := fmt.Sprintf("link '%v' already exists", name)
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateLink_UnknownSourceField(t *testing.T) {
	const name = "link-name"
	rpcManager := &mockreflection.MockManager{Rpcs: sync.Map{}}
	db, err := badger.Open(
		badger.DefaultOptions("").WithInMemory(true),
	)
	assert.NilError(t, err, "db creation")
	defer db.Close()
	s, err := NewBuilder().
		WithGrpc().
		WithDb(db).
		WithReflectionManager(rpcManager).
		WithLogger(testutil.NewLogger(t)).
		Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s, rpcManager)

	config := &api.CreateLinkRequest{
		Name:        name,
		SourceStage: "stage-1",
		SourceField: "unknown-field",
		TargetStage: "stage-2",
	}

	err = s.CreateLink(config)
	assert.Assert(t, errdefs.IsNotFound(err), "error is not NotFound")
	expectedMsg := fmt.Sprintf(
		"field with name unknown-field not found for message "+
			"pb.TestMessage1 for source stage in link %v",
		name,
	)
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateLink_UnknownTargetField(t *testing.T) {
	const name = "link-name"
	rpcManager := &mockreflection.MockManager{Rpcs: sync.Map{}}
	db, err := badger.Open(
		badger.DefaultOptions("").WithInMemory(true),
	)
	assert.NilError(t, err, "db creation")
	defer db.Close()
	s, err := NewBuilder().
		WithGrpc().
		WithDb(db).
		WithReflectionManager(rpcManager).
		WithLogger(testutil.NewLogger(t)).
		Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s, rpcManager)

	config := &api.CreateLinkRequest{
		Name:        name,
		SourceStage: "stage-1",
		TargetStage: "stage-2",
		TargetField: "unknown-field",
	}

	err = s.CreateLink(config)
	assert.Assert(t, errdefs.IsNotFound(err), "error is not NotFound")
	expectedMsg := fmt.Sprintf(
		"field with name unknown-field not found for message "+
			"pb.TestMessageDiffNames for target stage in link %v",
		name,
	)
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateLink_IncompatibleMessages(t *testing.T) {
	const name = "link-name"
	rpcManager := &mockreflection.MockManager{Rpcs: sync.Map{}}
	db, err := badger.Open(
		badger.DefaultOptions("").WithInMemory(true),
	)
	assert.NilError(t, err, "db creation")
	defer db.Close()
	s, err := NewBuilder().
		WithGrpc().
		WithDb(db).
		WithReflectionManager(rpcManager).
		WithLogger(testutil.NewLogger(t)).
		Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s, rpcManager)

	config := &api.CreateLinkRequest{
		Name:        name,
		SourceStage: "stage-1",
		TargetStage: "stage-3",
	}

	err = s.CreateLink(config)
	assert.Assert(
		t,
		errdefs.IsInvalidArgument(err),
		"error is not IsInvalidArgument",
	)
	expectedMsg := fmt.Sprintf(
		"incompatible message types between source output pb.TestMessage1 "+
			"and target input pb.TestWrongOuterFieldType in link %v",
		name,
	)
	assert.Error(t, err, expectedMsg)
}

// populateForLinks creates three stages. The first two are compatible, but the
// third is not.
func populateForLinks(
	t *testing.T,
	s *Server,
	rpcManager *mockreflection.MockManager,
) {
	stage1 := mockStage(t, 1, pb.TestMessage1{}, pb.TestMessage1{}, rpcManager)

	stage2 := mockStage(
		t,
		2,
		pb.TestMessageDiffNames{},
		pb.TestMessageDiffNames{},
		rpcManager,
	)
	stage3 := mockStage(
		t,
		3,
		pb.TestWrongOuterFieldType{},
		pb.TestWrongOuterFieldType{},
		rpcManager,
	)

	stages := []*storage.Stage{stage1, stage2, stage3}
	err := s.db.Update(
		func(txn *badger.Txn) error {
			populateStages(t, s, txn, stages)
			return nil
		},
	)
	assert.NilError(t, err, "Populate stages for links")
}
