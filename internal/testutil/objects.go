package testutil

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/link"
	"github.com/DuarteMRAlves/maestro/internal/orchestration"
	"github.com/DuarteMRAlves/maestro/internal/stage"
)

// AssetForNum deterministically creates an asset resource with the given
// number.
func AssetForNum(num int) *asset.Asset {
	return &asset.Asset{
		Name:  AssetNameForNum(num),
		Image: AssetImageForNum(num),
	}
}

// StageForNum deterministically creates a stage resource with the given number.
// The associated asset name is the one used in AssetForNum.
func StageForNum(num int) *stage.Stage {
	return &stage.Stage{
		Name:    StageNameForNum(num),
		Asset:   AssetNameForNum(num),
		Service: StageServiceForNum(num),
		Method:  StageMethodForNum(num),
		Address: StageAddressForNum(num),
	}
}

// LinkForNum deterministically creates a link resource with the given number.
// The associated source stage name is the one used in stageForNum with the num
// argument. The associated target stage name is the one used in the stageForNum
// with the num+1 argument.
func LinkForNum(num int) *link.Link {
	return &link.Link{
		Name:        LinkNameForNum(num),
		SourceStage: LinkSourceStageForNum(num),
		SourceField: LinkSourceFieldForNum(num),
		TargetStage: LinkTargetStageForNum(num),
		TargetField: LinkTargetFieldForNum(num),
	}
}

// OrchestrationForNum deterministically creates an orchestration object with
// the given number. The associated links names are created with LinkNameForNum
// with the received nums.
func OrchestrationForNum(
	num int,
	linkNums ...int,
) *orchestration.Orchestration {
	links := make([]string, 0, len(linkNums))
	for _, linkNum := range linkNums {
		links = append(links, LinkNameForNum(linkNum))
	}
	return &orchestration.Orchestration{
		Name:  OrchestrationNameForNum(num),
		Links: links,
	}
}

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
