package server

import (
	"fmt"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/testutil"
	"gotest.tools/v3/assert"
	"testing"
)

const assetName = "asset-name"

func TestServer_CreateAsset(t *testing.T) {
	tests := []struct {
		name   string
		config *apitypes.Asset
	}{
		{
			name: "correct nil image",
			config: &apitypes.Asset{
				Name: assetName,
			},
		},
		{
			name: "correct with empty image",
			config: &apitypes.Asset{
				Name:  assetName,
				Image: "",
			},
		},
		{
			name: "correct with image",
			config: &apitypes.Asset{
				Name:  assetName,
				Image: "image-name",
			},
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
				assert.NilError(t, err, "build server")
				err = s.CreateAsset(test.config)
				assert.NilError(t, err, "create asset error")
			})
	}
}

func TestServer_CreateAsset_NilConfig(t *testing.T) {
	s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
	assert.NilError(t, err, "build server")

	err = s.CreateAsset(nil)
	assert.Assert(
		t,
		errdefs.IsInvalidArgument(err),
		"error is not InvalidArgument")
	expectedMsg := "'config' is nil"
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateAsset_InvalidName(t *testing.T) {
	tests := []struct {
		name   string
		config *apitypes.Asset
	}{
		{
			name: "empty name",
			config: &apitypes.Asset{
				Name: "",
			},
		},
		{
			name: "invalid characters in name",
			config: &apitypes.Asset{
				Name: "%some-name%",
			},
		},
		{
			name: "invalid character sequence",
			config: &apitypes.Asset{
				Name: "invalid..name",
			},
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
				assert.NilError(t, err, "build server")

				err = s.CreateAsset(test.config)
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

func TestServer_CreateAsset_AlreadyExists(t *testing.T) {
	var err error

	s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
	assert.NilError(t, err, "build server")

	config := &apitypes.Asset{
		Name: assetName,
	}

	err = s.CreateAsset(config)
	assert.NilError(t, err, "first creation has an error")
	err = s.CreateAsset(config)
	assert.Assert(t, errdefs.IsAlreadyExists(err), "error is not NotFound")
	expectedMsg := fmt.Sprintf("asset '%v' already exists", assetName)
	assert.Error(t, err, expectedMsg)
}
