package orchestration

import (
	"fmt"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
)

type Link struct {
	name        apitypes.LinkName
	sourceStage apitypes.StageName
	sourceField string
	targetStage apitypes.StageName
	targetField string
}

func NewLink(
	name apitypes.LinkName,
	sourceStage apitypes.StageName,
	sourceField string,
	targetStage apitypes.StageName,
	targetField string,
) *Link {
	return &Link{
		name:        name,
		sourceStage: sourceStage,
		sourceField: sourceField,
		targetStage: targetStage,
		targetField: targetField,
	}
}

func (l *Link) Name() apitypes.LinkName {
	return l.name
}

func (l *Link) SourceStage() apitypes.StageName {
	return l.sourceStage
}

func (l *Link) SourceField() string {
	return l.sourceField
}

func (l *Link) TargetStage() apitypes.StageName {
	return l.targetStage
}

func (l *Link) TargetField() string {
	return l.targetField
}

func (l *Link) Clone() *Link {
	return &Link{
		name:        l.name,
		sourceStage: l.sourceStage,
		sourceField: l.sourceField,
		targetStage: l.targetStage,
		targetField: l.targetField,
	}
}

func (l *Link) ToApi() *apitypes.Link {
	return &apitypes.Link{
		Name:        l.name,
		SourceStage: l.sourceStage,
		SourceField: l.sourceField,
		TargetStage: l.targetStage,
		TargetField: l.targetField,
	}
}

func (l *Link) String() string {
	return fmt.Sprintf(
		"Link{Name:%v,SourceStage:%v,SourceField:%v,TargetStage:%v,TargetField:%v",
		l.name,
		l.sourceStage,
		l.sourceField,
		l.targetStage,
		l.targetField)
}
