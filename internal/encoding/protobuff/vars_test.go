package protobuff

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/blueprint"
)

const (
	assetName  = "Asset Name"
	assetImage = "user/assetImage:version"

	stageName    = "Stage Name"
	stageAsset   = "Stage Asset"
	stageService = "stageService"
	stageMethod  = "stageMethod"

	linkSourceStage = "linkSourceStage"
	linkSourceField = "linkSourceField"
	linkTargetStage = "linkSourceStage"
	linkTargetField = "linkTargetField"

	blueprintName = "BlueprintName"
	bpStage1      = "Stage Name 1"
	bpStage2      = "Stage Name 2"
	bpStage3      = "Stage Name 3"
)

var (
	link1 = &blueprint.Link{SourceField: "Source1"}
	link2 = &blueprint.Link{SourceField: "Source2"}
	link3 = &blueprint.Link{SourceField: "Source3"}

	pbLink1 = &pb.Link{
		SourceStage: "SourceStage1",
		SourceField: "Source1",
		TargetStage: "TargetStage1",
		TargetField: "Target1",
	}
	pbLink2 = &pb.Link{
		SourceStage: "SourceStage2",
		SourceField: "Source2",
		TargetStage: "TargetStage2",
		TargetField: "Target2",
	}
	pbLink3 = &pb.Link{
		SourceStage: "SourceStage3",
		SourceField: "Source3",
		TargetStage: "TargetStage3",
		TargetField: "Target3",
	}
)
