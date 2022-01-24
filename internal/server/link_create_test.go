package server

import (
	"fmt"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/orchestration"
	"github.com/DuarteMRAlves/maestro/internal/testutil"
	"github.com/DuarteMRAlves/maestro/tests/pb"
	"github.com/dgraph-io/badger/v3"
	"gotest.tools/v3/assert"
	"testing"
)

func TestServer_CreateLink(t *testing.T) {
	const name = "link-name"
	tests := []struct {
		name   string
		config *apitypes.Link
	}{
		{
			name: "correct with nil fields",
			config: &apitypes.Link{
				Name:        name,
				SourceStage: "stage-1",
				TargetStage: "stage-2",
			},
		},
		{
			name: "correct with empty fields",
			config: &apitypes.Link{
				Name:        name,
				SourceStage: "stage-1",
				SourceField: "",
				TargetStage: "stage-2",
				TargetField: "",
			},
		},
		{
			name: "correct with fields",
			config: &apitypes.Link{
				Name:        name,
				SourceStage: "stage-1",
				SourceField: "field4",
				TargetStage: "stage-2",
				TargetField: "fieldName4",
			},
		},
		{
			name: "incompatible outer but compatible inner",
			config: &apitypes.Link{
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
				db, err := badger.Open(
					badger.DefaultOptions("").WithInMemory(true))
				assert.NilError(t, err, "db creation")
				defer db.Close()
				s, err := NewBuilder().
					WithGrpc().
					WithDb(db).
					WithLogger(testutil.NewLogger(t)).
					Build()
				assert.NilError(t, err, "build server")

				populateForLinks(t, s)
				err = s.CreateLink(test.config)
				assert.NilError(t, err, "create link error")
			})
	}
}

func TestServer_CreateLink_NilConfig(t *testing.T) {
	db, err := badger.Open(
		badger.DefaultOptions("").WithInMemory(true))
	assert.NilError(t, err, "db creation")
	defer db.Close()
	s, err := NewBuilder().
		WithGrpc().
		WithDb(db).
		WithLogger(testutil.NewLogger(t)).
		Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s)

	err = s.CreateLink(nil)
	assert.Assert(
		t,
		errdefs.IsInvalidArgument(err),
		"error is not InvalidArgument")
	expectedMsg := "'config' is nil"
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateLink_InvalidName(t *testing.T) {
	tests := []struct {
		name   string
		config *apitypes.Link
	}{
		{
			name: "empty name",
			config: &apitypes.Link{
				Name:        "",
				SourceStage: "stage-1",
				TargetStage: "stage-2",
			},
		},
		{
			name: "invalid characters in name",
			config: &apitypes.Link{
				Name:        "some'character",
				SourceStage: "stage-1",
				TargetStage: "stage-2",
			},
		},
		{
			name: "invalid character sequence",
			config: &apitypes.Link{
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
					badger.DefaultOptions("").WithInMemory(true))
				assert.NilError(t, err, "db creation")
				defer db.Close()
				s, err := NewBuilder().
					WithGrpc().
					WithDb(db).
					WithLogger(testutil.NewLogger(t)).
					Build()
				assert.NilError(t, err, "build server")

				err = s.CreateLink(test.config)
				assert.Assert(
					t,
					errdefs.IsInvalidArgument(err),
					"error is not InvalidArgument")
				expectedMsg := fmt.Sprintf(
					"invalid name '%v'",
					test.config.Name)
				assert.Error(t, err, expectedMsg)
			})
	}
}

func TestServer_CreateLink_SourceEmpty(t *testing.T) {
	const name = "link-name"
	db, err := badger.Open(
		badger.DefaultOptions("").WithInMemory(true))
	assert.NilError(t, err, "db creation")
	defer db.Close()
	s, err := NewBuilder().
		WithGrpc().
		WithDb(db).
		WithLogger(testutil.NewLogger(t)).
		Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s)

	config := &apitypes.Link{
		Name:        name,
		SourceStage: "",
		TargetStage: "stage-2",
	}

	err = s.CreateLink(config)
	assert.Assert(
		t,
		errdefs.IsInvalidArgument(err),
		"error is not InvalidArgument")
	assert.Error(t, err, "empty source stage name")
}

func TestServer_CreateLink_TargetEmpty(t *testing.T) {
	const name = "link-name"
	db, err := badger.Open(
		badger.DefaultOptions("").WithInMemory(true))
	assert.NilError(t, err, "db creation")
	defer db.Close()
	s, err := NewBuilder().
		WithGrpc().
		WithDb(db).
		WithLogger(testutil.NewLogger(t)).
		Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s)

	config := &apitypes.Link{
		Name:        name,
		SourceStage: "stage-2",
		TargetStage: "",
	}

	err = s.CreateLink(config)
	assert.Assert(
		t,
		errdefs.IsInvalidArgument(err),
		"error is not InvalidArgument")
	assert.Error(t, err, "empty target stage name")
}

func TestServer_CreateLink_EqualSourceAndTarget(t *testing.T) {
	const name = "link-name"
	db, err := badger.Open(
		badger.DefaultOptions("").WithInMemory(true))
	assert.NilError(t, err, "db creation")
	defer db.Close()
	s, err := NewBuilder().
		WithGrpc().
		WithDb(db).
		WithLogger(testutil.NewLogger(t)).
		Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s)

	config := &apitypes.Link{
		Name:        name,
		SourceStage: "stage-1",
		TargetStage: "stage-1",
	}

	err = s.CreateLink(config)
	assert.Assert(t, errdefs.IsInvalidArgument(err), "error is not NotFound")
	assert.Error(t, err, "source and target stages are equal")
}

func TestServer_CreateLink_SourceNotFound(t *testing.T) {
	const name = "link-name"
	db, err := badger.Open(
		badger.DefaultOptions("").WithInMemory(true))
	assert.NilError(t, err, "db creation")
	defer db.Close()
	s, err := NewBuilder().
		WithGrpc().
		WithDb(db).
		WithLogger(testutil.NewLogger(t)).
		Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s)

	config := &apitypes.Link{
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
	db, err := badger.Open(
		badger.DefaultOptions("").WithInMemory(true))
	assert.NilError(t, err, "db creation")
	defer db.Close()
	s, err := NewBuilder().
		WithGrpc().
		WithDb(db).
		WithLogger(testutil.NewLogger(t)).
		Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s)

	config := &apitypes.Link{
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
	db, err := badger.Open(
		badger.DefaultOptions("").WithInMemory(true))
	assert.NilError(t, err, "db creation")
	defer db.Close()
	s, err := NewBuilder().
		WithGrpc().
		WithDb(db).
		WithLogger(testutil.NewLogger(t)).
		Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s)

	config := &apitypes.Link{
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
	db, err := badger.Open(
		badger.DefaultOptions("").WithInMemory(true))
	assert.NilError(t, err, "db creation")
	defer db.Close()
	s, err := NewBuilder().
		WithGrpc().
		WithDb(db).
		WithLogger(testutil.NewLogger(t)).
		Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s)

	config := &apitypes.Link{
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
		name)
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateLink_UnknownTargetField(t *testing.T) {
	const name = "link-name"
	db, err := badger.Open(
		badger.DefaultOptions("").WithInMemory(true))
	assert.NilError(t, err, "db creation")
	defer db.Close()
	s, err := NewBuilder().
		WithGrpc().
		WithDb(db).
		WithLogger(testutil.NewLogger(t)).
		Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s)

	config := &apitypes.Link{
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
		name)
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateLink_IncompatibleMessages(t *testing.T) {
	const name = "link-name"
	db, err := badger.Open(
		badger.DefaultOptions("").WithInMemory(true))
	assert.NilError(t, err, "db creation")
	defer db.Close()
	s, err := NewBuilder().
		WithGrpc().
		WithDb(db).
		WithLogger(testutil.NewLogger(t)).
		Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s)

	config := &apitypes.Link{
		Name:        name,
		SourceStage: "stage-1",
		TargetStage: "stage-3",
	}

	err = s.CreateLink(config)
	assert.Assert(t, errdefs.IsInvalidArgument(err), "error is not IsInvalidArgument")
	expectedMsg := fmt.Sprintf(
		"incompatible message types between source output pb.TestMessage1 "+
			"and target input pb.TestWrongOuterFieldType in link %v",
		name)
	assert.Error(t, err, expectedMsg)
}

// populateForLinks creates three stages. The first two are compatible, but the
// third is not.
func populateForLinks(t *testing.T, s *Server) {
	stage1 := mockStage(t, 1, pb.TestMessage1{}, pb.TestMessage1{})
	stage2 := mockStage(
		t,
		2,
		pb.TestMessageDiffNames{},
		pb.TestMessageDiffNames{})
	stage3 := mockStage(
		t,
		3,
		pb.TestWrongOuterFieldType{},
		pb.TestWrongOuterFieldType{})

	stages := []*orchestration.Stage{stage1, stage2, stage3}

	populateStages(t, s, stages)
}
