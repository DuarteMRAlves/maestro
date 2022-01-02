package protobuff

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"gotest.tools/v3/assert"
	"testing"
)

func TestMarshalStage(t *testing.T) {
	tests := []*types.Stage{
		{
			Name:    stageName,
			Asset:   stageAsset,
			Service: stageService,
			Method:  stageMethod,
			Address: stageAddress,
			Host:    stageHost,
			Port:    stagePort,
		},
		{
			Name:    "",
			Asset:   "",
			Service: "",
			Method:  "",
			Address: "",
			Host:    "",
			Port:    0,
		},
	}

	for _, s := range tests {
		testName := fmt.Sprintf("stage=%v", s)

		t.Run(
			testName, func(t *testing.T) {
				res, err := MarshalStage(s)
				assert.NilError(t, err, "marshal error")
				assertStage(t, s, res)
			})
	}
}

func TestUnmarshalStageCorrect(t *testing.T) {
	tests := []*pb.Stage{
		{
			Name:    stageName,
			Asset:   stageAsset,
			Service: stageService,
			Method:  stageMethod,
			Address: stageAddress,
			Host:    stageHost,
			Port:    stagePort,
		},
		{
			Name:    "",
			Asset:   "",
			Service: "",
			Method:  "",
			Address: "",
			Host:    "",
			Port:    0,
		},
	}
	for _, s := range tests {
		testName := fmt.Sprintf("stage=%v", s)

		t.Run(
			testName,
			func(t *testing.T) {
				res, err := UnmarshalStage(s)
				assert.Equal(t, nil, err, "unmarshal error")
				assertPbStage(t, s, res)
			})
	}
}

func TestMarshalStageNil(t *testing.T) {
	res, err := MarshalStage(nil)
	assert.Assert(t, errdefs.IsInvalidArgument(err), "err type")
	assert.ErrorContains(t, err, "'s' is nil")
	assert.Assert(t, res == nil, "nil return value")
}

func TestUnmarshalStageNil(t *testing.T) {
	res, err := UnmarshalStage(nil)
	assert.Assert(t, errdefs.IsInvalidArgument(err), "err type")
	assert.ErrorContains(t, err, "'p' is nil")
	assert.Assert(t, res == nil, "nil return value")
}

func assertStage(t *testing.T, expected *types.Stage, actual *pb.Stage) {
	assert.Equal(t, expected.Name, actual.Name, "stage assetName")
	assert.Equal(t, expected.Asset, actual.Asset, "asset id")
	assert.Equal(t, expected.Service, actual.Service, "stage service")
	assert.Equal(t, expected.Method, actual.Method, "stage method")
	assert.Equal(t, expected.Address, actual.Address, "stage address")
	assert.Equal(t, expected.Host, actual.Host, "stage host")
	assert.Equal(t, expected.Port, actual.Port, "stage port")
}

func assertPbStage(t *testing.T, expected *pb.Stage, actual *types.Stage) {
	assert.Equal(t, expected.Name, actual.Name, "stage assetName")
	assert.Equal(t, expected.Asset, actual.Asset, "asset id")
	assert.Equal(t, expected.Service, actual.Service, "stage service")
	assert.Equal(t, expected.Method, actual.Method, "stage method")
	assert.Equal(t, expected.Address, actual.Address, "stage address")
	assert.Equal(t, expected.Host, actual.Host, "stage host")
	assert.Equal(t, expected.Port, actual.Port, "stage port")
}
