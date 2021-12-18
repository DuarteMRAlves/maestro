package server

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/link"
	"github.com/DuarteMRAlves/maestro/internal/orchestration"
	"github.com/DuarteMRAlves/maestro/internal/stage"
	"github.com/DuarteMRAlves/maestro/internal/testutil"
	"gotest.tools/v3/assert"
	"testing"
)

const oName = "orchestration-name"

func TestServer_CreateOrchestration(t *testing.T) {
	tests := []struct {
		name   string
		config *orchestration.Orchestration
	}{
		{
			name: "correct with nil links",
			config: &orchestration.Orchestration{
				Name:  oName,
				Links: []string{},
			},
		},
		{
			name: "correct with empty links",
			config: &orchestration.Orchestration{
				Name:  oName,
				Links: []string{},
			},
		},
		{
			name: "correct with one link",
			config: &orchestration.Orchestration{
				Name:  oName,
				Links: []string{linkNameForNum(0)},
			},
		},
		{
			name: "correct with multiple links",
			config: &orchestration.Orchestration{
				Name: oName,
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
				s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
				assert.NilError(t, err, "build server")

				populateForOrchestrations(t, s)
				err = s.CreateOrchestration(test.config)
				assert.NilError(t, err, "create orchestration error")
			})
	}
}

func TestServer_CreateOrchestration_NilConfig(t *testing.T) {
	s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
	assert.NilError(t, err, "build server")
	populateForOrchestrations(t, s)

	err = s.CreateOrchestration(nil)
	assert.Assert(
		t,
		errdefs.IsInvalidArgument(err),
		"error is not InvalidArgument")
	expectedMsg := "'config' is nil"
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateOrchestration_InvalidName(t *testing.T) {
	tests := []struct {
		name   string
		config *orchestration.Orchestration
	}{
		{
			name: "empty name",
			config: &orchestration.Orchestration{
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
			config: &orchestration.Orchestration{
				Name: "?orchestration-name",
				Links: []string{
					linkNameForNum(0),
					linkNameForNum(1),
					linkNameForNum(2),
				},
			},
		},
		{
			name: "invalid character sequence",
			config: &orchestration.Orchestration{
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
				s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
				assert.NilError(t, err, "build server")
				populateForOrchestrations(t, s)

				err = s.CreateOrchestration(test.config)
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

func TestServer_CreateOrchestration_LinkNotFound(t *testing.T) {
	s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
	assert.NilError(t, err, "build server")
	populateForOrchestrations(t, s)

	config := &orchestration.Orchestration{
		Name: oName,
		Links: []string{
			linkNameForNum(0),
			// This link does not exist
			linkNameForNum(3),
			linkNameForNum(2),
		},
	}

	err = s.CreateOrchestration(config)
	assert.Assert(t, errdefs.IsNotFound(err), "error is not NotFound")
	expectedMsg := fmt.Sprintf("link '%v' not found", linkNameForNum(3))
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateOrchestration_AlreadyExists(t *testing.T) {
	var err error

	s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
	assert.NilError(t, err, "build server")
	populateForOrchestrations(t, s)

	config := &orchestration.Orchestration{
		Name: oName,
		Links: []string{
			linkNameForNum(0),
			linkNameForNum(1),
			linkNameForNum(2),
		},
	}

	err = s.CreateOrchestration(config)
	assert.NilError(t, err, "first creation has an error")
	err = s.CreateOrchestration(config)
	assert.Assert(t, errdefs.IsAlreadyExists(err), "error is not NotFound")
	expectedMsg := fmt.Sprintf("orchestration '%v' already exists", oName)
	assert.Error(t, err, expectedMsg)
}

func populateForOrchestrations(t *testing.T, s *Server) {
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
