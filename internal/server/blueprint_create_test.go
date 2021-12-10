package server

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/blueprint"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/link"
	"github.com/DuarteMRAlves/maestro/internal/stage"
	"gotest.tools/v3/assert"
	"testing"
)

const bpName = "blueprint-name"

func TestServer_CreateBlueprint(t *testing.T) {
	tests := []struct {
		name   string
		config *blueprint.Blueprint
	}{
		{
			name: "correct with nil links",
			config: &blueprint.Blueprint{
				Name:  bpName,
				Links: []string{},
			},
		},
		{
			name: "correct with empty links",
			config: &blueprint.Blueprint{
				Name:  bpName,
				Links: []string{},
			},
		},
		{
			name: "correct with one link",
			config: &blueprint.Blueprint{
				Name:  bpName,
				Links: []string{linkNameForNum(0)},
			},
		},
		{
			name: "correct with multiple links",
			config: &blueprint.Blueprint{
				Name: bpName,
				Links: []string{
					linkNameForNum(0),
					linkNameForNum(2),
					linkNameForNum(1),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				s := NewBuilder().WithGrpc().Build()
				populateForBlueprints(t, s)
				err := s.CreateBlueprint(test.config)
				assert.NilError(t, err, "create blueprint error")
			})
	}
}

func TestServer_CreateBlueprint_NilConfig(t *testing.T) {
	s := NewBuilder().WithGrpc().Build()
	populateForBlueprints(t, s)

	err := s.CreateBlueprint(nil)
	assert.Assert(
		t,
		errdefs.IsInvalidArgument(err),
		"error is not InvalidArgument")
	expectedMsg := "'config' is nil"
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateBlueprint_InvalidName(t *testing.T) {
	tests := []struct {
		name   string
		config *blueprint.Blueprint
	}{
		{
			name: "empty name",
			config: &blueprint.Blueprint{
				Name: "",
				Links: []string{
					linkNameForNum(0),
					linkNameForNum(1),
					linkNameForNum(2),
				},
			},
		},
		{
			name: "invalid characters in name",
			config: &blueprint.Blueprint{
				Name: "?blueprint-name",
				Links: []string{
					linkNameForNum(0),
					linkNameForNum(1),
					linkNameForNum(2),
				},
			},
		},
		{
			name: "invalid character sequence",
			config: &blueprint.Blueprint{
				Name: "invalid//name",
				Links: []string{
					linkNameForNum(0),
					linkNameForNum(1),
					linkNameForNum(2),
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				s := NewBuilder().WithGrpc().Build()
				populateForBlueprints(t, s)

				err := s.CreateBlueprint(test.config)
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

func TestServer_CreateBlueprint_LinkNotFound(t *testing.T) {
	s := NewBuilder().WithGrpc().Build()
	populateForBlueprints(t, s)

	config := &blueprint.Blueprint{
		Name: bpName,
		Links: []string{
			linkNameForNum(0),
			// This link does not exist
			linkNameForNum(3),
			linkNameForNum(2),
		},
	}

	err := s.CreateBlueprint(config)
	assert.Assert(t, errdefs.IsNotFound(err), "error is not NotFound")
	expectedMsg := fmt.Sprintf("link '%v' not found", linkNameForNum(3))
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateBlueprint_AlreadyExists(t *testing.T) {
	var err error

	s := NewBuilder().WithGrpc().Build()
	populateForBlueprints(t, s)

	config := &blueprint.Blueprint{
		Name: bpName,
		Links: []string{
			linkNameForNum(0),
			linkNameForNum(1),
			linkNameForNum(2),
		},
	}

	err = s.CreateBlueprint(config)
	assert.NilError(t, err, "first creation has an error")
	err = s.CreateBlueprint(config)
	assert.Assert(t, errdefs.IsAlreadyExists(err), "error is not NotFound")
	expectedMsg := fmt.Sprintf("blueprint '%v' already exists", bpName)
	assert.Error(t, err, expectedMsg)
}

func populateForBlueprints(t *testing.T, s *Server) {
	assets := []*asset.Asset{
		assetForNum(0),
		assetForNum(1),
		assetForNum(2),
		assetForNum(3),
	}
	stages := []*stage.Stage{
		stageForNum(0),
		stageForNum(1),
		stageForNum(2),
		stageForNum(3),
	}
	links := []*link.Link{linkForNum(0), linkForNum(1), linkForNum(2)}
	populateAssets(t, s, assets)
	populateStages(t, s, stages)
	populateLinks(t, s, links)
}
