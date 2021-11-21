package protobuff

import (
	"errors"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/stage"
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
				assert.DeepEqual(t, nil, err, "unmarshal error")
				assertPbStage(t, s, res)
			})
	}
}

func TestUnmarshalStageError(t *testing.T) {
	res, err := UnmarshalStage(nil)
	assert.DeepEqual(
		t,
		errdefs.InternalWithError(errors.New("p is nil")),
		err,
		"unmarshal error")
	assert.IsNil(t, res, "nil return value")
}

func assertStage(t *testing.T, expected *stage.Stage, actual *pb.Stage) {
	assert.DeepEqual(t, expected.Name, actual.Name, "stage assetName")
	assert.DeepEqual(t, expected.Asset, actual.Asset, "asset id")
	assert.DeepEqual(t, expected.Service, actual.Service, "stage service")
	assert.DeepEqual(t, expected.Method, actual.Method, "stage method")
}

func assertPbStage(t *testing.T, expected *pb.Stage, actual *stage.Stage) {
	assert.DeepEqual(t, expected.Name, actual.Name, "stage assetName")
	assert.DeepEqual(t, expected.Asset, actual.Asset, "asset id")
	assert.DeepEqual(t, expected.Service, actual.Service, "stage service")
	assert.DeepEqual(t, expected.Method, actual.Method, "stage method")
}
