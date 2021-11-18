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
)

var (
	stage1 = &blueprint.Stage{Name: "Stage Name 1"}
	stage2 = &blueprint.Stage{Name: "Stage Name 2"}
	stage3 = &blueprint.Stage{Name: "Stage Name 3"}

	link1 = &blueprint.Link{SourceField: "Source1"}
	link2 = &blueprint.Link{SourceField: "Source2"}
	link3 = &blueprint.Link{SourceField: "Source3"}

	pbStage1 = &pb.Stage{
		Name:    "Stage Name 1",
		Asset:   "Asset Name 1",
		Service: "Service1",
		Method:  "Method1",
	}
	pbStage2 = &pb.Stage{
		Name:    "Stage Name 2",
		Asset:   "Asset Name 2",
		Service: "Service2",
		Method:  "Method2",
	}
	pbStage3 = &pb.Stage{
		Name:    "Stage Name 3",
		Asset:   "Asset Name 3",
		Service: "Service3",
		Method:  "Method3",
	}

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
