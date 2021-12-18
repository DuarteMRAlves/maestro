package server

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/link"
	"github.com/DuarteMRAlves/maestro/internal/stage"
	"github.com/DuarteMRAlves/maestro/internal/testutil"
	"gotest.tools/v3/assert"
	"testing"
)

const linkName = "link-name"

func TestServer_CreateLink(t *testing.T) {
	tests := []struct {
		name   string
		config *link.Link
	}{
		{
			name: "correct with nil fields",
			config: &link.Link{
				Name:        linkName,
				SourceStage: stageNameForNum(0),
				TargetStage: stageNameForNum(1),
			},
		},
		{
			name: "correct with empty fields",
			config: &link.Link{
				Name:        linkName,
				SourceStage: stageNameForNum(0),
				SourceField: "",
				TargetStage: stageNameForNum(1),
				TargetField: "",
			},
		},
		{
			name: "correct with fields",
			config: &link.Link{
				Name:        linkName,
				SourceStage: stageNameForNum(0),
				SourceField: "SourceField",
				TargetStage: stageNameForNum(1),
				TargetField: "TargetField",
			},
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
				assert.NilError(t, err, "build server")

				populateForLinks(t, s)
				err = s.CreateLink(test.config)
				assert.NilError(t, err, "create link error")
			})
	}
}

func TestServer_CreateLink_NilConfig(t *testing.T) {
	s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
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
		config *link.Link
	}{
		{
			name: "empty name",
			config: &link.Link{
				Name:        "",
				SourceStage: stageNameForNum(0),
				TargetStage: stageNameForNum(1),
			},
		},
		{
			name: "invalid characters in name",
			config: &link.Link{
				Name:        "some'character",
				SourceStage: stageNameForNum(0),
				TargetStage: stageNameForNum(1),
			},
		},
		{
			name: "invalid character sequence",
			config: &link.Link{
				Name:        "//invalid-name",
				SourceStage: stageNameForNum(0),
				TargetStage: stageNameForNum(1),
			},
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
				assert.NilError(t, err, "build server")
				populateForOrchestrations(t, s)

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
	s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s)

	config := &link.Link{
		Name:        linkName,
		SourceStage: "",
		TargetStage: stageNameForNum(1),
	}

	err = s.CreateLink(config)
	assert.Assert(
		t,
		errdefs.IsInvalidArgument(err),
		"error is not InvalidArgument")
	assert.Error(t, err, "empty source stage name")
}

func TestServer_CreateLink_TargetEmpty(t *testing.T) {
	s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s)

	config := &link.Link{
		Name:        linkName,
		SourceStage: stageNameForNum(1),
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
	s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s)

	config := &link.Link{
		Name:        linkName,
		SourceStage: stageNameForNum(0),
		TargetStage: stageNameForNum(0),
	}

	err = s.CreateLink(config)
	assert.Assert(t, errdefs.IsInvalidArgument(err), "error is not NotFound")
	assert.Error(t, err, "source and target stages are equal")
}

func TestServer_CreateLink_SourceNotFound(t *testing.T) {
	s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s)

	config := &link.Link{
		Name:        linkName,
		SourceStage: stageNameForNum(3),
		TargetStage: stageNameForNum(1),
	}

	err = s.CreateLink(config)
	assert.Assert(t, errdefs.IsNotFound(err), "error is not NotFound")
	expectedMsg := fmt.Sprintf(
		"source stage '%v' not found",
		stageNameForNum(3))
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateLink_TargetNotFound(t *testing.T) {
	s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s)

	config := &link.Link{
		Name:        linkName,
		SourceStage: stageNameForNum(1),
		TargetStage: stageNameForNum(3),
	}

	err = s.CreateLink(config)
	assert.Assert(t, errdefs.IsNotFound(err), "error is not NotFound")
	expectedMsg := fmt.Sprintf(
		"target stage '%v' not found",
		stageNameForNum(3))
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateLink_AlreadyExists(t *testing.T) {
	var err error

	s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
	assert.NilError(t, err, "build server")
	populateForLinks(t, s)

	config := &link.Link{
		Name:        linkName,
		SourceStage: stageNameForNum(0),
		TargetStage: stageNameForNum(1),
	}

	err = s.CreateLink(config)
	assert.NilError(t, err, "first creation has an error")
	err = s.CreateLink(config)
	assert.Assert(t, errdefs.IsAlreadyExists(err), "error is not NotFound")
	expectedMsg := fmt.Sprintf("link '%v' already exists", linkName)
	assert.Error(t, err, expectedMsg)
}

func populateForLinks(t *testing.T, s *Server) {
	assets := []*asset.Asset{
		assetForNum(0),
		assetForNum(1),
	}
	stages := []*stage.Stage{
		stageForNum(0),
		stageForNum(1),
	}
	populateAssets(t, s, assets)
	populateStages(t, s, stages)
}
