package testutil

import (
	"fmt"
)

// AssetNameForNum deterministically creates an asset name for a given number.
func AssetNameForNum(num int) string {
	return fmt.Sprintf("asset-%v", num)
}

// AssetImageForNum deterministically creates an image for a given number.
func AssetImageForNum(num int) string {
	name := AssetNameForNum(num)
	return fmt.Sprintf("image-%v", name)
}

// StageNameForNum deterministically creates a stage name for a given number.
func StageNameForNum(num int) string {
	return fmt.Sprintf("stage-%v", num)
}

// StageServiceForNum deterministically creates a stage service for a given
// number.
func StageServiceForNum(num int) string {
	return fmt.Sprintf("service-%v", num)
}

// StageMethodForNum deterministically creates a stage method for a given
// number.
func StageMethodForNum(num int) string {
	return fmt.Sprintf("method-%v", num)
}

// StageAddressForNum deterministically creates a stage method for a given
// number.
func StageAddressForNum(num int) string {
	return fmt.Sprintf("address-%v", num)
}

// LinkNameForNum deterministically creates a link name for a given number.
func LinkNameForNum(num int) string {
	return fmt.Sprintf("link-%v", num)
}

// LinkSourceStageForNum deterministically creates a link source stage for a
// given number.
func LinkSourceStageForNum(num int) string {
	return StageNameForNum(num)
}

// LinkSourceFieldForNum deterministically creates a link source field for a
// given number.
func LinkSourceFieldForNum(num int) string {
	return fmt.Sprintf("source-field-%v", num)
}

// LinkTargetStageForNum deterministically creates a link target stage for a
// given number.
func LinkTargetStageForNum(num int) string {
	return StageNameForNum(num + 1)
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
