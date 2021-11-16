package protobuff

import (
	"errors"
	"fmt"
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/assert"
	"github.com/DuarteMRAlves/maestro/internal/blueprint"
	"github.com/DuarteMRAlves/maestro/internal/identifier"
	"testing"
)

func TestMarshalStage(t *testing.T) {
	id, err := identifier.Rand(5)
	assert.IsNil(t, err, "create stage id")
	assetId, err := identifier.Rand(5)
	assert.IsNil(t, err, "create asset id")

	tests := []*blueprint.Stage{
		{
			Id:      id,
			Name:    stageName,
			AssetId: assetId,
			Service: stageService,
			Method:  stageMethod,
		},
		{
			Id:      identifier.Empty(),
			Name:    "",
			AssetId: identifier.Empty(),
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
			Id:      &pb.Id{Val: stageId},
			Name:    stageName,
			AssetId: &pb.Id{Val: stageAssetId},
			Service: stageService,
			Method:  stageMethod,
		},
		{
			Id:      &pb.Id{Val: ""},
			Name:    "",
			AssetId: &pb.Id{Val: ""},
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
	assert.DeepEqual(t, errors.New("p is nil"), err, "unmarshal error")
	assert.IsNil(t, res, "nil return value")
}

func assertStage(t *testing.T, expected *blueprint.Stage, actual *pb.Stage) {
	assert.DeepEqual(t, expected.Id.Val, actual.Id.Val, "stage id")
	assert.DeepEqual(t, expected.Name, actual.Name, "stage assetName")
	assert.DeepEqual(t, expected.AssetId.Val, actual.AssetId.Val, "asset id")
	assert.DeepEqual(t, expected.Service, actual.Service, "stage service")
	assert.DeepEqual(t, expected.Method, actual.Method, "stage method")
}

func assertPbStage(t *testing.T, expected *pb.Stage, actual *blueprint.Stage) {
	fmt.Println(expected)
	fmt.Println(actual)
	assert.DeepEqual(t, expected.Id.Val, actual.Id.Val, "stage id")
	assert.DeepEqual(t, expected.Name, actual.Name, "stage assetName")
	assert.DeepEqual(t, expected.AssetId.Val, actual.AssetId.Val, "asset id")
	assert.DeepEqual(t, expected.Service, actual.Service, "stage service")
	assert.DeepEqual(t, expected.Method, actual.Method, "stage method")
}
