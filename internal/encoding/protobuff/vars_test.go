package protobuff

import (
	"github.com/DuarteMRAlves/maestro/api/pb"
	"github.com/DuarteMRAlves/maestro/internal/blueprint"
)

const (
	assetName  = "Asset Name"
	assetImage = "user/assetImage:version"

	stageId      = "QFE97FW2VD"
	stageName    = "Stage Name"
	stageAssetId = "WG623FEI7V"
	stageService = "stageService"
	stageMethod  = "stageMethod"

	linkSourceId    = "345EVD7FN9"
	linkSourceField = "linkSourceField"
	linkTargetId    = "GSV8S8S7CA"
	linkTargetField = "linkTargetField"

	blueprintId   = "FI97WC67CA"
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
		Id:      &pb.Id{Val: "GI89SVD65V"},
		Name:    "Stage Name 1",
		AssetId: &pb.Id{Val: "DG9GS7SW6D"},
		Service: "Service1",
		Method:  "Method1",
	}
	pbStage2 = &pb.Stage{
		Id:      &pb.Id{Val: "XD76VD9S7C"},
		Name:    "Stage Name 2",
		AssetId: &pb.Id{Val: "DEW98SNF67"},
		Service: "Service2",
		Method:  "Method2",
	}
	pbStage3 = &pb.Stage{
		Id:      &pb.Id{Val: "V9CD7S0XV3"},
		Name:    "Stage Name 3",
		AssetId: &pb.Id{Val: "15FW214F5G"},
		Service: "Service3",
		Method:  "Method3",
	}

	pbLink1 = &pb.Link{
		SourceId:    &pb.Id{Val: "VB87SA9V90"},
		SourceField: "Source1",
		TargetId:    &pb.Id{Val: "BF87VD9ZSG"},
		TargetField: "Target1",
	}
	pbLink2 = &pb.Link{
		SourceId:    &pb.Id{Val: "BF87XS7ASH"},
		SourceField: "Source2",
		TargetId:    &pb.Id{Val: "HL8S87B0TY"},
		TargetField: "Target2",
	}
	pbLink3 = &pb.Link{
		SourceId:    &pb.Id{Val: "XC946BO9P0"},
		SourceField: "Source3",
		TargetId:    &pb.Id{Val: "LF8SD80WF8"},
		TargetField: "Target3",
	}
)
