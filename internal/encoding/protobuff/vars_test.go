package protobuff

import apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"

const (
	assetName  apitypes.AssetName = "Asset Name"
	assetImage                    = "user/assetImage:version"

	stageName          = "Stage Name"
	stageAsset         = "Stage Asset"
	stageService       = "stageService"
	stageRpc           = "stageRpc"
	stageAddress       = "stageAddress"
	stageHost          = "stageHost"
	stagePort    int32 = 12345

	linkName        = "linkName"
	linkSourceStage = "linkSourceStage"
	linkSourceField = "linkSourceField"
	linkTargetStage = "linkSourceStage"
	linkTargetField = "linkTargetField"

	orchestrationName = "OrchestrationName"
	oLink1            = "Link Name 1"
	oLink2            = "Link Name 2"
	oLink3            = "Link Name 3"
)
