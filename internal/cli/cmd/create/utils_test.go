package create

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/testutil"
)

func equalCreateAssetRequest(
	expected *pb.CreateAssetRequest,
	actual *pb.CreateAssetRequest,
) bool {
	return expected.Name == actual.Name && expected.Image == actual.Image
}

func equalCreateStageRequest(
	expected *pb.CreateStageRequest,
	actual *pb.CreateStageRequest,
) bool {
	return expected.Name == actual.Name &&
		expected.Asset == actual.Asset &&
		expected.Service == actual.Service &&
		expected.Rpc == actual.Rpc &&
		expected.Address == actual.Address &&
		expected.Host == actual.Host &&
		expected.Port == actual.Port
}

func equalCreateLinkRequest(
	expected *pb.CreateLinkRequest,
	actual *pb.CreateLinkRequest,
) bool {
	return expected.Name == actual.Name &&
		expected.SourceStage == actual.SourceStage &&
		expected.SourceField == actual.SourceField &&
		expected.TargetStage == actual.TargetStage &&
		expected.TargetField == actual.TargetField
}

func equalCreateOrchestrationRequest(
	expected *pb.CreateOrchestrationRequest,
	actual *pb.CreateOrchestrationRequest,
) bool {
	return expected.Name == actual.Name &&
		testutil.ValidateEqualElementsString(expected.Links, actual.Links)
}
