package testutil

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/api"
)

// AssetNameForNum deterministically creates an asset name for a given number.
func AssetNameForNum(num int) api.AssetName {
	return api.AssetName(fmt.Sprintf("asset-%v", num))
}

// AssetNameForNumStr deterministically creates an asset name for a given number.
func AssetNameForNumStr(num int) string {
	return string(AssetNameForNum(num))
}

// AssetImageForNum deterministically creates an image for a given number.
func AssetImageForNum(num int) string {
	name := AssetNameForNum(num)
	return fmt.Sprintf("image-%v", name)
}

// StageNameForNum deterministically creates a stage name for a given number.
func StageNameForNum(num int) api.StageName {
	return api.StageName(fmt.Sprintf("stage-%v", num))
}

// StageNameForNumStr deterministically creates a stage name for a given number.
func StageNameForNumStr(num int) string {
	return string(StageNameForNum(num))
}

// StageServiceForNum deterministically creates a stage service for a given
// number.
func StageServiceForNum(num int) string {
	return fmt.Sprintf("service-%v", num)
}

// StageRpcForNum deterministically creates a stage rpc for a given
// number.
func StageRpcForNum(num int) string {
	return fmt.Sprintf("rpc-%v", num)
}

// StageAddressForNum deterministically creates a stage address for a given
// number.
func StageAddressForNum(num int) string {
	return fmt.Sprintf("address-%v", num)
}

// LinkNameForNum deterministically creates a link name for a given number.
func LinkNameForNum(num int) api.LinkName {
	return api.LinkName(fmt.Sprintf("link-%v", num))
}

// LinkNameForNumStr deterministically creates a link name for a given number.
func LinkNameForNumStr(num int) string {
	return string(LinkNameForNum(num))
}

// LinkSourceStageForNum deterministically creates a link source stage for a
// given number.
func LinkSourceStageForNum(num int) api.StageName {
	return StageNameForNum(num)
}

// LinkSourceStageForNumStr deterministically creates a link source stage for a
// given number.
func LinkSourceStageForNumStr(num int) string {
	return string(LinkSourceStageForNum(num))
}

// LinkSourceFieldForNum deterministically creates a link source field for a
// given number.
func LinkSourceFieldForNum(num int) string {
	return fmt.Sprintf("source-field-%v", num)
}

// LinkTargetStageForNum deterministically creates a link target stage for a
// given number.
func LinkTargetStageForNum(num int) api.StageName {
	return StageNameForNum(num + 1)
}

// LinkTargetStageForNumStr deterministically creates a link target stage for a
// given number.
func LinkTargetStageForNumStr(num int) string {
	return string(LinkTargetStageForNum(num + 1))
}

// LinkTargetFieldForNum deterministically creates a link target field for a
// given number.
func LinkTargetFieldForNum(num int) string {
	return fmt.Sprintf("target-field-%v", num)
}

// OrchestrationNameForNum deterministically creates an asset name for a given
// number.
func OrchestrationNameForNum(num int) string {
	return fmt.Sprintf("orchestration-%v", num)
}
