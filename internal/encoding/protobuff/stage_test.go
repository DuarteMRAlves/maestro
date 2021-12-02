package protobuff

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/stage"
	"gotest.tools/v3/assert"
	"testing"
)

func TestMarshalStage(t *testing.T) {
	tests := []*stage.Stage{
		{
			Name:    stageName,
			Asset:   stageAsset,
			Service: stageService,
			Method:  stageMethod,
		},
		{
			Name:    "",
			Asset:   "",
			Service: "",
			Method:  "",
		},
	}

	for _, s := range tests {
		testName := fmt.Sprintf("stage=%v", s)

		t.Run(
			testName, func(t *testing.T) {
				res := MarshalStage(s)
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
		},
		{
			Name:    "",
			Asset:   "",
			Service: "",
			Method:  "",
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

func TestUnmarshalStageError(t *testing.T) {
	res, err := UnmarshalStage(nil)
	assert.Assert(t, errdefs.IsInvalidArgument(err), "err type")
	assert.ErrorContains(t, err, "'p' is nil")
	assert.Assert(t, res == nil, "nil return value")
}

func assertStage(t *testing.T, expected *stage.Stage, actual *pb.Stage) {
	assert.Equal(t, expected.Name, actual.Name, "stage assetName")
	assert.Equal(t, expected.Asset, actual.Asset, "asset id")
	assert.Equal(t, expected.Service, actual.Service, "stage service")
	assert.Equal(t, expected.Method, actual.Method, "stage method")
}

func assertPbStage(t *testing.T, expected *pb.Stage, actual *stage.Stage) {
	assert.Equal(t, expected.Name, actual.Name, "stage assetName")
	assert.Equal(t, expected.Asset, actual.Asset, "asset id")
	assert.Equal(t, expected.Service, actual.Service, "stage service")
	assert.Equal(t, expected.Method, actual.Method, "stage method")
}
