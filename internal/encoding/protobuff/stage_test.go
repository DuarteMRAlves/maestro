package protobuff

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"gotest.tools/v3/assert"
	"testing"
)

func TestMarshalStage(t *testing.T) {
	const (
		stageName    api.StageName = "Stage Name"
		stagePhase                 = api.StagePending
		stageAsset                 = "Stage Asset"
		stageService               = "stageService"
		stageRpc                   = "stageRpc"
		stageAddress               = "stageAddress"
	)
	tests := []*api.Stage{
		{
			Name:    stageName,
			Phase:   stagePhase,
			Asset:   stageAsset,
			Service: stageService,
			Rpc:     stageRpc,
			Address: stageAddress,
		},
		{
			Name:    "",
			Phase:   "",
			Asset:   "",
			Service: "",
			Rpc:     "",
			Address: "",
		},
	}

	for _, s := range tests {
		testName := fmt.Sprintf("stage=%v", s)

		t.Run(
			testName, func(t *testing.T) {
				res, err := MarshalStage(s)
				assert.NilError(t, err, "marshal error")
				assertStage(t, s, res)
			},
		)
	}
}

func TestUnmarshalStageCorrect(t *testing.T) {
	const (
		stageName    = "Stage Name"
		stagePhase   = "Running"
		stageAsset   = "Stage Asset"
		stageService = "stageService"
		stageRpc     = "stageRpc"
		stageAddress = "stageAddress"
	)
	tests := []*pb.Stage{
		{
			Name:    stageName,
			Phase:   stagePhase,
			Asset:   stageAsset,
			Service: stageService,
			Rpc:     stageRpc,
			Address: stageAddress,
		},
		{
			Name:    "",
			Phase:   "",
			Asset:   "",
			Service: "",
			Rpc:     "",
			Address: "",
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
			},
		)
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

func assertStage(t *testing.T, expected *api.Stage, actual *pb.Stage) {
	assert.Equal(t, string(expected.Name), actual.Name, "stage assetName")
	assert.Equal(t, string(expected.Phase), actual.Phase, "stage phase")
	assert.Equal(t, string(expected.Asset), actual.Asset, "asset id")
	assert.Equal(t, expected.Service, actual.Service, "stage service")
	assert.Equal(t, expected.Rpc, actual.Rpc, "stage rpc")
	assert.Equal(t, expected.Address, actual.Address, "stage address")
}

func assertPbStage(t *testing.T, expected *pb.Stage, actual *api.Stage) {
	assert.Equal(t, expected.Name, string(actual.Name), "stage assetName")
	assert.Equal(t, expected.Phase, string(actual.Phase), "stage phase")
	assert.Equal(t, expected.Asset, string(actual.Asset), "asset id")
	assert.Equal(t, expected.Service, actual.Service, "stage service")
	assert.Equal(t, expected.Rpc, actual.Rpc, "stage rpc")
	assert.Equal(t, expected.Address, actual.Address, "stage address")
}
