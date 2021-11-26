package protobuff

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/stage"
	testing2 "github.com/DuarteMRAlves/maestro/internal/testing"
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
				testing2.DeepEqual(t, nil, err, "unmarshal error")
				assertPbStage(t, s, res)
			})
	}
}

func TestUnmarshalStageError(t *testing.T) {
	res, err := UnmarshalStage(nil)
	expectedErr := errdefs.InvalidArgumentWithMsg("'p' is nil")
	testing2.DeepEqual(t, expectedErr, err, "unmarshal error")
	testing2.IsNil(t, res, "nil return value")
}

func assertStage(t *testing.T, expected *stage.Stage, actual *pb.Stage) {
	testing2.DeepEqual(t, expected.Name, actual.Name, "stage assetName")
	testing2.DeepEqual(t, expected.Asset, actual.Asset, "asset id")
	testing2.DeepEqual(t, expected.Service, actual.Service, "stage service")
	testing2.DeepEqual(t, expected.Method, actual.Method, "stage method")
}

func assertPbStage(t *testing.T, expected *pb.Stage, actual *stage.Stage) {
	testing2.DeepEqual(t, expected.Name, actual.Name, "stage assetName")
	testing2.DeepEqual(t, expected.Asset, actual.Asset, "asset id")
	testing2.DeepEqual(t, expected.Service, actual.Service, "stage service")
	testing2.DeepEqual(t, expected.Method, actual.Method, "stage method")
}
