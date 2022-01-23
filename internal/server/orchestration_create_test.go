package server

import (
	"fmt"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/testutil"
	"gotest.tools/v3/assert"
	"testing"
)

func TestServer_CreateOrchestration(t *testing.T) {
	const name = "orchestration-name"
	tests := []struct {
		name   string
		config *apitypes.Orchestration
	}{
		{
			name: "correct with name",
			config: &apitypes.Orchestration{
				Name: name,
			},
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
				assert.NilError(t, err, "build server")

				err = s.CreateOrchestration(test.config)
				assert.NilError(t, err, "create orchestration error")
			})
	}
}

func TestServer_CreateOrchestration_NilConfig(t *testing.T) {
	s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
	assert.NilError(t, err, "build server")

	err = s.CreateOrchestration(nil)
	assert.Assert(
		t,
		errdefs.IsInvalidArgument(err),
		"error is not InvalidArgument")
	expectedMsg := "'cfg' is nil"
	assert.Error(t, err, expectedMsg)
}

func TestServer_CreateOrchestration_InvalidName(t *testing.T) {
	tests := []struct {
		name   string
		config *apitypes.Orchestration
	}{
		{
			name: "empty name",
			config: &apitypes.Orchestration{
				Name: "",
				Links: []apitypes.LinkName{
					testutil.LinkNameForNum(0),
					testutil.LinkNameForNum(1),
					testutil.LinkNameForNum(2),
				},
			},
		},
		{
			name: "invalid characters in name",
			config: &apitypes.Orchestration{
				Name: "?orchestration-name",
				Links: []apitypes.LinkName{
					testutil.LinkNameForNum(0),
					testutil.LinkNameForNum(1),
					testutil.LinkNameForNum(2),
				},
			},
		},
		{
			name: "invalid character sequence",
			config: &apitypes.Orchestration{
				Name: "invalid//name",
				Links: []apitypes.LinkName{
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
				s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
				assert.NilError(t, err, "build server")

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

func TestServer_CreateOrchestration_AlreadyExists(t *testing.T) {
	var err error
	const name = "orchestration-name"
	s, err := NewBuilder().WithGrpc().WithLogger(testutil.NewLogger(t)).Build()
	assert.NilError(t, err, "build server")

	config := &apitypes.Orchestration{
		Name: name,
		Links: []apitypes.LinkName{
			testutil.LinkNameForNum(0),
			testutil.LinkNameForNum(1),
		},
	}

	err = s.CreateOrchestration(config)
	assert.NilError(t, err, "first creation has an error")
	err = s.CreateOrchestration(config)
	assert.Assert(t, errdefs.IsAlreadyExists(err), "error is not NotFound")
	expectedMsg := fmt.Sprintf("orchestration '%v' already exists", name)
	assert.Error(t, err, expectedMsg)
}
