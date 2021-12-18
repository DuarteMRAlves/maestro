package server

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/stage"
	"github.com/DuarteMRAlves/maestro/internal/testutil"
	"gotest.tools/v3/assert"
	"testing"
)

const stageName = "stage-name"

func TestServer_CreateStage(t *testing.T) {
	tests := []struct {
		name   string
		config *stage.Stage
	}{
		{
			name: "correct with nil asset, service, method and address",
			config: &stage.Stage{
				Name: stageName,
			},
		},
		{
			name: "correct with empty asset, method and address",
			config: &stage.Stage{
				Name:    stageName,
				Asset:   "",
				Service: "",
				Method:  "",
				Address: "",
			},
		},
		{
			name: "correct with asset, service, method and address",
			config: &stage.Stage{
				Name:    stageName,
				Asset:   assetNameForNum(0),
				Service: "ServiceName",
				Method:  "MethodName",
				Address: "Address",
			},
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
				assert.NilError(t, err, "build server")
				populateForStages(t, s)
				err = s.CreateStage(test.config)
				assert.NilError(t, err, "create stage error")
			})
	}
}

func TestServer_CreateStage_NilConfig(t *testing.T) {
	s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
	assert.NilError(t, err, "build server")
	populateForOrchestrations(t, s)

	err = s.CreateStage(nil)
	assert.Assert(
		t,
		errdefs.IsInvalidArgument(err),
		"error is not InvalidArgument")
	expectedMsg := "'config' is nil"
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateStage_InvalidName(t *testing.T) {
	tests := []struct {
		name   string
		config *stage.Stage
	}{
		{
			name: "empty name",
			config: &stage.Stage{
				Name:  "",
				Asset: assetNameForNum(0),
			},
		},
		{
			name: "invalid characters in name",
			config: &stage.Stage{
				Name:  "some@name",
				Asset: assetNameForNum(0),
			},
		},
		{
			name: "invalid character sequence",
			config: &stage.Stage{
				Name:  "other-/name",
				Asset: assetNameForNum(0),
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

				err = s.CreateStage(test.config)
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

func TestServer_CreateStage_AssetNotFound(t *testing.T) {
	s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
	assert.NilError(t, err, "build server")
	populateForStages(t, s)

	config := &stage.Stage{
		Name:  stageName,
		Asset: assetNameForNum(1),
	}

	err = s.CreateStage(config)
	assert.Assert(t, errdefs.IsNotFound(err), "error is not NotFound")
	expectedMsg := fmt.Sprintf("asset '%v' not found", assetNameForNum(1))
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateStage_AlreadyExists(t *testing.T) {
	var err error

	s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
	assert.NilError(t, err, "build server")
	populateForStages(t, s)

	config := &stage.Stage{
		Name:  stageName,
		Asset: assetNameForNum(0),
	}

	err = s.CreateStage(config)
	assert.NilError(t, err, "first creation has an error")
	err = s.CreateStage(config)
	assert.Assert(t, errdefs.IsAlreadyExists(err), "error is not NotFound")
	expectedMsg := fmt.Sprintf("stage '%v' already exists", stageName)
	assert.Error(t, err, expectedMsg)
}

func populateForStages(t *testing.T, s *Server) {
	assets := []*asset.Asset{
		assetForNum(0),
	}
	populateAssets(t, s, assets)
}
